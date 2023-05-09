package vaults

import (
	"cloud.google.com/go/bigquery"
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	bigquerytypes "github.com/arhamchordia/chain-details/types/bigquery"
	"google.golang.org/api/iterator"
	"log"
	"regexp"
)

// QueryBond returns a file with the bond events in all the blocks
func QueryBond(addressQuery string, confirmedQuery bool) error {
	addressFilterString := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependQueryVaultsBond
	if len(addressQuery) > 0 {
		addressFilterString = fmt.Sprintf(bigquerytypes.QueryVaultsBondAddressFilter, addressQuery)
		filename = fmt.Sprintf("%s_%s", filename, addressQuery)
	}
	if confirmedQuery {
		filename = fmt.Sprintf("%s_%s", filename, "confirmed")
	}

	rows, err := executeQueryAndFetchRows(bigquerytypes.QueryVaultsBond, addressFilterString, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if confirmedQuery {
		confirmedRows, err := executeQueryAndFetchRows(bigquerytypes.QueryVaultsBondConfirmedFilter, "", false)
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
			if len(match) > 1 {
				bondID := match[1]
				if _, exists := confirmedBondIDs[bondID]; exists {
					filteredRows = append(filteredRows, row)
				}
			}
		}
		rows = filteredRows
	}

	err = createCSV(filename, rows)
	if err != nil {
		log.Printf("Warning: %v", err)
		return err
	}

	return nil
}

// QueryUnbond returns a file with the unbond events in all the blocks
func QueryUnbond(addressQuery string) error {
	addressFilterString := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependQueryVaultsUnbond
	if len(addressQuery) > 0 {
		addressFilterString = fmt.Sprintf(bigquerytypes.QueryVaultsUnbondAddressFilter, addressQuery)
		filename = fmt.Sprintf("%s_%s", filename, addressQuery)
	}

	rows, err := executeQueryAndFetchRows(bigquerytypes.QueryVaultsUnbond, addressFilterString, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = createCSV(filename, rows)
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

	rows, err := executeQueryAndFetchRows(bigquerytypes.QueryVaultsWithdraw, addressFilterString, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = createCSV(filename, rows)
	if err != nil {
		log.Printf("Warning: %v", err)
		return err
	}

	return nil
}

func executeQueryAndFetchRows(query, addressFilter string, applyAddressFilter bool) ([][]string, error) {
	bqQuerier, _ := internal.NewBigQueryQuerier()

	// apply addressFilter only if applyAddressFilter is true
	if applyAddressFilter {
		query = fmt.Sprintf(query, addressFilter)
	}

	it, err := bqQuerier.ExecuteQuery(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to execute BigQuery query: %v", err)
	}
	defer bqQuerier.Close()

	var rows [][]string

	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Failed to iterate results: %v", err)
		}

		var csvRow []string
		for _, val := range row {
			csvRow = append(csvRow, fmt.Sprintf("%v", val))
		}
		rows = append(rows, csvRow)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("No rows returned by query")
	}

	return rows, nil
}

func createCSV(filename string, rows [][]string) error {
	if len(rows) == 0 {
		return fmt.Errorf("No rows to write to the CSV")
	}

	headerRow := make([]string, len(rows[0]))
	for i := range rows[0] {
		headerRow[i] = fmt.Sprintf("column_%d", i+1)
	}

	err := internal.WriteCSV(filename, headerRow, rows)
	if err != nil {
		return err
	}

	return nil
}
