# `chain-details`

## Overview 
This is a basic project which generates scv files for any CosmosSDK based chain.

## Capabilities
Currently, it provides 2 commands :
- `validators-data` : It outputs a csv file with validators detail like moniker, self-delegation, percentage voting power etc.
- `delegators-data` : It outputs 2 csv files with the following details
  - One file with entry of delegator with multiple entries corresponding to every validator that delegator has delegated to.
  - Second file with aggregated shares of every delegator present on chain.
- `vesting-accounts` : It outputs a csv file with information on all type of vesting accounts like `DelayedVestingAccount`, `PeriodicVestingAccount`, `PermanentLockedAccount` and `PeriodicVestingAccount`
  - It has information like Address, Vesting Tokens, Vesting End Time, Vesting Start Time, Tokens freed up on every block, tokens freed up every day. 

## How to use 

```
$ go run main.go valdiators-data "grpc-url-target-chain" "account-prefix-target-chain"
```

```
$ go run main.go delegators-data "grpc-url-target-chain"
```

```
$ go run main.go vesting-accounts "genesis-json-data-url" "denom-for-chain"
```