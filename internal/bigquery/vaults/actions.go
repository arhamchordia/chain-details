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
func QueryBond(AddressQuery string) error {
	bqQuerier, _ := internal.NewBigQueryQuerier()

	addressFilter := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependQueryVaultsBond
	if len(AddressQuery) > 0 {
		addressFilter = fmt.Sprintf(bigquerytypes.QueryVaultsBondAddressFilter, AddressQuery)
		filename = fmt.Sprintf("%s_%s", filename, AddressQuery)
	}

	it, err := bqQuerier.ExecuteQuery(fmt.Sprintf(bigquerytypes.QueryVaultsBond, addressFilter))
	if err != nil {
		log.Fatalf("Failed to execute BigQuery query: %v", err)
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
			log.Fatalf("Failed to iterate results: %v", err)
		}

		var csvRow []string
		for _, val := range row {
			csvRow = append(csvRow, fmt.Sprintf("%v", val))
		}
		rows = append(rows, csvRow)
	}

	if len(rows) == 0 {
		log.Fatalf("No rows returned by query")
	}

	headerRow := make([]string, len(rows[0]))
	for i := range rows[0] {
		headerRow[i] = fmt.Sprintf("column_%d", i+1)
	}

	err = internal.WriteCSV(filename, headerRow, rows)
	if err != nil {
		return err
	}

	return nil
}

// QueryUnbond returns a file with the unbond events in all the blocks
func QueryUnbond(AddressQuery string) error {
	bqQuerier, _ := internal.NewBigQueryQuerier()

	addressFilter := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependQueryVaultsUnbond
	if len(AddressQuery) > 0 {
		addressFilter = fmt.Sprintf(bigquerytypes.QueryVaultsUnbondAddressFilter, AddressQuery)
		filename = fmt.Sprintf("%s_%s", filename, AddressQuery)
	}

	it, err := bqQuerier.ExecuteQuery(fmt.Sprintf(bigquerytypes.QueryVaultsUnbond, addressFilter))
	if err != nil {
		log.Fatalf("Failed to execute BigQuery query: %v", err)
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
			log.Fatalf("Failed to iterate results: %v", err)
		}

		var csvRow []string
		for _, val := range row {
			csvRow = append(csvRow, fmt.Sprintf("%v", val))
		}
		rows = append(rows, csvRow)
	}

	if len(rows) == 0 {
		log.Fatalf("No rows returned by query")
	}

	headerRow := make([]string, len(rows[0]))
	for i := range rows[0] {
		headerRow[i] = fmt.Sprintf("column_%d", i+1)
	}

	err = internal.WriteCSV(filename, headerRow, rows)
	if err != nil {
		return err
	}

	return nil
}

func QueryWithdraw(AddressQuery string) error {
	bqQuerier, _ := internal.NewBigQueryQuerier()

	addressFilter := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependQueryVaultsWithdraw
	if len(AddressQuery) > 0 {
		addressFilter = fmt.Sprintf(bigquerytypes.QueryVaultsWithdrawAddressFilter, AddressQuery)
		filename = fmt.Sprintf("%s_%s", filename, AddressQuery)
	}

	it, err := bqQuerier.ExecuteQuery(fmt.Sprintf(bigquerytypes.QueryVaultsWithdraw, addressFilter))
	if err != nil {
		log.Fatalf("Failed to execute BigQuery query: %v", err)
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
			log.Fatalf("Failed to iterate results: %v", err)
		}

		var csvRow []string
		for _, val := range row {
			csvRow = append(csvRow, fmt.Sprintf("%v", val))
		}
		rows = append(rows, csvRow)
	}

	if len(rows) == 0 {
		log.Fatalf("No rows returned by query")
	}

	headerRow := make([]string, len(rows[0]))
	for i := range rows[0] {
		headerRow[i] = fmt.Sprintf("column_%d", i+1)
	}

	err = internal.WriteCSV(filename, headerRow, rows)
	if err != nil {
		return err
	}

	return nil
}
