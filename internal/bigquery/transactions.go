package bigquery

import (
	"cloud.google.com/go/bigquery"
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	"github.com/arhamchordia/chain-details/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/api/iterator"
	"log"
)

func TransactionsQuery(AddressQuery string) error {
	bqQuerier, _ := internal.NewBigQueryQuerier()

	addr, err := sdk.AccAddressFromBech32(AddressQuery)
	if err != nil {
		return err
	}

	it, err := bqQuerier.ExecuteQuery("SELECT block_height, tx_id, message, ingestion_timestamp  " +
		"FROM `numia-data.quasar.quasar_tx_messages` " +
		"WHERE (" +
		"	SELECT COUNT(*)" +
		"	FROM UNNEST(REGEXP_EXTRACT_ALL(TO_JSON_STRING(message), r':\\s*\"([^\"]*)\"')) AS json_values" +
		"	WHERE json_values = '" + addr.String() + "'" +
		") > 0 " +
		"ORDER BY block_height ASC")
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

	err = internal.WriteCSV(types.PrefixBigQuery+"transactions_"+AddressQuery, headerRow, rows)
	if err != nil {
		return err
	}

	return nil
}
