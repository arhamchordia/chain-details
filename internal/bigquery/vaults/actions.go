package vaults

import (
	"cloud.google.com/go/bigquery"
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	bigquerytypes "github.com/arhamchordia/chain-details/types/bigquery"
	"google.golang.org/api/iterator"
	"log"
)

// QueryBond returns a file with the bond events in all the blocks
func QueryBond(addressQuery string) error {
	addressFilterString := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependQueryVaultsBond
	if len(addressQuery) > 0 {
		addressFilterString = fmt.Sprintf(bigquerytypes.QueryVaultsBondAddressFilter, addressQuery)
		filename = fmt.Sprintf("%s_%s", filename, addressQuery)
	}

	rows, err := executeQueryAndFetchRows(bigquerytypes.QueryVaultsBond, addressFilterString)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = createCSV(filename, rows)
	if err != nil {
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

	rows, err := executeQueryAndFetchRows(bigquerytypes.QueryVaultsUnbond, addressFilterString)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = createCSV(filename, rows)
	if err != nil {
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

	rows, err := executeQueryAndFetchRows(bigquerytypes.QueryVaultsWithdraw, addressFilterString)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = createCSV(filename, rows)
	if err != nil {
		return err
	}

	return nil
}

func executeQueryAndFetchRows(query, addressFilter string) ([][]string, error) {
	bqQuerier, _ := internal.NewBigQueryQuerier()

	it, err := bqQuerier.ExecuteQuery(fmt.Sprintf(query, addressFilter))
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
