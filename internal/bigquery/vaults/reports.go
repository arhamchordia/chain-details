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

type UpdateIndex struct {
	VaultTokenBalance  int
	IngestionTimestamp time.Time
}

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

	// Response variables
	// - General
	var generalUsersBonded, generalUsersExited, generalAverageBondAmount, generalAverageTxNumber int
	// - Latest 24h
	var dailyBondNewUsersCount, dailyBondNewUsersAmount, dailyBondOldUsersCount, dailyBondOldUsersAmount int
	var dailyUnbondUsersCount, dailyUnbondUsersAmount, dailyExitUsersAmount int
	// - Wall of fame
	var biggestSingleDepositor, biggestHolder string
	var biggestSingleDeposit, biggestBalance int

	// Utility variables
	var totalBondAmount, totalTxCount int
	dailyExitUsersCount := make(map[string]bool)
	userFirstDeposit := make(map[string]time.Time)

	// Computation of statistics
	for user, transactions := range rewardsUpdateUser {
		// double check transactions are sorted by time
		sort.Slice(transactions, func(i, j int) bool {
			return transactions[i].IngestionTimestamp.Before(transactions[j].IngestionTimestamp)
		})

		// declaring a 0 previousBalance foreach user we iterate its transactions
		var previousBalance int
		for _, transaction := range transactions {
			// incr total transaction count
			totalTxCount++

			// state bond or unbond based on balance change
			change := transaction.VaultTokenBalance - previousBalance
			if change > 0 { // Bond
				// increas total bond amount
				totalBondAmount += change

				// Check if this is user's first deposit
				// TODO: this is wrong, wtf
				if _, ok := userFirstDeposit[user]; !ok {
					// new user's deposit
					dailyBondNewUsersCount++          // TODO: this should increase only if less than 24h ago, not related to the userFirstDeposit
					dailyBondNewUsersAmount += change // TODO: this should increase only if less than 24h ago, not related to the userFirstDeposit
					userFirstDeposit[user] = transaction.IngestionTimestamp
					// increase total bonded users count
					generalUsersBonded++
				} else {
					//old user's deposit
					dailyBondOldUsersCount++          // TODO: this should increase only if less than 24h ago, not related to the userFirstDeposit
					dailyBondOldUsersAmount += change // TODO: this should increase only if less than 24h ago, not related to the userFirstDeposit
				}

				// check if is the biggest single deposit
				if change > biggestSingleDeposit {
					biggestSingleDeposit = change
					biggestSingleDepositor = user
				}
			} else if change <= 0 { // Unbond TODO: double check if <= is fine or we need <
				// TODO: the dailyUnbondUsersCount and dailyUnbondUsersAmount increment should be done only if less than 24h
				time.Since(transaction.IngestionTimestamp)
				dailyUnbondUsersCount++
				dailyUnbondUsersAmount += -change // Convert negative to positive
				// check if user completely exited
				// TODO consider that this worth only if in addition to be 0 it is the latest transaction of the user, or he could have joined again afterward
				if transaction.VaultTokenBalance == 0 && !dailyExitUsersCount[user] {
					dailyExitUsersCount[user] = true // TODO: this should be set only if less than 24h
					generalUsersExited++
				}
			}

			// update previous balance for next iteration
			previousBalance = transaction.VaultTokenBalance
		}

		userBalance := transactions[len(transactions)-1].VaultTokenBalance
		if userBalance > biggestBalance {
			biggestBalance = userBalance
			biggestHolder = user
		}
	}

	// Compute average bond amount and tx number per user
	if dailyBondNewUsersCount+dailyBondOldUsersCount != 0 {
		generalAverageBondAmount = totalBondAmount / (dailyBondNewUsersCount + dailyBondOldUsersCount)
	}
	if len(rewardsUpdateUser) != 0 {
		generalAverageTxNumber = totalTxCount / len(rewardsUpdateUser)
	}

	// Count exited users and amount in the last 24h
	for user, exited := range dailyExitUsersCount {
		if exited && time.Since(userFirstDeposit[user]).Hours() < 24 {
			dailyExitUsersAmount += rewardsUpdateUser[user][len(rewardsUpdateUser[user])-1].VaultTokenBalance
		}
	}

	headers := []string{
		// General
		"general_users_bonded", // this is a general count of bonding users so far since the start of the vault
		"general_users_exited", // this is a general count of exited users so far since the start of the vault
		"general_users_active", // this is a general count of actually active users
		"general_users_average_bond_amount",
		"general_users_average_tx_number",
		// Latest 24h
		"24_bond_new_users_count", "24_bond_new_users_amount", "24_bond_old_users_count", "24_bond_old_users_amount", // bond from new and old users
		"24_unbond_users_count", "24_unbond_users_amount", // unbond from old users (no new users here, they are already known since bond)
		"24_exit_users_count", "24_exit_users_amount", // complete exits from old users (the ones whom exited completely in the last 24h)
		// Wall of fame
		"wall_biggest_deposit_user", "wall_biggest_deposit_amount", // biggest deposit
		"wall_biggest_holder_user", "wall_biggest_holder_amount", // biggest hodler
	}
	rows := [][]string{
		{
			// General
			strconv.Itoa(generalUsersBonded),
			strconv.Itoa(generalUsersExited),
			strconv.Itoa(generalUsersBonded - generalUsersExited),
			strconv.Itoa(generalAverageBondAmount),
			strconv.Itoa(generalAverageTxNumber),
			// Latest 24h
			strconv.Itoa(dailyBondNewUsersCount), strconv.Itoa(dailyBondNewUsersAmount), strconv.Itoa(dailyBondOldUsersCount), strconv.Itoa(dailyBondOldUsersAmount),
			strconv.Itoa(dailyUnbondUsersCount), strconv.Itoa(dailyUnbondUsersAmount),
			strconv.Itoa(len(dailyExitUsersCount)), strconv.Itoa(dailyExitUsersAmount),
			// Wall of fame
			biggestSingleDepositor, strconv.Itoa(biggestSingleDeposit),
			biggestHolder, strconv.Itoa(biggestBalance),
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

		// TODO: we could remove the blockHeight from here which is uncomputed and slowing things
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
