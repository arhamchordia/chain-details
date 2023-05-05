package vaults

import (
	"cloud.google.com/go/bigquery"
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	"github.com/arhamchordia/chain-details/types"
	"google.golang.org/api/iterator"
	"log"
)

// QueryBond returns a file with the bond events in all the blocks
func QueryBond(AddressQuery string) error {
	bqQuerier, _ := internal.NewBigQueryQuerier()

	addressFilter := ""
	filename := types.PrefixBigQuery + "vaults_bond"
	if len(AddressQuery) > 0 {
		addressFilter = fmt.Sprintf("AND EXISTS (SELECT 1 FROM combined_rows c WHERE c.tx_id = combined_rows.tx_id AND c.attribute_key = 'spender' AND c.attribute_value = '%s')", AddressQuery)
		filename = fmt.Sprintf("%s_%s", filename, AddressQuery)
	}

	it, err := bqQuerier.ExecuteQuery(fmt.Sprintf(`WITH combined_rows AS (
  SELECT
    block_height,
    tx_id,
    event_type,
    event_source,
    attribute_key,
    attribute_value,
    ingestion_timestamp,
    ROW_NUMBER() OVER (PARTITION BY tx_id, attribute_key ORDER BY ingestion_timestamp) AS row_num
  FROM
    `+"`numia-data.quasar.quasar_event_attributes`"+`
  WHERE
    event_type = 'wasm' OR event_type = 'coin_spent'
  ORDER BY
    block_height ASC
),
filtered_tx_ids AS (
  SELECT DISTINCT
    tx_id
  FROM
    combined_rows
  WHERE
    attribute_key = 'bond_id'
    %s
),
valid_tx_ids AS (
  SELECT DISTINCT
    tx_id
  FROM
    combined_rows
  WHERE
    tx_id NOT IN (
      SELECT tx_id
      FROM combined_rows
      WHERE attribute_key = 'action' AND attribute_value = 'start_unbond'
    )
),
filtered_combined_rows AS (
  SELECT *
  FROM combined_rows
  WHERE attribute_key IN ('spender', 'amount', 'bond_id', 'deposit')
),
key_value_pairs AS (
  SELECT
    fcr.block_height,
    fcr.tx_id,
    STRING_AGG(DISTINCT fcr.event_type) AS event_types,
    STRING_AGG(DISTINCT fcr.event_source) AS event_sources,
    fcr.attribute_key,
    fcr.attribute_value,
    MAX(fcr.ingestion_timestamp) AS latest_ingestion_timestamp
  FROM
    filtered_combined_rows fcr
  JOIN
    filtered_tx_ids ft
  ON
    fcr.tx_id = ft.tx_id
  JOIN
    valid_tx_ids vt
  ON
    fcr.tx_id = vt.tx_id
  GROUP BY
    fcr.block_height,
    fcr.tx_id,
    fcr.attribute_key,
    fcr.attribute_value
)
SELECT
  block_height,
  tx_id,
  MAX(event_types) AS event_types,
  MAX(event_sources) AS event_sources,
  ARRAY_AGG(STRUCT(attribute_key, attribute_value)) AS attribute_pairs,
  MAX(latest_ingestion_timestamp) AS latest_ingestion_timestamp
FROM
  key_value_pairs
GROUP BY
  block_height,
  tx_id
ORDER BY
  block_height ASC;`, addressFilter))
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
	filename := types.PrefixBigQuery + "vaults_unbond"
	if len(AddressQuery) > 0 {
		addressFilter = fmt.Sprintf("AND EXISTS (SELECT 1 FROM combined_rows c WHERE c.tx_id = combined_rows.tx_id AND c.attribute_key = 'spender' AND c.attribute_value = '%s')", AddressQuery)
		filename = fmt.Sprintf("%s_%s", filename, AddressQuery)
	}

	it, err := bqQuerier.ExecuteQuery(fmt.Sprintf(`WITH combined_rows AS (
  SELECT
    block_height,
    tx_id,
    event_type,
    event_source,
    attribute_key,
    attribute_value,
    ingestion_timestamp,
    ROW_NUMBER() OVER (PARTITION BY tx_id, attribute_key ORDER BY ingestion_timestamp) AS row_num
  FROM
    `+"`numia-data.quasar.quasar_event_attributes`"+`
  WHERE
    event_type = 'wasm' OR event_type = 'coin_spent'
  ORDER BY
    block_height ASC
),
filtered_tx_ids AS (
  SELECT DISTINCT
    tx_id
  FROM
    combined_rows
  WHERE
    attribute_key = 'bond_id'
    %s
),
valid_tx_ids AS (
  SELECT DISTINCT
    tx_id
  FROM
    combined_rows
  WHERE
    tx_id IN (
      SELECT tx_id
      FROM combined_rows
      WHERE attribute_key = 'action' AND attribute_value = 'start_unbond'
    )
),
filtered_combined_rows AS (
  SELECT *
  FROM combined_rows
  WHERE attribute_key IN ('spender', 'amount', 'bond_id')
),
key_value_pairs AS (
  SELECT
    fcr.block_height,
    fcr.tx_id,
    STRING_AGG(DISTINCT fcr.event_type) AS event_types,
    STRING_AGG(DISTINCT fcr.event_source) AS event_sources,
    fcr.attribute_key,
    fcr.attribute_value,
    MAX(fcr.ingestion_timestamp) AS latest_ingestion_timestamp
  FROM
    filtered_combined_rows fcr
  JOIN
    filtered_tx_ids ft
  ON
    fcr.tx_id = ft.tx_id
  JOIN
    valid_tx_ids vt
  ON
    fcr.tx_id = vt.tx_id
  GROUP BY
    fcr.block_height,
    fcr.tx_id,
    fcr.attribute_key,
    fcr.attribute_value
)
SELECT
  block_height,
  tx_id,
  MAX(event_types) AS event_types,
  MAX(event_sources) AS event_sources,
  ARRAY_AGG(STRUCT(attribute_key, attribute_value)) AS attribute_pairs,
  MAX(latest_ingestion_timestamp) AS latest_ingestion_timestamp
FROM
  key_value_pairs
GROUP BY
  block_height,
  tx_id
ORDER BY
  block_height ASC;`, addressFilter))
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
	filename := types.PrefixBigQuery + "vaults_withdraw"
	if len(AddressQuery) > 0 {
		addressFilter = fmt.Sprintf("WHERE EXISTS (SELECT 1 FROM combined_rows c WHERE c.tx_id = combined_rows.tx_id AND c.attribute_key = 'spender' AND c.attribute_value = '%s')", AddressQuery)
		filename = fmt.Sprintf("%s_%s", filename, AddressQuery)
	}

	it, err := bqQuerier.ExecuteQuery(fmt.Sprintf(`WITH combined_rows AS (
  SELECT
    block_height,
    tx_id,
    event_type,
    event_source,
    attribute_key,
    attribute_value,
    ingestion_timestamp,
    ROW_NUMBER() OVER (PARTITION BY tx_id, attribute_key ORDER BY ingestion_timestamp) AS row_num
  FROM
    `+"`numia-data.quasar.quasar_event_attributes`"+`
  WHERE
    event_type = 'wasm' OR event_type = 'coin_spent' OR event_type = 'send_packet'
  ORDER BY
    block_height ASC
),
filtered_tx_ids AS (
  SELECT DISTINCT
    tx_id
  FROM
    combined_rows
  %s
),
valid_tx_ids AS (
  SELECT DISTINCT
    tx_id
  FROM
    combined_rows
  WHERE
    tx_id IN (
      SELECT tx_id
      FROM combined_rows
      WHERE attribute_key = 'action' AND attribute_value = 'unbond'
    )
),
filtered_combined_rows AS (
  SELECT *
  FROM combined_rows
  WHERE attribute_key IN ('spender', 'amount', 'packet_sequence', 'pending-msg')
),
key_value_pairs AS (
  SELECT
    fcr.block_height,
    fcr.tx_id,
    STRING_AGG(DISTINCT fcr.event_type) AS event_types,
    STRING_AGG(DISTINCT fcr.event_source) AS event_sources,
    fcr.attribute_key,
    fcr.attribute_value,
    MAX(fcr.ingestion_timestamp) AS latest_ingestion_timestamp
  FROM
    filtered_combined_rows fcr
  JOIN
    filtered_tx_ids ft
  ON
    fcr.tx_id = ft.tx_id
  JOIN
    valid_tx_ids vt
  ON
    fcr.tx_id = vt.tx_id
  GROUP BY
    fcr.block_height,
    fcr.tx_id,
    fcr.attribute_key,
    fcr.attribute_value
)
SELECT
  block_height,
  tx_id,
  MAX(event_types) AS event_types,
  MAX(event_sources) AS event_sources,
  ARRAY_AGG(STRUCT(attribute_key, attribute_value)) AS attribute_pairs,
  MAX(latest_ingestion_timestamp) AS latest_ingestion_timestamp
FROM
  key_value_pairs
GROUP BY
  block_height,
  tx_id
ORDER BY
  block_height ASC;`, addressFilter))
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
