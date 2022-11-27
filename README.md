# `chain-details`

## Overview 
This is a basic project which generates scv files for any CosmosSDK based chain.

## Capabilities
Currently, it provides 2 commands :
- `validators-data` : It outputs a csv file with validators detail like moniker, self-delegation, percentage voting power etc.
- `delegators-data` : It outputs 2 csv files with the following details
  - One file with entry of delegator with multiple entries corresponding to every validator that delegator has delegated to.
  - Second file with aggregated shares of every delegator present on chain.

## How to use 

```
$ go bin main.go valdiators-data "grpc-url-target-chain" "account-prefix-target-chain"
```

```
$ go bin main.go delegators-data "grpc-url-target-chain"
```