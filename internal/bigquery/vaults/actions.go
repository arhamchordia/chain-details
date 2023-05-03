package vaults

import (
	"cloud.google.com/go/bigquery"
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	"github.com/arhamchordia/chain-details/types"
	"google.golang.org/api/iterator"
	"log"
)

// QueryDepositorsBond returns a file with the bond events in all the blocks given as startingHeight and endHeight
func QueryBond() error {
	bqQuerier, _ := internal.NewBigQueryQuerier()

	// TODO: This query is giving 2 rows per result, one with the sender and one with the amount/denom. Merge them.
	it, err := bqQuerier.ExecuteQuery("SELECT" +
		"		j1.block_height," +
		"		j1.tx_id," +
		"		j1.event_index," +
		"		j1.event_type," +
		"		j1.attribute_key," +
		"		j1.attribute_value," +
		"		j2.attribute_value AS sender," +
		"		j3.attribute_value AS coin_spent" +
		"	FROM `numia-data.quasar.quasar_event_attributes`" +
		"	JOIN `numia-data.quasar.quasar_event_attributes` j1" +
		"		ON j1.tx_id = `numia-data.quasar.quasar_event_attributes`.tx_id" +
		"		AND j1.block_height = `numia-data.quasar.quasar_event_attributes`.block_height" +
		"		AND j1.event_index = `numia-data.quasar.quasar_event_attributes`.event_index + 1" + // TODO: I don't like this + 1
		"		AND j1.event_type = 'message'" +
		"		AND j1.attribute_value = 'wasm'" +
		"	JOIN `numia-data.quasar.quasar_event_attributes` j2" +
		"		ON j2.tx_id = `numia-data.quasar.quasar_event_attributes`.tx_id" +
		"		AND j2.block_height = `numia-data.quasar.quasar_event_attributes`.block_height" +
		"		AND j2.event_index = `numia-data.quasar.quasar_event_attributes`.event_index + 2" + // TODO: I don't like this + 2
		"		AND j2.event_type = 'coin_spent'" +
		"		JOIN `numia-data.quasar.quasar_event_attributes` j3" +
		"	ON j3.tx_id = `numia-data.quasar.quasar_event_attributes`.tx_id" +
		"		AND j3.block_height = `numia-data.quasar.quasar_event_attributes`.block_height" +
		"		AND j3.event_index = `numia-data.quasar.quasar_event_attributes`.event_index + 3" + // TODO: I don't like this + 3
		"		AND j3.event_type = 'coin_received'" +
		"		AND j3.attribute_value = 'quasar18a2u6az6dzw528rptepfg6n49ak6hdzkf8ewf0n5r0nwju7gtdgqamr7qu'" + // TODO: Make this dynamic by cmd argument or flag
		"	WHERE `numia-data.quasar.quasar_event_attributes`.event_type = 'message'" +
		"		AND `numia-data.quasar.quasar_event_attributes`.attribute_key = 'action'" +
		"		AND `numia-data.quasar.quasar_event_attributes`.attribute_value = '/cosmwasm.wasm.v1.MsgExecuteContract'" +
		"	ORDER BY `numia-data.quasar.quasar_event_attributes`.block_height DESC")
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

	err = internal.WriteCSV(types.PrefixBigQuery+"vaults_bond", headerRow, rows)
	if err != nil {
		return err
	}

	return nil
}

func QueryUnbond() error {
	bqQuerier, _ := internal.NewBigQueryQuerier()

	// TODO query
	it, err := bqQuerier.ExecuteQuery("")
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

	err = internal.WriteCSV(types.PrefixBigQuery+"vaults_unbond", headerRow, rows)
	if err != nil {
		return err
	}

	return nil
}

func QueryWithdraw() error {
	bqQuerier, _ := internal.NewBigQueryQuerier()

	// TODO query
	it, err := bqQuerier.ExecuteQuery("")
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

	err = internal.WriteCSV(types.PrefixBigQuery+"vaults_withdraw", headerRow, rows)
	if err != nil {
		return err
	}

	return nil
}
