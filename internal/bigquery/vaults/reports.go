package vaults

import (
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	bigquerytypes "github.com/arhamchordia/chain-details/types/bigquery"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// QueryDailyReport returns a file with the last 24h statistics of a given vaultAddress
func QueryDailyReport(addressQuery string, outputFormat string) error {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}

	filename := fmt.Sprintf("%s_%s", bigquerytypes.PrefixBigQuery+bigquerytypes.PrependQueryDailyReport, addressQuery)

	rewardsUpdateUser, err := queryDailyReportRewardsUpdateUser(addressQuery)
	if err != nil {
		return err
	}

	// init stat vars for bond
	var bondNewUsersCount, bondNewUsersAmount, bondOldUsersCount, bondOldUsersAmount int
	// init stat vars for unbond
	var unbondNewUsersCount, unbondNewUsersAmount, unbondOldUsersCount, unbondOldUsersAmount int
	// init stat vars for bond
	var bondTotalUsers int
	biggestSingleDeposit := 0
	biggestSingleDepositor := ""
	biggestHolder := ""
	biggestBalance := 0

	// init user deposit history
	userFirstDeposit := make(map[string]time.Time)

	for user, transactions := range rewardsUpdateUser {
		// double check transactions are sorted by time
		sort.Slice(transactions, func(i, j int) bool {
			return transactions[i].IngestionTimestamp.Before(transactions[j].IngestionTimestamp)
		})

		previousBalance := 0
		for _, transaction := range transactions {
			// compute bond or unbond based on balance change
			change := transaction.VaultTokenBalance - previousBalance
			if change > 0 { // Bond
				// if this is user's first deposit
				if _, ok := userFirstDeposit[user]; !ok {
					// new user's deposit
					bondNewUsersCount++
					bondNewUsersAmount += change
					userFirstDeposit[user] = transaction.IngestionTimestamp
				} else {
					// old user's deposit
					bondOldUsersCount++
					bondOldUsersAmount += change
				}

				// Check if this is the biggest single deposit
				if change > biggestSingleDeposit {
					biggestSingleDeposit = change
					biggestSingleDepositor = user
				}
			} else if change < 0 { // Unbond
				// This is an user's withdrawal
				unbondOldUsersCount++
				unbondOldUsersAmount += -change // Convert negative to positive
			}

			// Update previous balance
			previousBalance = transaction.VaultTokenBalance
		}

		// check if this user is the biggest hodler
		if previousBalance > biggestBalance {
			biggestBalance = previousBalance
			biggestHolder = user
		}

		// count total users with not 0 balance
		if previousBalance > 0 {
			bondTotalUsers++
		}
	}

	// Count new users' withdrawals
	for user, firstDepositTime := range userFirstDeposit {
		// If first deposit was less than 24 hours ago, this is a new user's withdrawal
		if time.Since(firstDepositTime).Hours() < 24 {
			unbondNewUsersCount++
			unbondNewUsersAmount += rewardsUpdateUser[user][len(rewardsUpdateUser[user])-1].VaultTokenBalance
		}
	}

	// Print wall of fame
	fmt.Printf("Biggest single deposit: %d by user %s\n", biggestSingleDeposit, biggestSingleDepositor)
	fmt.Printf("Biggest actual holder: %s with balance %d\n", biggestHolder, biggestBalance)

	headers := []string{
		// General
		"general_users_bonded",              // this is a general count of bonding users so far since the start of the vault
		"general_users_exited",              // this is a general count of exited users so far since the start of the vault
		"general_users_average_bond_amount", // this is a general average of tx per user regardless are bond or unbonds
		"general_users_average_tx_number",   // this is a general average of tx per user regardless are bond or unbonds
		// Latest 24h
		"bond_new_users_count", "bond_new_users_amount", "bond_old_users_count", "bond_old_users_amount", // bond from new and old users
		"unbond_users_count", "unbond_users_amount", // unbond from old users (no new users here, they are already known since bond)
		"exit_users_count", "exit_users_amount", // complete exits from old users (the ones whom exited completely in the last 24h)
		// Wall of fame
		"biggest_deposit_user", "biggest_deposit_amount", // biggest deposit
		"biggest_holder_user", "biggest_holder_amount", // biggest hodler
	}
	rows := [][]string{
		{
			// Bond
			strconv.Itoa(bondNewUsersCount), strconv.Itoa(bondNewUsersAmount), strconv.Itoa(bondOldUsersCount), strconv.Itoa(bondOldUsersAmount),
			// Unbond
			strconv.Itoa(unbondNewUsersCount), strconv.Itoa(unbondNewUsersAmount), strconv.Itoa(unbondOldUsersCount), strconv.Itoa(unbondOldUsersAmount),
			// General
			strconv.Itoa(bondTotalUsers),
		},
	}

	if outputFormat == "csv" {
		err = internal.WriteCSV(filename, headers, rows)
	} else {
		err = internal.WriteJSON(filename, headers, rows)
	}
	if err != nil {
		log.Printf("Warning: %v", err)
		return err
	}

	return nil
}

type UpdateIndex struct {
	VaultTokenBalance  int
	IngestionTimestamp time.Time
}

func parseTransactions(input string) ([][]string, error) {
	r := regexp.MustCompile(`\[(\d+) (\d+) ([^\]]+)\]`)
	matches := r.FindAllStringSubmatch(input, -1)
	if matches == nil {
		return nil, fmt.Errorf("no matches found")
	}

	transactions := make([][]string, 0)
	for _, match := range matches {
		transaction := []string{match[1], match[2], match[3]}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func queryDailyReportRewardsUpdateUser(addressQuery string) (map[string][]UpdateIndex, error) {
	_, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryDailyReportRewardsUpdateUser, addressQuery, true)
	if err != nil {
		log.Fatalf("%v", err)
		return nil, err
	}

	Addresses := make(map[string][]UpdateIndex)

	for _, row := range rows {
		user := row[0] // is i.e. quasar103dsgfltsaykm0x4sd0mf4yj3wjht9ruyv3ckl

		transactions, err := parseTransactions(row[1]) // is i.e. [[211310 1919557 "2023-04-05 16:38:38 +0000 UTC"] [286844 0 "2023-04-10 15:07:56 +0000 UTC"]]
		if err != nil {
			log.Fatalf("%v", err)
			return nil, err
		}

		for _, transaction := range transactions {
			VaultTokenBalance, err := strconv.Atoi(transaction[1])
			if err != nil {
				log.Fatalf("%v", err)
				return nil, err
			}

			timestamp, err := time.Parse("2006-01-02 15:04:05 -0700 MST", strings.TrimSpace(transaction[2]))
			if err != nil {
				log.Fatalf("%v", err)
				return nil, err
			}
			Addresses[user] = append(Addresses[user], UpdateIndex{
				VaultTokenBalance:  VaultTokenBalance,
				IngestionTimestamp: timestamp,
			})
		}
	}

	return Addresses, nil
}
