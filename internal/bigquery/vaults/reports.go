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

type Unbond struct {
	Amount             int
	IngestionTimestamp time.Time
}

type UpdateIndex struct {
	VaultTokenBalance  int
	IngestionTimestamp time.Time
}

// QueryDailyReport returns a file with the last 24h statistics of a given vaultAddress
func QueryDailyReport(blockHeight int, addressQuery string, outputFormat string) error {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}
	if blockHeight < 1 {
		log.Fatal("Block height should be higher than 0")
	}

	filename := fmt.Sprintf("%s_%s", bigquerytypes.PrefixBigQuery+bigquerytypes.PrependQueryDailyReport, addressQuery)

	rewardsUpdateUser, err := queryDailyReportRewardsUpdateUser(blockHeight, addressQuery)
	if err != nil {
		return err
	}

	// Response variables
	// - General
	var generalUsersBonded, generalUsersExited, generalUnbondAmountPending, generalAverageBondAmount, generalAverageTxNumber int
	// - Latest 24h
	var dailyBondNewUsersCount, dailyBondNewUsersAmount, dailyBondOldUsersCount, dailyBondOldUsersAmount int
	var dailyUnbondUsersCount, dailyUnbondUsersAmount int
	dailyExitUsersCount := make(map[string]bool)
	var dailyExitUsersAmount int
	// - Wall of fame
	var biggestSingleDepositor, biggestHolder string
	var biggestSingleDeposit, biggestBalance int

	// Utility variables
	var totalBondCount, totalBondAmount, totalTxCount int
	totalUnbondCount := make(map[int]Unbond)
	userFirstDeposit := make(map[string]time.Time)

	// Computation of statistics
	for user, transactions := range rewardsUpdateUser {
		// double check transactions are sorted by time
		sort.Slice(transactions, func(i, j int) bool {
			return transactions[i].IngestionTimestamp.Before(transactions[j].IngestionTimestamp)
		})

		// declaring a 0 previousBalance foreach user we iterate its transactions
		var previousBalance int
		for i, transaction := range transactions {
			// incr total transaction count
			totalTxCount++

			// state bond or unbond based on balance change
			change := transaction.VaultTokenBalance - previousBalance
			// state if the current tx is older than 24h from now
			isDailyTransaction := time.Since(transaction.IngestionTimestamp).Hours() <= 24

			// Bond if change respect previous balance is positive
			if change > 0 {
				// increase total bond count and amount
				totalBondCount++
				totalBondAmount += change

				// Check if this is user's first deposit
				if _, ok := userFirstDeposit[user]; !ok {
					// new user's deposit as this is the first deposit for him
					if isDailyTransaction {
						// only if since less than 24h we increase daily variables
						dailyBondNewUsersCount++
						dailyBondNewUsersAmount += change
					}
					// we set the ingestion timestamp for the user's first bond
					userFirstDeposit[user] = transaction.IngestionTimestamp
					// increase total bonded users count
					generalUsersBonded++
				} else if ok && isDailyTransaction { //old user's deposit
					// if we are here we know this is NOT the first bond
					// so only if the first bond is since more than 24h we increase daily variables
					if time.Since(userFirstDeposit[user]).Hours() > 24 {
						dailyBondOldUsersCount++
						dailyBondOldUsersAmount += change
					}
					// TODO: check if else needed in order to increase general stats
				}

				// check if is the biggest single deposit
				if change > biggestSingleDeposit {
					biggestSingleDeposit = change
					biggestSingleDepositor = user
				}
			} else if change <= 0 { // else is an Unbond
				// increase total unbond count and amount
				totalUnbondCount[len(totalUnbondCount)] = Unbond{ // this struct serves for later use
					Amount:             -change,
					IngestionTimestamp: transaction.IngestionTimestamp,
				}
				//totalUnbondAmount += -change // Convert negative to positive TODO check if needced or not! if yes reimplement above on vars initialization

				// we check if the unbond is from the latest 24h
				if isDailyTransaction {
					// so we increase the count and amount
					dailyUnbondUsersCount++
					dailyUnbondUsersAmount += -change // Convert negative to positive
				}
				// then we check if user completely exited by checking if the current tx balance is 0 && is the latest user tx
				if transaction.VaultTokenBalance == 0 && i == len(transactions)-1 {
					if isDailyTransaction {
						// filling map for future computation outside this whole for loop
						dailyExitUsersCount[user] = true
					}
					// increasing general count of users exited completely regardless the timeframe of 24h
					generalUsersExited++
				}
			}

			// update previous balance for next iteration
			previousBalance = transaction.VaultTokenBalance
		}

		// here we take the latest balance of the current iter user to determine who is the biggest current holder of shares
		userBalance := transactions[len(transactions)-1].VaultTokenBalance
		if userBalance > biggestBalance {
			biggestBalance = userBalance
			biggestHolder = user
		}
	}

	// Count exited amount in the last 24h
	for user := range dailyExitUsersCount {
		if time.Since(rewardsUpdateUser[user][len(rewardsUpdateUser[user])-1].IngestionTimestamp).Hours() <= 24 { // taking -1 as it is the exit itself
			dailyExitUsersAmount += rewardsUpdateUser[user][len(rewardsUpdateUser[user])-2].VaultTokenBalance // taking -2 as it is the previous to latest tx by user whom exited completely
		}
	}

	// check for unbonds pending 14 days
	for _, unbond := range totalUnbondCount {
		// if the unbond is from less than twoWeeksInHours means that is pending to claim (osmosis pools unbond time)
		if time.Since(unbond.IngestionTimestamp).Hours() < 336 { // 336 = 24hrs * 14days
			generalUnbondAmountPending += unbond.Amount
		}
	}

	// Compute average bond amount and tx number per user
	if totalBondCount != 0 {
		generalAverageBondAmount = totalBondAmount / totalBondCount
	}
	if len(rewardsUpdateUser) != 0 {
		generalAverageTxNumber = totalTxCount / len(rewardsUpdateUser)
	}

	headers := []string{
		// Latest 24h
		"24_bond_new_users_count",
		"24_bond_new_users_amount",
		"24_bond_old_users_count",
		"24_bond_old_users_amount",
		"24_unbond_users_count",
		"24_unbond_users_amount", // unbond from old users (no new users here, they are already known since bond)
		"24_exit_users_count",
		"24_exit_users_amount", // complete exits from old users (the ones whom exited completely in the last 24h)
		// General
		"general_users_bonded",          // this is a general count of bonding users so far since the start of the vault
		"general_users_exited",          // this is a general count of exited users so far since the start of the vault
		"general_users_active",          // this is a general count of actually active users
		"general_unbond_amount_pending", // the unbonding amount pending to be claimable
		"general_users_average_bond_amount",
		"general_users_average_tx_number",
		// Wall of fame
		"wall_biggest_deposit_user",
		"wall_biggest_deposit_amount",
		"wall_biggest_holder_user",
		"wall_biggest_holder_amount",
	}
	rows := [][]string{
		{
			// Latest 24h
			strconv.Itoa(dailyBondNewUsersCount),
			strconv.Itoa(dailyBondNewUsersAmount), // TODO: here we should convert shares in to denom
			strconv.Itoa(dailyBondOldUsersCount),
			strconv.Itoa(dailyBondOldUsersAmount), // TODO: here we should convert shares in to denom
			strconv.Itoa(dailyUnbondUsersCount),
			strconv.Itoa(dailyUnbondUsersAmount), // TODO: here we should convert shares in to denom
			strconv.Itoa(len(dailyExitUsersCount)),
			strconv.Itoa(dailyExitUsersAmount), // TODO: here we should convert shares in to denom
			// General
			strconv.Itoa(generalUsersBonded),
			strconv.Itoa(generalUsersExited),
			strconv.Itoa(generalUsersBonded - generalUsersExited),
			strconv.Itoa(generalUnbondAmountPending),
			strconv.Itoa(generalAverageBondAmount), // TODO: here we should convert shares in to denom
			strconv.Itoa(generalAverageTxNumber),
			// Wall of fame
			biggestSingleDepositor,
			strconv.Itoa(biggestSingleDeposit),
			biggestHolder,
			strconv.Itoa(biggestBalance),
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

func queryDailyReportRewardsUpdateUser(blockHeight int, addressQuery string) (map[string][]UpdateIndex, error) {
	fmt.Println(fmt.Sprintf(bigquerytypes.QueryDailyReportRewardsUpdateUser, blockHeight, addressQuery))
	_, rows, err := internal.ExecuteQueryAndFetchRows(fmt.Sprintf(bigquerytypes.QueryDailyReportRewardsUpdateUser, blockHeight, addressQuery), "", false)
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
