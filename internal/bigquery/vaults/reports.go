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
func QueryDailyReportBond(outputFormat string) error {
	//filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependQueryDailyReportBond

	// query all the bonds before last 24h with distinct on userAddressSender,
	// this will retrieve a list of all the current bonders without repeated values
	_, rowsBefore, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryDailyReportBondBefore, "", false)
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
	_, rowsAfter, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryDailyReportBondAfter, "", false)
	if err != nil {
		log.Fatalf("%v", err)
	}

	NewUsersCount := 0
	NewUsersAmount := 0

	OldUsersCount := 0
	OldUsersAmount := 0

	for _, bond := range rowsAfter {
		tempAddr := strings.ReplaceAll(bond[0], "\"", "")
		_, ok := AddressBefore[tempAddr]
		if ok {
			// found is an old user
			parseInt, err := strconv.ParseInt(bond[1][:len(bond[1])-2], 10, 64)
			if err != nil {
				return err
			}
			OldUsersCount++
			OldUsersAmount += int(parseInt)
		} else {
			// not found is a new user
			parseInt, err := strconv.ParseInt(bond[1][:len(bond[1])-2], 10, 64)
			if err != nil {
				return err
			}
			NewUsersCount++
			NewUsersAmount += int(parseInt)
		}
	}

	TotalNumberUsers := len(rowsBefore) + NewUsersCount

	fmt.Println(NewUsersCount, NewUsersAmount, OldUsersCount, OldUsersAmount, TotalNumberUsers)

	if outputFormat == "csv" {
		//err = internal.WriteCSV(filename, headers, rows)
	} else {
		//err = internal.WriteJSON(filename, headers, rows)
	}
	if err != nil {
		log.Printf("Warning: %v", err)
		return err
	}

	return nil
}
