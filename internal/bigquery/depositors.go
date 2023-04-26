package bigquery

import (
	"cloud.google.com/go/bigquery"
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	"google.golang.org/api/iterator"
	"log"
)

// QueryDepositorsBond returns a file with the bond events in all the blocks given as startingHeight and endHeight
func QueryDepositorsBond() error {
	bqQuerier, _ := internal.NewBigQueryQuerier()

	it, err := bqQuerier.ExecuteQuery("SELECT * " +
		"FROM `numia-data.quasar.quasar_event_attributes` " +
		"WHERE `event_type` = 'message'" +
		"	AND `attribute_key` = 'action'" +
		"	AND `attribute_value` = '/cosmwasm.wasm.v1.MsgExecuteContract'" +
		"	AND EXISTS (" +
		"		SELECT 1" +
		"		FROM `numia-data.quasar.quasar_event_attributes` j1" +
		"		WHERE j1.`tx_id` = `numia-data.quasar.quasar_event_attributes`.`tx_id`" +
		"			AND j1.`block_height` = `numia-data.quasar.quasar_event_attributes`.`block_height`" +
		"			AND j1.`event_index` = `numia-data.quasar.quasar_event_attributes`.`event_index` + 1" +
		"			AND j1.`event_type` = 'message'" +
		"			AND j1.`attribute_value` = 'wasm'" +
		"		)" +
		"		AND EXISTS (" +
		"			SELECT 1" +
		"			FROM `numia-data.quasar.quasar_event_attributes` j2" +
		"			WHERE j2.`tx_id` = `numia-data.quasar.quasar_event_attributes`.`tx_id`" +
		"				AND j2.`block_height` = `numia-data.quasar.quasar_event_attributes`.`block_height`" +
		"				AND j2.`event_index` = `numia-data.quasar.quasar_event_attributes`.`event_index` + 2" +
		"				AND j2.`event_type` = 'coin_spent'" +
		"			)" +
		"			AND EXISTS (" +
		"				SELECT 1" +
		"				FROM `numia-data.quasar.quasar_event_attributes` j3" +
		"				WHERE j3.`tx_id` = `numia-data.quasar.quasar_event_attributes`.`tx_id`" +
		"					AND j3.`block_height` = `numia-data.quasar.quasar_event_attributes`.`block_height`" +
		"					AND j3.`event_index` = `numia-data.quasar.quasar_event_attributes`.`event_index` + 3" +
		"					AND j3.`event_type` = 'coin_received'" +
		"					AND j3.`attribute_value` = 'quasar18a2u6az6dzw528rptepfg6n49ak6hdzkf8ewf0n5r0nwju7gtdgqamr7qu'" +
		"			)" +
		"		ORDER BY block_height DESC")
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

	err = internal.WriteCSV("bigquery_depositors_bond", headerRow, rows)
	if err != nil {
		return err
	}

	return nil
}
