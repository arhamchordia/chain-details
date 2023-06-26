# Chain Details Project

- [Requirements](#requirements)
- [Usage](#usage)
- [Capabilities](#capabilities)
    - [gRPC](#grpc)
    - [BigQuery](#bigquery)

---

## Overview

This is a basic project which generates CSV and JSON files of the requested data for any CosmosSDK based chain.

---

## Requirements

In order to use `chain-details`, you need to have Go installed on your system. If you don't have it already, you can
follow the instructions provided in [this link](https://go.dev/doc/install).

### gRPC Endpoint

To use the gRPC commands, you need to obtain the grpc urls for the desired chain. These can be found in the
chain-registry repository. You also need to obtain the genesis json url for the chain. This can usually be found in the
[chain-registry repository](https://github.com/cosmos/chain-registry).

### BigQuery with Numia Data

If you want to use `chain-details` to query data from the Numia Data collection on Google Cloud BigQuery, you need to
first import the `numia-data` collection on BigQuery, you can follow
the [instructions on their documentation](https://docs.numia.xyz/using-numia/querying-numia-datasets).

Once you have imported the collection, you can use `chain-details` to execute SQL queries against it. To authenticate
`chain-details` with your Google Cloud account, you can use the gcloud CLI, and run `gcloud auth login` and `gcloud auth
application-default login`.

You also need to set the `GOOGLE_CLOUD_PROJECT_ID` environment variable to your project ID by running:

```bash
export GOOGLE_CLOUD_PROJECT_ID=<your_project_id>
```

---

## Usage

To use chain-details, you need to first gather the required information and resources for the specific chain you are
working with. This includes grpc endpoints, genesis urls, address prefixes, and denoms.

Once you have gathered all the necessary information, you can run the relevant commands to generate the required
CSV/JSON files.

The output CSV/JSON files can be found in the `/output` directory in the root of the chain-details project.

---

## Capabilities

The chain-details project offers various commands to generate CSV and JSON files of the requested data for any
CosmosSDK-based chain.

By default, the output format is CSV, but it can be changed to JSON using the `--output [csv/json]` flag.

### gRPC

#### Delegators

`delegators-data`: Generates two files with information about delegators present on the chain. One file contains a list
of delegators along with the validator they have delegated to, and the other contains the total aggregated shares of
every delegator on the chain.

```bash
go run main.go grpc delegators-data [grpc-url]
```

#### Depositors

`grpc depositors-bond`: Generates a file with the details of all deposits made to the validators from a given height to
an end height.

```bash
go run main.go grpc depositors-bond [rpc-url] [start-height] [end-height]
```

`grpc depositors-unbond`: Generates a file with the details of all unbonds made to the validators from a given height to
an end height.

```bash
go run main.go grpc depositors-unbond [rpc-url] [start-height] [end-height]
```

`grpc depositors-locked-tokens`: Generates a file with the details of all tokens locked for a given duration by the
depositors from a given height to an end height.

```bash
go run main.go grpc depositors-locked-tokens [rpc-url] [start-height] [end-height]
```

`grpc depositors-mints`: Generates a file with the details of all minted tokens by the validators from a given height to
an end height.

```bash
go run main.go grpc depositors-mints [rpc-url] [start-height] [end-height]
```

`grpc depositors-callback-info`: Generates a file with the details of all callbacks made by the validators from a given
height to an end height.

```bash
go run main.go grpc depositors-callback-info [rpc-url] [start-height] [end-height]
```

`grpc depositors-begin-unlocking`: Generates a file with the details of all begin-unlocking events initiated by the
validators from a given height to an end height.

```bash
go run main.go grpc depositors-begin-unlocking [rpc-url] [start-height] [end-height]
```

#### Genesis

`grpc genesis-vesting-accounts`: Generates a file with the details of all vesting accounts on the chain like
DelayedVestingAccount, PeriodicVestingAccount, PermanentLockedAccount and PeriodicVestingAccount.

```bash
go run main.go grpc genesis-vesting-accounts [json-url] [denom]
```

#### Validators

`grpc validators-data`: Generates a file with details about all validators on the chain, including moniker,
self-delegation,
and percentage voting power.

```bash
go run main.go grpc validators-data [grpc-url] [account-address-prefix]
```

### BigQuery

(--flags) are optionals

#### Raw Query

`bigquery raw`: Executes a raw SQL query on Google Cloud BigQuery and generates a file with the results. The query
should be enclosed in double quotes and any backticks within the query must be escaped with a backslash.

```bash
go run main.go bigquery raw --query "SELECT * FROM \`numia-data.quasar.quasar_transactions\` ORDER BY \`block_height\` DESC LIMIT 1000"
```

#### Quasar Transactions

`bigquery transactions`: Executes a SQL query to scrape all the transactions for a given account address and generates a file with the results. Provide the address to query with the --address flag.
```bash
go run main.go bigquery transactions --address <account_address>
```

#### Quasar Vault Actions

##### Bond

`bigquery bond`: Executes a SQL query to scrape all the bond transactions for a given account address and generates a file with the results. Providing the address to query with the --address flag is optional.
```bash
go run main.go bigquery bond (--address <account_address>) (--confirmed || --pending)
```

##### Unbond

`bigquery unbond`: Executes a SQL query to scrape all the unbond transactions for a given account address and generates a file with the results. Providing the address to query with the --address flag is optional.
```bash
go run main.go bigquery unbond (--address <account_address>) (--confirmed || --pending)
```

##### Claim

`bigquery claim`: Executes a SQL query to scrape all the claim transactions for a given account address and generates a file with the results. Providing the address to query with the --address flag is optional.
```bash
go run main.go bigquery claim (--address <account_address>)
```