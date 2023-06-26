package vaults

import (
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	bigquerytypes "github.com/arhamchordia/chain-details/types/bigquery"
	"log"
	"strconv"
	"strings"
)

type DailyReportBond struct {
	NewUsersCount  int
	NewUsersAmount int
	OldUsersCount  int
	OldUsersAmount int
	TotalUsers     int
}

type DailyReportUnbond struct {
	NewUsersCount  int
	NewUsersAmount int
	OldUsersCount  int
	OldUsersAmount int
	TotalUsers     int
}

type DailyReportClaim struct {
	NewUsersCount int
	OldUsersCount int
	TotalUsers    int
}

//TODO:
// - Remove distinct from BondBefore query in order to be able ask for block_timestamp and OSMO bonding amount per-bond
// - Merge all the functions in one to post-compute between different actions as bond, unbond and claim
// - - Query all the bonds before 24h
// - - Query all the bonds in the last 24h
// - - Query all the unbonds before 24h
// - - Query all the unbonds in the last 24
// - - Query all the claims before 24h
// - - Query all the claims in the last 24h
// - - - Choose if we want to convert OSMO to Shares during bonding, or Shares to OSMO during unbond/claim
// - - - Use the converted value to post compute between bonds and unbonds to know which user exited completely and which one not
// - - - Use the claims to check if in the effective claim timestamp the user had pending unbonds to know if the claim was about all the unbonds or not

// QueryDailyReport returns a file with the last 24h statistics of a given vaultAddress
func QueryDailyReport(addressQuery string, outputFormat string) error {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}

	filename := fmt.Sprintf("%s_%s", bigquerytypes.PrefixBigQuery+bigquerytypes.PrependQueryDailyReport, addressQuery)

	// Run queries and get return values
	bond, err := queryDailyReportBond(addressQuery)
	if err != nil {
		return err
	}
	// TODO: convertDenomToSharesAmount(denomAmount) iteration

	// TODO: check against the contracts rewards as sender and shares_amount

	unbond, err := queryDailyReportUnbond(addressQuery)
	if err != nil {
		return err
	}
	claim, err := queryDailyReportClaim(addressQuery)
	if err != nil {
		return err
	}

	// TODO: Headers and rows
	headers := []string{
		// Bond
		"bond_new_users_count", "bond_new_users_amount", "bond_old_users_count", "bond_old_users_amount",
		// Unbond
		"unbond_new_users_count", "unbond_new_users_amount", "unbond_old_users_count", "unbond_old_users_amount",
		// General
		"total_users_bond",
	}
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

func queryDailyReportBond(addressQuery string) (*DailyReportBond, error) {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}

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
				return nil, err
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
				return nil, err
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

	return &DailyReportBond{
		newUsersCount,
		newUsersAmount,
		oldUsersCount,
		oldUsersAmount,
		totalUsers,
	}, nil
}

func queryDailyReportUnbond(addressQuery string) (*DailyReportUnbond, error) {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}

	// query all the unbonds before last 24h with distinct on userAddressSender,
	// this will retrieve a list of all the current unbonders without repeated values
	_, rowsBefore, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryDailyReportUnbondBefore, addressQuery, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Creating a string/bool map for faster checking after " character removal (Numia is retrieving address sometimes embedded by double quotes)
	AddressBefore := make(map[string]bool)
	for _, address := range rowsBefore {
		tempAddr := strings.ReplaceAll(strings.Join(address, ""), "\"", "")
		AddressBefore[tempAddr] = true
	}

	// query all the unbonds after last 24h without distinct, then compare it with the previous result set
	// and compute who is new, who not and the related amounts
	_, rowsAfter, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryDailyReportUnbondAfter, addressQuery, true)
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

	for _, unbond := range rowsAfter {
		tempAddr := strings.ReplaceAll(unbond[0], "\"", "")
		// looking in AddressBefore to check if the user was already existing before 24h or not
		_, ok := AddressBefore[tempAddr]
		if ok {
			// found is an old user
			parseInt, err := strconv.ParseInt(unbond[1][:len(unbond[1])-2], 10, 64)
			if err != nil {
				return nil, err
			}
			// check if address is already counted
			if _, ok := AddressAfter[tempAddr]; !ok {
				oldUsersCount++
				AddressAfter[tempAddr] = true
			}
			oldUsersAmount += int(parseInt)
		} else {
			// not found is a new user
			parseInt, err := strconv.ParseInt(unbond[1][:len(unbond[1])-2], 10, 64)
			if err != nil {
				return nil, err
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

	return &DailyReportUnbond{
		newUsersCount,
		newUsersAmount,
		oldUsersCount,
		oldUsersAmount,
		totalUsers,
	}, nil
}

func queryDailyReportClaim(addressQuery string) (*DailyReportClaim, error) {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}

	// query all the claims before last 24h with distinct on userAddressSender,
	// this will retrieve a list of all the current claimers without repeated values
	_, rowsBefore, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryDailyReportClaimBefore, addressQuery, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Creating a string/bool map for faster checking after " character removal (Numia is retrieving address sometimes embedded by double quotes)
	AddressBefore := make(map[string]bool)
	for _, address := range rowsBefore {
		tempAddr := strings.ReplaceAll(strings.Join(address, ""), "\"", "")
		AddressBefore[tempAddr] = true
	}

	// query all the claims after last 24h without distinct, then compare it with the previous result set
	// and compute who is new, who not and the related amounts
	_, rowsAfter, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryDailyReportClaimAfter, addressQuery, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Variables to increment
	newUsersCount := 0
	oldUsersCount := 0

	// AddressAfter map to check duplicates in rowsAfter
	AddressAfter := make(map[string]bool)

	fmt.Println(rowsAfter)

	for _, claim := range rowsAfter {
		tempAddr := strings.ReplaceAll(claim[0], "\"", "")
		// looking in AddressBefore to check if the user was already existing before 24h or not
		_, ok := AddressBefore[tempAddr]
		if ok {
			// check if address is already counted
			if _, ok := AddressAfter[tempAddr]; !ok {
				oldUsersCount++
				AddressAfter[tempAddr] = true
			}
		} else {
			// check if address is already counted
			if _, ok := AddressAfter[tempAddr]; !ok {
				newUsersCount++
				AddressAfter[tempAddr] = true
			}
		}
	}

	// Variables to compute
	totalUsers := len(rowsBefore) + newUsersCount

	return &DailyReportClaim{
		newUsersCount,
		oldUsersCount,
		totalUsers,
	}, nil
}

func convertDenomToSharesAmount(denomAmount int) (int, error) {
	// TODO: conversion from denom bonded to vault shares amount

	return denomAmount, nil
}
