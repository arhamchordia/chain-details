package vaults

import (
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	bigquerytypes "github.com/arhamchordia/chain-details/types/bigquery"
	"log"
	"regexp"
)

// QueryBond returns a file with the bond events in all the blocks
func QueryBond(addressQuery string, confirmedQuery bool, pendingQuery bool) error {
	if confirmedQuery && pendingQuery {
		return fmt.Errorf("--confirmed and --pending flags cannot be used together")
	}

	addressFilterString := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependQueryVaultsBond
	if len(addressQuery) > 0 {
		addressFilterString = fmt.Sprintf(bigquerytypes.QueryVaultsBondAddressFilter, addressQuery)
		filename = fmt.Sprintf("%s_%s", filename, addressQuery)
	}
	if confirmedQuery {
		filename = fmt.Sprintf("%s_%s", filename, "confirmed")
	}
	if pendingQuery {
		filename = fmt.Sprintf("%s_%s", filename, "pending")
	}

	headers, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryVaultsBond, addressFilterString, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if confirmedQuery || pendingQuery {
		_, confirmedRows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryVaultsBondConfirmedFilter, "", false)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// creating map storing the bond_id from the second query
		confirmedBondIDs := make(map[string]bool)
		for _, row := range confirmedRows {
			bondID := row[0]
			confirmedBondIDs[bondID] = true
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
					if _, exists := confirmedBondIDs[bondID]; exists {
						filteredRows = append(filteredRows, row)
					}
				} else if pendingQuery {
					if _, exists := confirmedBondIDs[bondID]; !exists {
						filteredRows = append(filteredRows, row)
					}
				}
			}
		}
		rows = filteredRows
	}

	err = internal.WriteCSV(filename, headers, rows)
	if err != nil {
		log.Printf("Warning: %v", err)
		return err
	}

	return nil
}

// QueryUnbond returns a file with the unbond events in all the blocks
func QueryUnbond(addressQuery string, confirmedQuery bool, pendingQuery bool) error {
	addressFilterString := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependQueryVaultsUnbond
	if len(addressQuery) > 0 {
		addressFilterString = fmt.Sprintf(bigquerytypes.QueryVaultsUnbondAddressFilter, addressQuery)
		filename = fmt.Sprintf("%s_%s", filename, addressQuery)
	}
	if confirmedQuery {
		filename = fmt.Sprintf("%s_%s", filename, "confirmed")
	}
	if pendingQuery {
		filename = fmt.Sprintf("%s_%s", filename, "pending")
	}

	headers, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryVaultsUnbond, addressFilterString, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if confirmedQuery || pendingQuery {
		_, confirmedRows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryVaultsUnbondConfirmedFilter, "", false)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// creating map storing the bond_id from the second query
		confirmedBondIDs := make(map[string]bool)
		for _, row := range confirmedRows {
			bondID := row[0]
			confirmedBondIDs[bondID] = true
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
					if _, exists := confirmedBondIDs[bondID]; exists {
						filteredRows = append(filteredRows, row)
					}
				} else if pendingQuery {
					if _, exists := confirmedBondIDs[bondID]; !exists {
						filteredRows = append(filteredRows, row)
					}
				}
			}
		}
		rows = filteredRows
	}

	err = internal.WriteCSV(filename, headers, rows)
	if err != nil {
		log.Printf("Warning: %v", err)
		return err
	}

	return nil
}

// QueryWithdraw returns a file with the withdraw events in all the blocks
func QueryWithdraw(addressQuery string) error {
	addressFilterString := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependQueryVaultsWithdraw
	if len(addressQuery) > 0 {
		addressFilterString = fmt.Sprintf(bigquerytypes.QueryVaultsWithdrawAddressFilter, addressQuery)
		filename = fmt.Sprintf("%s_%s", filename, addressQuery)
	}

	headers, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryVaultsWithdraw, addressFilterString, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = internal.WriteCSV(filename, headers, rows)
	if err != nil {
		log.Printf("Warning: %v", err)
		return err
	}

	return nil
}
