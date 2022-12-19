# `chain-details`

## Overview

This is a basic project which generates CSV files of the requested data for any CosmosSDK based chain.

## Install Golang :

Follow the instruction on how to install golang on your system through this
link [to install go-lang on Mac/Windows/Linux](https://go.dev/doc/install)

## Capabilities

Currently, it provides 3 commands :

- `validators-data` : It outputs a csv file with validators detail like moniker, self-delegation, percentage voting
  power etc.
- `delegators-data` : It outputs 2 csv files with the following details
    - One file with entry of delegator with multiple entries corresponding to every validator that delegator has
      delegated to.
    - Second file with aggregated shares of every delegator present on chain.
- `vesting-accounts` : It outputs a csv file with information on all type of vesting accounts
  like `DelayedVestingAccount`, `PeriodicVestingAccount`, `PermanentLockedAccount` and `PeriodicVestingAccount`
    - It has information like Address, Vesting Tokens, Vesting End Time, Vesting Start Time, Tokens freed up on every
      block, tokens freed up every day.

## How to use

- Get the required links for data collections :
    - Get the grpc urls of the desired chain from [chain-registry](https://github.com/cosmos/chain-registry).
    - Get the prefix for the specific chain from the same repo mentioned for grpc url.
    - Get the denom of the specific chain from the same repo mentioned for grpc url.
    - For genesis json url, please contact the respective chains or browse through the chain's github repo to find a
      link to genesis json.


- Once you have the required details. Just clone this repository on your system using :

```
git clone https://github.com/arhamchordia/chain-details
```

- After cloning go to the folder and open a terminal there.

- In the terminal, run the following commands with the details collected above.

```
$ go run main.go valdiators-data "grpc-url-target-chain" "account-prefix-target-chain"
```

```
$ go run main.go delegators-data "grpc-url-target-chain"
```

```
$ go run main.go vesting-accounts "genesis-json-data-url" "denom-for-chain"
```

- Once the commands have completed all the operations. You can find the respective CSV files in the cloned repo itself.
