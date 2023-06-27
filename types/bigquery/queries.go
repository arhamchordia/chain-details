package bigquery

const (
	// QueryTransactions bigquery transactions --address
	QueryTransactions = "SELECT block_height, tx_id, message, ingestion_timestamp  " +
		"FROM `numia-data.quasar.quasar_tx_messages` " +
		"WHERE (" +
		"	SELECT COUNT(*)" +
		"	FROM UNNEST(REGEXP_EXTRACT_ALL(TO_JSON_STRING(message), r':\\s*\"([^\"]*)\"')) AS json_values" +
		"	WHERE json_values = '%s'" +
		") > 0 " +
		"ORDER BY block_height ASC"

	// QueryVaultsBond bigquery bond (main query)
	QueryVaultsBond = `WITH combined_rows AS (
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
    ` + "`numia-data.quasar.quasar_event_attributes`" + `
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
  ARRAY_AGG(STRUCT(attribute_key, attribute_value)) AS attribute_pairs,
  MAX(latest_ingestion_timestamp) AS latest_ingestion_timestamp
FROM
  key_value_pairs
GROUP BY
  block_height,
  tx_id
ORDER BY
  block_height ASC;`

	// QueryVaultsBondAddressFilter bigquery bond optional query for flag --address
	QueryVaultsBondAddressFilter = "AND EXISTS (SELECT 1 FROM combined_rows c WHERE c.tx_id = combined_rows.tx_id AND c.attribute_key = 'spender' AND c.attribute_value = '%s')"

	// QueryVaultsBondResponseFilter bigquery bond optional query for --confirmed and --pending flags
	QueryVaultsBondResponseFilter = `
WITH extracted_data AS (
  SELECT
    block_height,
    tx_id,
    REGEXP_EXTRACT(attribute_value, r'bond_id: "([^"]+)"') AS bond_id,
    REGEXP_EXTRACT(attribute_value, r'share_amount: Uint128\((\d+)\)') AS share_amount,
    REGEXP_EXTRACT(attribute_value, r'owner: Addr\("([^"]+)"\)') AS owner_addr,
    CAST(ingestion_timestamp AS STRING) AS ingestion_timestamp
  FROM
    ` + "`numia-data.quasar.quasar_message_event_attributes`" + `
  WHERE
    event_type = 'wasm'
    AND attribute_key = 'callback-info'
    AND attribute_value LIKE '%BondResponse%'
    AND attribute_value LIKE '%bond_id%'
)
SELECT
  bond_id,
  STRING_AGG(share_amount, ', ') AS share_amounts,
  STRING_AGG(owner_addr, ', ') AS owner_addrs,
  STRING_AGG(ingestion_timestamp, ', ') AS ingestion_timestamps,
  STRING_AGG(CAST(block_height AS STRING), ', ') AS block_heights,
  STRING_AGG(tx_id, ', ') AS tx_ids
FROM
  extracted_data
GROUP BY
  bond_id
ORDER BY
  bond_id ASC;
`
	QueryVaultsBondShareAmountsTxIds = `
WITH 
message_actions AS (
    SELECT DISTINCT G.tx_id 
    FROM ` + "`numia-data.quasar.quasar_event_attributes`" + ` AS G
    WHERE G.event_type = 'message' 
      AND G.attribute_key = 'action' 
      AND G.attribute_value = '/cosmwasm.wasm.v1.MsgMigrateContract'
      AND EXISTS (
        SELECT 1 
        FROM ` + "`numia-data.quasar.quasar_event_attributes`" + ` AS H
        WHERE H.tx_id = G.tx_id AND H.attribute_key = '_contract_address'
      )
      AND EXISTS (
        SELECT 1 
        FROM ` + "`numia-data.quasar.quasar_event_attributes`" + ` AS I
        WHERE I.tx_id = G.tx_id AND I.attribute_key = 'user'
      )
      AND EXISTS (
        SELECT 1 
        FROM ` + "`numia-data.quasar.quasar_event_attributes`" + ` AS J
        WHERE J.tx_id = G.tx_id AND J.attribute_key = 'vault_token_balance'
      )
)
SELECT A.tx_id
FROM message_actions AS A
`

	// QueryVaultsUnbond bigquery unbond (main query)
	QueryVaultsUnbond = `
WITH combined_rows AS (
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
    ` + "`numia-data.quasar.quasar_event_attributes`" + `
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
  WHERE attribute_key IN ('spender', 'burnt', 'bond_id')
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
  ARRAY_AGG(STRUCT(attribute_key, attribute_value)) AS attribute_pairs,
  MAX(latest_ingestion_timestamp) AS latest_ingestion_timestamp
FROM
  key_value_pairs
GROUP BY
  block_height,
  tx_id
ORDER BY
  block_height ASC;`

	// QueryVaultsUnbondAddressFilter bigquery unbond optional query for flag --address
	QueryVaultsUnbondAddressFilter = "AND EXISTS (SELECT 1 FROM combined_rows c WHERE c.tx_id = combined_rows.tx_id AND c.attribute_key = 'spender' AND c.attribute_value = '%s')"

	// QueryVaultsUnbondConfirmedFilter bigquery unbond optional query for --confirmed and --pending flags
	QueryVaultsUnbondConfirmedFilter = `
WITH extracted_data AS (
  SELECT
    block_height,
    tx_id,
    REGEXP_EXTRACT(attribute_value, r'unbond_id: "([^"]+)"') AS unbond_id,
    REGEXP_EXTRACT(attribute_value, r'amount: Some\(Uint128\((\d+)\)\)') AS share_amount,
    REGEXP_EXTRACT(attribute_value, r'owner: Addr\("([^"]+)"\)') AS owner_addr,
    CAST(ingestion_timestamp AS STRING) AS ingestion_timestamp
  FROM
    ` + "`numia-data.quasar.quasar_event_attributes`" + `
  WHERE
    event_type = 'wasm'
    AND attribute_key = 'callback-info'
    AND attribute_value LIKE '%UnbondResponse%'
    AND attribute_value LIKE '%unbond_id%'
)
SELECT
  unbond_id,
  STRING_AGG(share_amount, ', ') AS share_amounts,
  STRING_AGG(owner_addr, ', ') AS owner_addrs,
  STRING_AGG(ingestion_timestamp, ', ') AS ingestion_timestamps,
  STRING_AGG(CAST(block_height AS STRING), ', ') AS block_heights,
  STRING_AGG(tx_id, ', ') AS tx_ids
FROM
  extracted_data
GROUP BY
  unbond_id
ORDER BY
  unbond_id ASC;
`

	// QueryVaultsClaim bigquery claim (main query)
	QueryVaultsClaim = `WITH combined_rows AS (
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
    ` + "`numia-data.quasar.quasar_event_attributes`" + `
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
  WHERE attribute_key IN ('spender', 'amount', 'packet_sequence', 'packet_src_channel')
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
  ARRAY_AGG(STRUCT(attribute_key, attribute_value)) AS attribute_pairs,
  MAX(latest_ingestion_timestamp) AS latest_ingestion_timestamp
FROM
  key_value_pairs
GROUP BY
  block_height,
  tx_id
ORDER BY
  block_height ASC;`

	// QueryVaultsClaimAddressFilter bigquery claim optional query for flag --address
	QueryVaultsClaimAddressFilter = "WHERE EXISTS (SELECT 1 FROM combined_rows c WHERE c.tx_id = combined_rows.tx_id AND c.attribute_key = 'spender' AND c.attribute_value = '%s')"

	// QueryDailyReportRewardsUpdateUser retrieve all the rewards contracts action:update_user_index SubMsgResponses
	QueryDailyReportRewardsUpdateUser = `WITH combined_rows AS (
  SELECT
    block_height,
    tx_id,
    attribute_key,
    attribute_value,
    MAX(ingestion_timestamp) OVER (PARTITION BY tx_id, attribute_key) AS latest_ingestion_timestamp
  FROM
    ` + "`numia-data.quasar.quasar_event_attributes`" + `
  WHERE
    event_type = 'wasm'
),

grouped_data AS (
  SELECT
    block_height,
    tx_id,
    ARRAY_AGG(STRUCT(attribute_key, attribute_value)) AS attributes,
    MAX(latest_ingestion_timestamp) AS latest_ingestion_timestamp
  FROM
    combined_rows
  WHERE
    tx_id IN (
      SELECT tx_id
      FROM combined_rows
      WHERE (attribute_key = 'action' AND attribute_value = 'update_user_index')
        AND tx_id IN (
          SELECT tx_id 
          FROM combined_rows 
          WHERE attribute_key = '_contract_address' AND attribute_value = '%s'
        )
    )
  GROUP BY
    block_height,
    tx_id
),

flattened_data AS (
  SELECT
    block_height,
    tx_id,
    attr.attribute_value AS user,
    (
      SELECT attribute_value 
      FROM UNNEST(attributes) attr
      WHERE attr.attribute_key = 'vault_token_balance'
      LIMIT 1
    ) AS vault_token_balance,
    latest_ingestion_timestamp
  FROM
    grouped_data,
    UNNEST(attributes) attr
  WHERE
    attr.attribute_key = 'user'
),

distinct_flattened_data AS (
  SELECT DISTINCT
    block_height,
    user,
    vault_token_balance,
    latest_ingestion_timestamp
  FROM
    flattened_data
)

SELECT
  user,
  ARRAY_AGG(
    STRUCT(block_height, vault_token_balance, latest_ingestion_timestamp) 
    ORDER BY block_height ASC
  ) AS user_transactions
FROM
  distinct_flattened_data
GROUP BY
  user
ORDER BY
  user ASC;`
)
