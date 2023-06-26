package vaults

import (
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	bigquerytypes "github.com/arhamchordia/chain-details/types/bigquery"
	"log"
	"strconv"
	"strings"
)

// QueryDailyReportBond returns a file with the last 24h statistics about:
// - Number of new users bonding
// - Amount bonded by new users
// - Number of old users bonding
// - Amount bonded by old users
// - Number of total users till today
func QueryDailyReportBond(addressQuery string, outputFormat string) error {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}

	filename := fmt.Sprintf("%s_%s", bigquerytypes.PrefixBigQuery+bigquerytypes.PrependQueryDailyReportBond, addressQuery)

	// query all the bonds before last 24h with distinct on userAddressSender,
	// this will retrieve a list of all the current bonders without repeated values
	_, rowsBefore, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryDailyReportBondBefore, addressQuery, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Creating a string/bool map for faster checking after " character removal (Numia is retrieving address sometimes embedded by double quotes)
	AddressBefore := make(map[string]bool)
	for _, address := range rowsBefore {
		tempAddr := strings.ReplaceAll(strings.Join(address, ""), "\"", "")
		AddressBefore[tempAddr] = true
	}

	// query all the bonds after last 24h without distinct, then compare it with the previous result set
	// and compute who is new, who not and the related amounts
	_, rowsAfter, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryDailyReportBondAfter, addressQuery, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Variables to increment
	newUsersCount := 0
	newUsersAmount := 0
	oldUsersCount := 0
	oldUsersAmount := 0

	// AddressAfter map to check duplicates in rowsAfter
	AddressAfter := make(map[string]bool)

	for _, bond := range rowsAfter {
		tempAddr := strings.ReplaceAll(bond[0], "\"", "")
		// looking in AddressBefore to check if the user was already existing before 24h or not
		_, ok := AddressBefore[tempAddr]
		if ok {
			// found is an old user
			parseInt, err := strconv.ParseInt(bond[1][:len(bond[1])-2], 10, 64)
			if err != nil {
				return err
			}
			// check if address is already counted
			if _, ok := AddressAfter[tempAddr]; !ok {
				oldUsersCount++
				AddressAfter[tempAddr] = true
			}
			oldUsersAmount += int(parseInt)
		} else {
			// not found is a new user
			parseInt, err := strconv.ParseInt(bond[1][:len(bond[1])-2], 10, 64)
			if err != nil {
				return err
			}
			// check if address is already counted
			if _, ok := AddressAfter[tempAddr]; !ok {
				newUsersCount++
				AddressAfter[tempAddr] = true
			}
			newUsersAmount += int(parseInt)
		}
	}

	// Variables to compute
	totalUsers := len(rowsBefore) + newUsersCount

	// Headers and rows
	headers := []string{"new_users_count", "new_users_amount", "old_users_count", "old_users_amount", "total_users"}
	rows := [][]string{
		{strconv.Itoa(newUsersCount), strconv.Itoa(newUsersAmount), strconv.Itoa(oldUsersCount), strconv.Itoa(oldUsersAmount), strconv.Itoa(totalUsers)},
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
