package vaults

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	"github.com/arhamchordia/chain-details/internal/export"
	bigquerytypes "github.com/arhamchordia/chain-details/types/bigquery"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// LPQueryBond returns a file with the bond events in all the blocks
func LPQueryBond(addressQuery string, confirmedQuery bool, pendingQuery bool, outputFormat string) error {
	if confirmedQuery && pendingQuery {
		return fmt.Errorf("--confirmed and --pending flags cannot be used together")
	}

	addressFilterString := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependLPQueryVaultsBond
	if len(addressQuery) > 0 {
		addressFilterString = fmt.Sprintf(bigquerytypes.QueryLPVaultsBondAddressFilter, addressQuery)
		filename = fmt.Sprintf("%s_%s", filename, addressQuery)
	}
	if confirmedQuery {
		filename = fmt.Sprintf("%s_%s", filename, "confirmed")
	}
	if pendingQuery {
		filename = fmt.Sprintf("%s_%s", filename, "pending")
	}

	headers, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryLPVaultsBond, addressFilterString, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if confirmedQuery || pendingQuery {
		_, bondResponses, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryLPVaultsBondResponseFilter, "", false)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// creating map storing the bond_id and share amounts from the second query
		bondResponsesIDs := make(map[string]int)
		for _, row := range bondResponses {
			bondID := row[0]
			shareAmounts := strings.Split(row[1], ", ") // assuming share_amounts is at index 1, consider changing this if the query changes
			bondResponsesIDs[bondID] = len(shareAmounts)
		}

		_, bondShareAmountsTxIds, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryLPVaultsBondShareAmountsTxIds, "", false)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Convert txIDs into a slice of strings for the query
		var txIDs []string
		for _, row := range bondShareAmountsTxIds {
			txIDs = append(txIDs, fmt.Sprintf("'%s'", row[0]))
		}
		txIDsStr := strings.Join(txIDs, ", ")

		// Running the second query SELECT message FROM numia-data.quasar.quasar_tx_messages WHERE tx_id IN (%s)
		_, messageRows, err := internal.ExecuteQueryAndFetchRows(fmt.Sprintf("SELECT message FROM numia-data.quasar.quasar_tx_messages WHERE tx_id IN (%s)", txIDsStr), "", false)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Create a map to store bond_ids that are included in the second query
		includedInSecondQuery := make(map[string]bool)

		// Create an array of bond_ids from the JSON struct retrieved, and decoding the "msg" from base64 if it is encoded, otherwise work with JSON already
		for _, row := range messageRows {
			var jsonData map[string]interface{}
			var jsonMsgData map[string]interface{}
			var data []byte

			message := row[0]

			if err := json.Unmarshal([]byte(message), &jsonData); err != nil {
				continue
			}

			msg, ok := jsonData["msg"]
			if !ok {
				continue
			}

			switch v := msg.(type) {
			case string:
				decodedMessage, err := base64.StdEncoding.DecodeString(v)
				if err == nil {
					data = decodedMessage
				} else {
					data = []byte(v)
				}
			case map[string]interface{}:
				jsonMsgData, err := json.Marshal(v)
				if err != nil {
					continue
				}
				data = jsonMsgData
			default:
				continue
			}

			if err := json.Unmarshal(data, &jsonMsgData); err != nil {
				continue
			}

			callbacks, ok := jsonMsgData["callbacks"].([]interface{})
			if ok {
				for _, callback := range callbacks {
					callbackMap, ok := callback.(map[string]interface{})
					if ok {
						bondID, ok := callbackMap["bond_id"].(string)
						if ok {
							includedInSecondQuery[bondID] = true
						}
					}
				}
			}
		}

		filteredRows := [][]string{}
		bondIDRegex := regexp.MustCompile(`bond_id (\d+)`)
		for _, row := range rows {
			column3 := row[2]
			match := bondIDRegex.FindStringSubmatch(column3)
			if len(match) > 1 && bondIDRegex.MatchString(column3) {
				bondID := match[1]
				shareAmounts, exists := bondResponsesIDs[bondID]
				_, included := includedInSecondQuery[bondID]
				if confirmedQuery {
					if (exists && shareAmounts >= 3) || included {
						filteredRows = append(filteredRows, row)
					}
				} else if pendingQuery {
					if (!exists && !included) || (exists && shareAmounts < 3 && !included) {
						filteredRows = append(filteredRows, row)
					}
				}
			}
		}
		rows = filteredRows
	}

	err = export.ExportFile(outputFormat, filename, headers, rows)
	if err != nil {
		return err
	}

	return nil
}

// LPQueryUnbond returns a file with the unbond events in all the blocks
func LPQueryUnbond(addressQuery string, confirmedQuery bool, pendingQuery bool, outputFormat string) error {
	addressFilterString := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependLPQueryVaultsUnbond
	if len(addressQuery) > 0 {
		addressFilterString = fmt.Sprintf(bigquerytypes.QueryLPVaultsUnbondAddressFilter, addressQuery)
		filename = fmt.Sprintf("%s_%s", filename, addressQuery)
	}
	if confirmedQuery {
		filename = fmt.Sprintf("%s_%s", filename, "confirmed")
	}
	if pendingQuery {
		filename = fmt.Sprintf("%s_%s", filename, "pending")
	}

	headers, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryLPVaultsUnbond, addressFilterString, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if confirmedQuery || pendingQuery {
		_, confirmedRows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryLPVaultsUnbondConfirmedFilter, "", false)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// creating map storing the bond_id from the second query
		confirmedUnbondIDs := make(map[string]bool)
		for _, row := range confirmedRows {
			bondID := row[0]
			confirmedUnbondIDs[bondID] = true
		}

		// filtering rows from the first query by checking if bond_id exists
		filteredRows := [][]string{}
		bondIDRegex := regexp.MustCompile(`bond_id (\d+)`)
		for _, row := range rows {
			column3 := row[2]
			match := bondIDRegex.FindStringSubmatch(column3)
			if len(match) > 1 && bondIDRegex.MatchString(column3) {
				bondID := match[1]
				if confirmedQuery {
					if _, exists := confirmedUnbondIDs[bondID]; exists {
						filteredRows = append(filteredRows, row)
					}
				} else if pendingQuery {
					if _, exists := confirmedUnbondIDs[bondID]; !exists {
						filteredRows = append(filteredRows, row)
					}
				}
			}
		}
		rows = filteredRows
	}

	err = export.ExportFile(outputFormat, filename, headers, rows)
	if err != nil {
		return err
	}

	return nil
}

// LPQueryClaim returns a file with the claim events in all the blocks
func LPQueryClaim(addressQuery string, outputFormat string) error {
	addressFilterString := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependLPQueryVaultsClaim

	if len(addressQuery) > 0 {
		addressFilterString = fmt.Sprintf(bigquerytypes.QueryLPVaultsClaimAddressFilter, addressQuery)
		filename = fmt.Sprintf("%s_%s", filename, addressQuery)
	}

	headers, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryLPVaultsClaim, addressFilterString, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = export.ExportFile(outputFormat, filename, headers, rows)
	if err != nil {
		return err
	}

	return nil
}

type Unbond struct {
	Amount             int
	IngestionTimestamp time.Time
}

type UpdateIndex struct {
	VaultTokenBalance  int
	IngestionTimestamp time.Time
}

// LPQueryReport returns a file with the last 24h statistics of a given vaultAddress
func LPQueryReport(blockHeight int, addressQuery string, outputFormat string) error {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}
	if blockHeight < 1 {
		log.Fatal("Block height should be higher than 0")
	}

	filename := fmt.Sprintf("%s_%s", bigquerytypes.PrefixBigQuery+bigquerytypes.PrependLPQueryReport, addressQuery)

	rewardsUpdateUser, err := queryReportRewardsUpdateUser(blockHeight, addressQuery)
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
	var totalBondCount, totalBondAmount, totalUnbondAmount, totalTxCount int
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
						// only if since less than 24h we increase variables
						dailyBondNewUsersCount++
						dailyBondNewUsersAmount += change
					}
					// we set the ingestion timestamp for the user's first bond
					userFirstDeposit[user] = transaction.IngestionTimestamp
					// increase total bonded users count
					generalUsersBonded++
				} else if ok && isDailyTransaction { //old user's deposit
					// if we are here we know this is NOT the first bond
					// so only if the first bond is since more than 24h we increase variables
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
				totalUnbondAmount += -change
				totalUnbondCount[len(totalUnbondCount)] = Unbond{ // this struct serves for later use
					Amount:             -change,
					IngestionTimestamp: transaction.IngestionTimestamp,
				}

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
		"general_amount_active",         // this is a general count of actually active users
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
			strconv.Itoa(totalBondAmount - totalUnbondAmount),
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

	err = export.ExportFile(outputFormat, filename, headers, rows)
	if err != nil {
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

func queryReportRewardsUpdateUser(blockHeight int, addressQuery string) (map[string][]UpdateIndex, error) {
	_, rows, err := internal.ExecuteQueryAndFetchRows(fmt.Sprintf(bigquerytypes.QueryLPReportRewardsUpdateUser, blockHeight, addressQuery), "", false)
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
