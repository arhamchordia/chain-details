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

	// QueryLPVaultsBond bigquery bond (main query)
	QueryLPVaultsBond = `WITH combined_rows AS (
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

	// QueryLPVaultsBondAddressFilter bigquery bond optional query for flag --address
	QueryLPVaultsBondAddressFilter = "AND EXISTS (SELECT 1 FROM combined_rows c WHERE c.tx_id = combined_rows.tx_id AND c.attribute_key = 'spender' AND c.attribute_value = '%s')"

	// QueryLPVaultsBondResponseFilter bigquery bond optional query for --confirmed and --pending flags
	QueryLPVaultsBondResponseFilter = `
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
	// QueryLPVaultsBondShareAmountsTxIds
	QueryLPVaultsBondShareAmountsTxIds = `
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

	// QueryLPVaultsUnbond bigquery unbond (main query)
	QueryLPVaultsUnbond = `
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

	// QueryLPVaultsUnbondAddressFilter bigquery unbond optional query for flag --address
	QueryLPVaultsUnbondAddressFilter = "AND EXISTS (SELECT 1 FROM combined_rows c WHERE c.tx_id = combined_rows.tx_id AND c.attribute_key = 'spender' AND c.attribute_value = '%s')"

	// QueryLPVaultsUnbondConfirmedFilter bigquery unbond optional query for --confirmed and --pending flags
	QueryLPVaultsUnbondConfirmedFilter = `
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

	// QueryLPVaultsClaim bigquery claim (main query)
	QueryLPVaultsClaim = `WITH combined_rows AS (
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

	// QueryLPVaultsClaimAddressFilter bigquery claim optional query for flag --address
	QueryLPVaultsClaimAddressFilter = "WHERE EXISTS (SELECT 1 FROM combined_rows c WHERE c.tx_id = combined_rows.tx_id AND c.attribute_key = 'spender' AND c.attribute_value = '%s')"

	// QueryLPReportRewardsUpdateUser retrieve all the rewards contracts action:update_user_index SubMsgResponses
	QueryLPReportRewardsUpdateUser = `WITH combined_rows AS (
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
	AND block_height >= %d
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

	QueryCLVaultsDeposit = `
WITH filtered_events AS (
  SELECT 
    tx_id,
    MAX(ingestion_timestamp) AS ingestion_timestamp,
    MAX(CASE WHEN attribute_key='receiver' THEN attribute_value END) AS sender,
    MAX(CASE WHEN attribute_key='_contract_address' AND attribute_value='%s' THEN TRUE ELSE FALSE END) AS is_target_contract,
    MAX(CASE WHEN attribute_key='action' AND attribute_value='exact_deposit' THEN TRUE ELSE FALSE END) AS is_exact_deposit,
    MAX(CASE WHEN attribute_key='amount0' THEN CAST(attribute_value AS INT64) ELSE 0 END) AS amount0,
    MAX(CASE WHEN attribute_key='amount1' THEN CAST(attribute_value AS INT64) ELSE 0 END) AS amount1,
    MAX(CASE WHEN attribute_key='refund0_amount' THEN CAST(attribute_value AS INT64) ELSE 0 END) AS refund0_amount,
    MAX(CASE WHEN attribute_key='refund1_amount' THEN CAST(attribute_value AS INT64) ELSE 0 END) AS refund1_amount
  FROM
    ` + "`numia-data.osmosis.osmosis_message_event_attributes`" + `
  WHERE event_type = 'wasm'
  GROUP BY tx_id
)
SELECT 
  tx_id,
  ingestion_timestamp,
  sender,
  IFNULL(amount0 - refund0_amount, 0) AS final_amount0,
  IFNULL(amount1 - refund1_amount, 0) AS final_amount1
FROM filtered_events
WHERE is_target_contract AND is_exact_deposit
ORDER BY ingestion_timestamp DESC;
`
	QueryCLVaultsWithdraw = `
WITH 
filtered_events AS (
  SELECT 
    tx_id,
    MAX(ingestion_timestamp) AS ingestion_timestamp,
    MAX(CASE WHEN attribute_key='_contract_address' AND attribute_value='%s' THEN TRUE ELSE FALSE END) AS is_target_contract,
    MAX(CASE WHEN attribute_key='action' AND attribute_value='withdraw' THEN TRUE ELSE FALSE END) AS is_withdraw,
    MAX(CASE WHEN attribute_key='token0_amount' THEN CAST(attribute_value AS INT64) ELSE 0 END) AS token0_amount,
    MAX(CASE WHEN attribute_key='token1_amount' THEN CAST(attribute_value AS INT64) ELSE 0 END) AS token1_amount,
    MAX(CASE WHEN attribute_key='liquidity_amount' THEN CAST(attribute_value AS INT64) ELSE 0 END) AS liquidity_amount,
    MAX(CASE WHEN attribute_key='share_amount' THEN CAST(attribute_value AS INT64) ELSE 0 END) AS share_amount
  FROM
    ` + "`numia-data.osmosis.osmosis_message_event_attributes`" + `
  WHERE event_type = 'wasm'
  GROUP BY tx_id
),
sender_info AS (
  SELECT
    tx_id,
    MAX(ingestion_timestamp) AS ingestion_timestamp,
    MAX(CASE WHEN attribute_key='sender' THEN attribute_value END) AS sender
  FROM
    ` + "`numia-data.osmosis.osmosis_message_event_attributes`" + `
  WHERE event_type = 'message'
  GROUP BY tx_id
)
SELECT
  f.tx_id,
  f.ingestion_timestamp,
  s.sender,
  f.token0_amount,
  f.token1_amount
FROM filtered_events f
JOIN sender_info s ON f.tx_id = s.tx_id
WHERE f.is_target_contract AND f.is_withdraw
ORDER BY f.ingestion_timestamp DESC;
`
	QueryCLVaultsClaim = `
WITH 
claim_rewards_events AS (
  SELECT 
    tx_id,
    MAX(ingestion_timestamp) AS ingestion_timestamp,
    MAX(CASE WHEN attribute_key='_contract_address' AND attribute_value='%s' THEN TRUE ELSE FALSE END) AS is_target_contract,
    MAX(CASE WHEN attribute_key='action' AND attribute_value='claim_user_rewards' THEN TRUE ELSE FALSE END) AS is_claim_rewards,
    MAX(CASE WHEN attribute_key='recipient' THEN attribute_value END) AS recipient
  FROM
    ` + "`numia-data.osmosis.osmosis_message_event_attributes`" + `
  WHERE event_type = 'wasm'
  GROUP BY tx_id
),
transfer_amounts AS (
  SELECT
    tx_id,
    MAX(ingestion_timestamp) AS ingestion_timestamp,
    SPLIT(MAX(CASE WHEN attribute_key='amount' THEN attribute_value END), ',') AS amount_array
  FROM
    ` + "`numia-data.osmosis.osmosis_message_event_attributes`" + `
  WHERE event_type = 'transfer'
  GROUP BY tx_id
)
SELECT 
  c.tx_id,
  c.ingestion_timestamp,
  c.recipient,
  t.amount_array
FROM claim_rewards_events c
JOIN transfer_amounts t ON c.tx_id = t.tx_id
WHERE c.is_target_contract AND c.is_claim_rewards
ORDER BY c.ingestion_timestamp DESC;`

	QueryCLVaultsDistributeRewards = `
WITH sender_message_index AS (
    SELECT 
        tx_id,
        message_index,
        attribute_index AS sender_attribute_index
    FROM
    ` + "`numia-data.osmosis.osmosis_message_event_attributes`" + `
    WHERE event_type IN ('total_collect_incentives', 'total_collect_spread_rewards')
        AND attribute_key = 'sender'
        AND attribute_value = '%s'
)

, amount_data AS (
    SELECT 
        a.tx_id,
        a.message_index,
        a.event_type,
        MAX(a.ingestion_timestamp) AS ingestion_timestamp,
        ARRAY_AGG(a.attribute_value) AS amount_array
    FROM
    ` + "`numia-data.osmosis.osmosis_message_event_attributes`" + `
	a
    JOIN sender_message_index s
    ON a.tx_id = s.tx_id 
        AND a.message_index = s.message_index
    WHERE a.event_type IN ('total_collect_incentives', 'total_collect_spread_rewards')
        AND a.attribute_key = 'tokens_out'
        AND a.ingestion_timestamp > TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL %d DAY) -- Added this line to filter by the last N days interval
    GROUP BY a.tx_id, a.message_index, a.event_type
)

, unnest_amounts AS (
    SELECT 
        tx_id,
        event_type,
        ingestion_timestamp,
        amount
    FROM amount_data,
    UNNEST(SPLIT(ARRAY_TO_STRING(amount_array, ','))) AS amount
)

, incentives AS (
    SELECT 
        tx_id,
        ingestion_timestamp,
        ARRAY_AGG(DISTINCT amount ORDER BY amount ASC) AS amount_incentives
    FROM unnest_amounts
    WHERE event_type = 'total_collect_incentives'
    GROUP BY tx_id, ingestion_timestamp
)

, spread_rewards AS (
    SELECT 
        tx_id,
        ingestion_timestamp,
        ARRAY_AGG(DISTINCT amount ORDER BY amount ASC) AS amount_spread_rewards
    FROM unnest_amounts
    WHERE event_type = 'total_collect_spread_rewards'
    GROUP BY tx_id, ingestion_timestamp
)

SELECT 
    COALESCE(i.tx_id, sr.tx_id) AS tx_id,
    COALESCE(i.ingestion_timestamp, sr.ingestion_timestamp) AS ingestion_timestamp,
    i.amount_incentives,
    sr.amount_spread_rewards
FROM incentives i
FULL JOIN spread_rewards sr
ON i.tx_id = sr.tx_id AND i.ingestion_timestamp = sr.ingestion_timestamp
ORDER BY ingestion_timestamp DESC;`
)
