package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	grpctypes "github.com/arhamchordia/chain-details/types/grpc"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/arhamchordia/chain-details/types"
	"github.com/tendermint/tendermint/rpc/client/http"
)

// QueryDepositorsBond returns a file with the bond events in all the blocks given as startingHeight and endHeight
func QueryDepositorsBond(RPCAddress string, startingHeight, endHeight int64) error {
	// create an rpcClient with the given RPCAddress
	rpcClient, err := http.New(RPCAddress, grpctypes.Websocket)
	if err != nil {
		return err
	}

	var depositorDetails []types.DepositorDetailsBond
	for i := startingHeight; i <= endHeight; i++ {
		if i%1000 == 0 {
			fmt.Println(i)
		}
		time.Sleep(time.Millisecond * 10)

		// get block results
		blockResults, err := rpcClient.BlockResults(context.Background(), &i)
		if err != nil {
			return err
		}

		if blockResults.Height == 0 {
			return fmt.Errorf("cannot read height %d", i)
		}

		// iterate all the block transaction results to match the field we are looking for
		for _, j := range blockResults.TxsResults {
			var tempDepositorDetails []types.DepositorDetailsBond
			var tempBondIDs []int64
			if strings.Contains(string(j.Data), grpctypes.IdentifierMsgExecuteContract) {
				for o, k := range j.Events {
					if k.Type == grpctypes.Message && string(k.Attributes[0].Value) == grpctypes.IdentifierMsgExecuteContract {
						if len(j.Events) >= o+3 {
							if j.Events[o+1].Type == grpctypes.Message && string(j.Events[o+1].Attributes[0].Value) == grpctypes.Wasm {
								if j.Events[o+2].Type == grpctypes.CoinSpent {
									if j.Events[o+3].Type == grpctypes.CoinReceived && string(j.Events[o+3].Attributes[0].Value) == grpctypes.VaultAddress {
										tempDepositorDetails = append(tempDepositorDetails, types.DepositorDetailsBond{
											Address:      string(j.Events[o+2].Attributes[0].Value),
											Amount:       string(j.Events[o+2].Attributes[1].Value),
											BlockHeight:  i,
											VaultAddress: grpctypes.VaultAddress,
										})
									}
								}
							}
						} else {
							fmt.Println("couldn't find the next 2 events on block height :", i)
						}
					}
				}

				for _, q := range j.Events {
					if q.Type == grpctypes.Wasm && string(q.Attributes[0].Value) == grpctypes.VaultAddress && string(q.Attributes[1].Key) == grpctypes.BondID {
						tempBondID, err := strconv.ParseInt(string(q.Attributes[1].Value), 10, 64)
						if err != nil {
							return fmt.Errorf("incorrect bond ID at height %d", i)
						}
						tempBondIDs = append(tempBondIDs, tempBondID)
					}
				}

				if len(tempBondIDs) != len(tempDepositorDetails) {
					fmt.Println(tempBondIDs)
					fmt.Println(tempDepositorDetails)
					return fmt.Errorf("mismatch in the counting of bond IDs %d", i)
				}

				for p := range tempBondIDs {
					depositorDetails = append(depositorDetails, types.DepositorDetailsBond{
						Address:      tempDepositorDetails[p].Address,
						Amount:       tempDepositorDetails[p].Amount,
						BlockHeight:  tempDepositorDetails[p].BlockHeight,
						VaultAddress: tempDepositorDetails[p].VaultAddress,
						BondID:       tempBondIDs[p],
					})
				}
			}
		}
	}

	// marshal and write the contents in a file
	file, err := json.MarshalIndent(depositorDetails, "", " ")
	if err != nil {
		return err
	}

	// TODO add prefix and rename
	err = os.WriteFile("replay"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

// QueryDepositorsUnbond returns a file with the unbond events in all the blocks given as startingHeight and endHeight
func QueryDepositorsUnbond(RPCAddress string, startingHeight, endHeight int64) error {
	// create an rpcClient with the given RPCAddress
	rpcClient, err := http.New(RPCAddress, grpctypes.Websocket)
	if err != nil {
		return err
	}

	var depositorDetailsUnbond []types.DepositorDetailsUnbond

	for i := startingHeight; i <= endHeight; i++ {
		if i%1000 == 0 {
			fmt.Println(i)
		}
		time.Sleep(time.Millisecond * 10)

		// get block results
		blockResults, err := rpcClient.BlockResults(context.Background(), &i)
		if err != nil {
			return err
		}

		if blockResults.Height == 0 {
			return fmt.Errorf("cannot read height %d", i)
		}

		// iterate all the block transaction results to match the field we are looking for
		for _, j := range blockResults.TxsResults {
			if strings.Contains(string(j.Data), grpctypes.IdentifierMsgExecuteContract) {
				for _, k := range j.Events {
					if k.Type == grpctypes.Wasm {
						if len(k.Attributes) == 5 {
							unbondID, err := strconv.ParseInt(string(k.Attributes[4].Value), 10, 64)
							if err != nil {
								return fmt.Errorf("incorrect unbond ID at height %d", i)
							}
							depositorDetailsUnbond = append(depositorDetailsUnbond, types.DepositorDetailsUnbond{
								Address:      string(k.Attributes[2].Value),
								BlockHeight:  i,
								BurntShares:  string(k.Attributes[3].Value),
								VaultAddress: string(k.Attributes[0].Value),
								UnbondID:     unbondID,
							})
						}
					}
				}
			}
		}
	}

	// marshal and write the contents in a file
	file, err := json.MarshalIndent(depositorDetailsUnbond, "", " ")
	if err != nil {
		return err
	}

	// TODO add prefix and rename

	err = os.WriteFile("replay-unbond"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

// QueryDepositorsLockedTokens returns a file with the locked tokens events in all the blocks given as startingHeight and endHeight
func QueryDepositorsLockedTokens(RPCAddress string, startingHeight, endHeight int64) error {
	// create an rpcClient with the given RPCAddress
	rpcClient, err := http.New(RPCAddress, grpctypes.Websocket)
	if err != nil {
		return err
	}

	var lockDetailsByHeight []types.LockDetailsByHeight

	for i := startingHeight; i <= endHeight; i++ {
		time.Sleep(time.Millisecond * 50)

		// get block results
		blockResults, err := rpcClient.BlockResults(context.Background(), &i)
		if err != nil {
			return err
		}

		if blockResults.Height == 0 {
			return fmt.Errorf("cannot read height %d", i)
		}

		tempContractDetailsMap := make(map[string]types.ContractDetails)
		var tempContractDetails []types.ContractDetails

		// iterate all the block transaction results to match the field we are looking for
		for _, j := range blockResults.TxsResults {
			if strings.Contains(string(j.Data), grpctypes.IdentifierMsgUpdateClient) && strings.Contains(string(j.Data), grpctypes.IdentifierMsgAcknowledgement) {
				for _, k := range j.Events {
					if k.Type == grpctypes.Wasm && len(k.Attributes) == 3 {
						if string(k.Attributes[0].Key) == grpctypes.ContractAddress && string(k.Attributes[1].Key) == grpctypes.LockID && string(k.Attributes[2].Key) == grpctypes.LockedTokens {
							lockID, err := strconv.ParseInt(string(k.Attributes[1].Value), 10, 64)
							if err != nil {
								return fmt.Errorf("incorrect lock ID at height %d", i)
							}
							tempContractDetailsMap[string(k.Attributes[0].Value)] = types.ContractDetails{
								Address:                 string(k.Attributes[0].Value),
								LockID:                  lockID,
								LockedTokensProtoString: string(k.Attributes[2].Value),
							}
						}
					}
					if k.Type == grpctypes.Wasm && len(k.Attributes) == 5 {
						if string(k.Attributes[0].Key) == grpctypes.ContractAddress && string(k.Attributes[1].Key) == grpctypes.Action &&
							string(k.Attributes[2].Key) == grpctypes.CallbackInfo && string(k.Attributes[3].Key) == grpctypes.ReplyMsgID &&
							string(k.Attributes[4].Key) == grpctypes.ReplyResult {
							value, ok := tempContractDetailsMap[string(k.Attributes[0].Value)]
							if ok {
								value.Action = string(k.Attributes[1].Value)
								value.CallbackInfo = string(k.Attributes[2].Value)
								value.ReplyMessageID = string(k.Attributes[3].Value)
								value.ReplyResult = string(k.Attributes[4].Value)
								tempContractDetailsMap[string(k.Attributes[0].Value)] = value
							} else {
								return fmt.Errorf("unable to find the primitve address in the map at hegiht %d", i)
							}
						}
					}
				}
			}
		}

		if len(tempContractDetailsMap) > 0 {
			for _, value := range tempContractDetailsMap {
				tempContractDetails = append(tempContractDetails, value)
			}
			lockDetailsByHeight = append(lockDetailsByHeight, types.LockDetailsByHeight{
				Height:          i,
				ContractDetails: tempContractDetails,
			})
		}
	}

	// marshal and write the contents in a file
	file, err := json.MarshalIndent(lockDetailsByHeight, "", " ")
	if err != nil {
		return err
	}

	// TODO add prefix and rename
	err = os.WriteFile("lock-details"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

// QueryDepositorsMints returns a file with the mint tokens in incentive contract events in all the blocks given as startingHeight and endHeight
func QueryDepositorsMints(RPCAddress string, startingHeight, endHeight int64) error {
	// create an rpcClient with the given RPCAddress
	rpcClient, err := http.New(RPCAddress, grpctypes.Websocket)
	if err != nil {
		return err
	}

	addressToSharesMap := make(map[string]types.AddressSharesInIncentiveContract)

	for i := startingHeight; i <= endHeight; i++ {
		time.Sleep(time.Millisecond * 50)

		// get block results
		blockResults, err := rpcClient.BlockResults(context.Background(), &i)
		if err != nil {
			return err
		}

		if blockResults.Height == 0 {
			return fmt.Errorf("cannot read height %d", i)
		}

		// iterate all the block transaction results to match the field we are looking for
		for _, j := range blockResults.TxsResults {
			if strings.Contains(j.String(), grpctypes.VaultTokenBalance) {
				for _, k := range j.Events {
					if k.Type == grpctypes.Wasm && len(k.Attributes) > 3 {
						if string(k.Attributes[0].Key) == grpctypes.ContractAddress && string(k.Attributes[1].Key) == grpctypes.Action &&
							string(k.Attributes[2].Key) == grpctypes.User && string(k.Attributes[3].Key) == grpctypes.VaultTokenBalance {
							if len(k.Attributes) > 4 {
								fmt.Println("found a block with multiple mints at height :", i)
							}
							value, ok := addressToSharesMap[string(k.Attributes[2].Value)]
							if ok {
								value.Shares = append(value.Shares, string(k.Attributes[3].Value))
								value.LastUpdatedHeight = append(value.LastUpdatedHeight, i)
								addressToSharesMap[string(k.Attributes[2].Value)] = value
							} else {
								addressToSharesMap[string(k.Attributes[2].Value)] = types.AddressSharesInIncentiveContract{
									Shares:            []string{string(k.Attributes[3].Value)},
									LastUpdatedHeight: []int64{i},
								}
							}
						}
					}
				}
			}
		}
	}

	// marshal and write the contents in a file
	file, err := json.MarshalIndent(addressToSharesMap, "", " ")
	if err != nil {
		return err
	}

	// TODO add prefix and rename
	err = os.WriteFile("minted-shares"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

// QueryDepositorsCallbackInfo returns a file with the callback info of primitives events in all the blocks given as startingHeight and endHeight
func QueryDepositorsCallbackInfo(RPCAddress string, startingHeight, endHeight int64) error {
	// create an rpcClient with the given RPCAddress
	rpcClient, err := http.New(RPCAddress, grpctypes.Websocket)
	if err != nil {
		return err
	}

	var callBackInfoWithHeight []types.CallBackInfoWithHeight
	for i := startingHeight; i <= endHeight; i++ {
		if i%1000 == 0 {
			fmt.Println(i)
		}
		time.Sleep(time.Millisecond * 20)

		// get block results
		blockResults, err := rpcClient.BlockResults(context.Background(), &i)
		if err != nil {
			return err
		}

		if blockResults.Height == 0 {
			return fmt.Errorf("cannot read height %d", i)
		}

		var tempCallBackInfos []types.CallBackInfo

		// iterate all the block transaction results to match the field we are looking for
		for _, j := range blockResults.TxsResults {
			if strings.Contains(j.String(), grpctypes.ReplyResult) && strings.Contains(j.String(), grpctypes.CallbackInfo) {
				for _, k := range j.Events {
					if k.Type == grpctypes.Wasm && len(k.Attributes) == 5 {
						if string(k.Attributes[0].Key) == grpctypes.ContractAddress && string(k.Attributes[1].Key) == grpctypes.Action &&
							string(k.Attributes[2].Key) == grpctypes.CallbackInfo && string(k.Attributes[3].Key) == grpctypes.ReplyMsgID &&
							string(k.Attributes[4].Key) == grpctypes.ReplyResult {
							tempCallBackInfos = append(tempCallBackInfos, types.CallBackInfo{
								ContractAddress:    string(k.Attributes[0].Value),
								Action:             string(k.Attributes[1].Value),
								CallBackInfoString: string(k.Attributes[2].Value),
								ReplyMsgID:         string(k.Attributes[3].Value),
								ReplyResult:        string(k.Attributes[4].Value),
							})
						}
					}
				}
			}
		}

		if len(tempCallBackInfos) > 0 {
			callBackInfoWithHeight = append(callBackInfoWithHeight, types.CallBackInfoWithHeight{
				Height:        i,
				CallBackInfos: tempCallBackInfos,
			})
		}
	}

	// marshal and write the contents in a file
	file, err := json.MarshalIndent(callBackInfoWithHeight, "", " ")
	if err != nil {
		return err
	}

	// TODO add prefix and rename
	err = os.WriteFile("callback-infos"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func QueryDepositorsBeginUnlocking(RPCAddress string, startingHeight, endHeight int64) error {
	rpcClient, err := http.New(RPCAddress, grpctypes.Websocket)
	if err != nil {
		return err
	}

	var BeginUnlocking []types.BeginUnlocking
	for i := startingHeight; i <= endHeight; i++ {
		if i%1000 == 0 {
			fmt.Println(i)
		}
		time.Sleep(time.Millisecond * 20)

		// get block results
		blockResults, err := rpcClient.BlockResults(context.Background(), &i)
		if err != nil {
			return err
		}

		if blockResults.Height == 0 {
			return fmt.Errorf("cannot read height %d", i)
		}

		// iterate all the block transaction results to match the field we are looking for
		for _, j := range blockResults.TxsResults {
			for _, k := range j.Events {
				if k.Type == "wasm" && len(k.Attributes) == 3 {
					if string(k.Attributes[2].Key) == "step" && strings.Contains(string(k.Attributes[2].Value), "BeginUnlocking") {
						BeginUnlocking = append(BeginUnlocking, types.BeginUnlocking{
							Height:          i,
							Step:            string(k.Attributes[2].Value),
							ContractAddress: string(k.Attributes[0].Value),
							PendingMsg:      string(k.Attributes[1].Value),
						})
					}
				}

			}
		}
	}

	// marshal and write the contents in a file
	file, err := json.MarshalIndent(BeginUnlocking, "", " ")
	if err != nil {
		return err
	}

	// TODO add prefix and rename
	err = os.WriteFile("begin-unlocking"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func QueryDepositorsReplayChain(RPCAddress string, startingHeight, endHeight int64) error {
	// create an rpcClient with the given RPCAddress
	rpcClient, err := http.New(RPCAddress, grpctypes.Websocket)
	if err != nil {
		return err
	}

	var depositorDetails []types.DepositorDetailsBond
	var depositorDetailsUnbond []types.DepositorDetailsUnbond
	var lockDetailsByHeight []types.LockDetailsByHeight
	addressToSharesMap := make(map[string]types.AddressSharesInIncentiveContract)
	var callBackInfoWithHeight []types.CallBackInfoWithHeight
	var beginUnlocking []types.BeginUnlocking

	for i := startingHeight; i <= endHeight; i++ {
		if i%1000 == 0 {
			fmt.Println(i)
		}
		time.Sleep(time.Millisecond * 10)

		// get block results
		blockResults, err := rpcClient.BlockResults(context.Background(), &i)
		if err != nil {
			return err
		}

		if blockResults.Height == 0 {
			return fmt.Errorf("cannot read height %d", i)
		}

		// filter bonds, unbonds etc
		for _, j := range blockResults.TxsResults {
			var tempDepositorDetails []types.DepositorDetailsBond
			var tempBondIDs []int64
			tempContractDetailsMap := make(map[string]types.ContractDetails)
			var tempContractDetails []types.ContractDetails
			var tempCallBackInfos []types.CallBackInfo

			for o, k := range j.Events {
				// bond filters
				if k.Type == grpctypes.Message && string(k.Attributes[0].Value) == grpctypes.IdentifierMsgExecuteContract {
					if len(j.Events) >= o+3 {
						if j.Events[o+1].Type == grpctypes.Message && string(j.Events[o+1].Attributes[0].Value) == grpctypes.Wasm {
							if j.Events[o+2].Type == grpctypes.CoinSpent {
								if j.Events[o+3].Type == grpctypes.CoinReceived && string(j.Events[o+3].Attributes[0].Value) == grpctypes.VaultAddress {
									tempDepositorDetails = append(tempDepositorDetails, types.DepositorDetailsBond{
										Address:      string(j.Events[o+2].Attributes[0].Value),
										Amount:       string(j.Events[o+2].Attributes[1].Value),
										BlockHeight:  i,
										VaultAddress: grpctypes.VaultAddress,
									})
								}
							}
						}
					} else {
						fmt.Println("couldn't find the next 2 events on block height :", i)
					}
				}

				// unbond filters
				if k.Type == grpctypes.Wasm && len(k.Attributes) == 5 {
					if string(k.Attributes[1].Key) == "action" && string(k.Attributes[1].Value) == "start_unbond" &&
						string(k.Attributes[3].Key) == "burnt" {
						unbondID, _ := strconv.ParseInt(string(k.Attributes[4].Value), 10, 64)
						if err != nil {
							return fmt.Errorf("incorrect unbond ID at height %d", i)
						}
						depositorDetailsUnbond = append(depositorDetailsUnbond, types.DepositorDetailsUnbond{
							Address:      string(k.Attributes[2].Value),
							BlockHeight:  i,
							BurntShares:  string(k.Attributes[3].Value),
							VaultAddress: string(k.Attributes[0].Value),
							UnbondID:     unbondID,
						})
					}
				}

				// lock tokens
				if strings.Contains(string(j.Data), grpctypes.IdentifierMsgUpdateClient) && strings.Contains(string(j.Data), grpctypes.IdentifierMsgAcknowledgement) {
					if k.Type == grpctypes.Wasm && len(k.Attributes) == 3 {
						if string(k.Attributes[0].Key) == grpctypes.ContractAddress && string(k.Attributes[1].Key) == grpctypes.LockID && string(k.Attributes[2].Key) == grpctypes.LockedTokens {
							lockID, err := strconv.ParseInt(string(k.Attributes[1].Value), 10, 64)
							if err != nil {
								return fmt.Errorf("incorrect lock ID at height %d", i)
							}
							tempContractDetailsMap[string(k.Attributes[0].Value)] = types.ContractDetails{
								Address:                 string(k.Attributes[0].Value),
								LockID:                  lockID,
								LockedTokensProtoString: string(k.Attributes[2].Value),
							}
						}
					}
					if k.Type == grpctypes.Wasm && len(k.Attributes) == 5 {
						if string(k.Attributes[0].Key) == grpctypes.ContractAddress && string(k.Attributes[1].Key) == grpctypes.Action &&
							string(k.Attributes[2].Key) == grpctypes.CallbackInfo && string(k.Attributes[3].Key) == grpctypes.ReplyMsgID &&
							string(k.Attributes[4].Key) == grpctypes.ReplyResult {
							value, ok := tempContractDetailsMap[string(k.Attributes[0].Value)]
							if ok {
								value.Action = string(k.Attributes[1].Value)
								value.CallbackInfo = string(k.Attributes[2].Value)
								value.ReplyMessageID = string(k.Attributes[3].Value)
								value.ReplyResult = string(k.Attributes[4].Value)
								tempContractDetailsMap[string(k.Attributes[0].Value)] = value
							} else {
								return fmt.Errorf("unable to find the primitve address in the map at hegiht %d", i)
							}
						}
					}
				}

				// parse mints
				if strings.Contains(j.String(), grpctypes.VaultTokenBalance) {
					if k.Type == grpctypes.Wasm && len(k.Attributes) > 3 {
						if string(k.Attributes[0].Key) == grpctypes.ContractAddress && string(k.Attributes[1].Key) == grpctypes.Action &&
							string(k.Attributes[2].Key) == grpctypes.User && string(k.Attributes[3].Key) == grpctypes.VaultTokenBalance {
							if len(k.Attributes) > 4 {
								fmt.Println("found a block with multiple mints at height :", i)
							}
							value, ok := addressToSharesMap[string(k.Attributes[2].Value)]
							if ok {
								value.Shares = append(value.Shares, string(k.Attributes[3].Value))
								value.LastUpdatedHeight = append(value.LastUpdatedHeight, i)
								addressToSharesMap[string(k.Attributes[2].Value)] = value
							} else {
								addressToSharesMap[string(k.Attributes[2].Value)] = types.AddressSharesInIncentiveContract{
									Shares:            []string{string(k.Attributes[3].Value)},
									LastUpdatedHeight: []int64{i},
								}
							}
						}
					}
				}

				// callback infos
				if strings.Contains(j.String(), grpctypes.ReplyResult) && strings.Contains(j.String(), grpctypes.CallbackInfo) {
					if k.Type == grpctypes.Wasm && len(k.Attributes) == 5 {
						if string(k.Attributes[0].Key) == grpctypes.ContractAddress && string(k.Attributes[1].Key) == grpctypes.Action &&
							string(k.Attributes[2].Key) == grpctypes.CallbackInfo && string(k.Attributes[3].Key) == grpctypes.ReplyMsgID &&
							string(k.Attributes[4].Key) == grpctypes.ReplyResult {
							tempCallBackInfos = append(tempCallBackInfos, types.CallBackInfo{
								ContractAddress:    string(k.Attributes[0].Value),
								Action:             string(k.Attributes[1].Value),
								CallBackInfoString: string(k.Attributes[2].Value),
								ReplyMsgID:         string(k.Attributes[3].Value),
								ReplyResult:        string(k.Attributes[4].Value),
							})
						}
					}
				}

				// begin unlocking
				if k.Type == "wasm" && len(k.Attributes) == 3 {
					if string(k.Attributes[2].Key) == "step" && strings.Contains(string(k.Attributes[2].Value), "BeginUnlocking") {
						beginUnlocking = append(beginUnlocking, types.BeginUnlocking{
							Height:          i,
							Step:            string(k.Attributes[2].Value),
							ContractAddress: string(k.Attributes[0].Value),
							PendingMsg:      string(k.Attributes[1].Value),
						})
					}
				}
			}

			for _, q := range j.Events {
				if q.Type == "wasm" && string(q.Attributes[0].Value) == grpctypes.VaultAddress && string(q.Attributes[1].Key) == "bond_id" {
					tempBondID, err := strconv.ParseInt(string(q.Attributes[1].Value), 10, 64)
					if err != nil {
						return fmt.Errorf("incorrect bond ID at height %d", i)
					}
					tempBondIDs = append(tempBondIDs, tempBondID)
				}
			}

			if len(tempBondIDs) != len(tempDepositorDetails) {
				fmt.Println(tempBondIDs)
				fmt.Println(tempDepositorDetails)
				return fmt.Errorf("mismatch in the counting of bond IDs %d", i)
			}

			for p := range tempBondIDs {
				depositorDetails = append(depositorDetails, types.DepositorDetailsBond{
					Address:      tempDepositorDetails[p].Address,
					Amount:       tempDepositorDetails[p].Amount,
					BlockHeight:  tempDepositorDetails[p].BlockHeight,
					VaultAddress: tempDepositorDetails[p].VaultAddress,
					BondID:       tempBondIDs[p],
				})
			}

			if len(tempContractDetailsMap) > 0 {
				for _, value := range tempContractDetailsMap {
					tempContractDetails = append(tempContractDetails, value)
				}
				lockDetailsByHeight = append(lockDetailsByHeight, types.LockDetailsByHeight{
					Height:          i,
					ContractDetails: tempContractDetails,
				})
			}

			if len(tempCallBackInfos) > 0 {
				callBackInfoWithHeight = append(callBackInfoWithHeight, types.CallBackInfoWithHeight{
					Height:        i,
					CallBackInfos: tempCallBackInfos,
				})
			}
		}
	}

	// store filter bonds
	// marshal and write the contents in a file
	bondFile, err := json.MarshalIndent(depositorDetails, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile("replay-bond"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", bondFile, 0644)
	if err != nil {
		return err
	}

	// store filter unbonds
	// marshal and write the contents in a file
	unbondFile, err := json.MarshalIndent(depositorDetailsUnbond, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile("replay-unbond"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", unbondFile, 0644)
	if err != nil {
		return err
	}

	// store locked tokens
	// marshal and write the contents in a file
	lockDetailsFile, err := json.MarshalIndent(lockDetailsByHeight, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile("lock-details"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", lockDetailsFile, 0644)
	if err != nil {
		return err
	}

	// store mints
	// marshal and write the contents in a file
	mintFile, err := json.MarshalIndent(addressToSharesMap, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile("minted-shares"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", mintFile, 0644)
	if err != nil {
		return err
	}

	// store call back infos
	// marshal and write the contents in a file
	callbackInfoFile, err := json.MarshalIndent(callBackInfoWithHeight, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile("callback-infos"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", callbackInfoFile, 0644)
	if err != nil {
		return err
	}

	// store begin unlocking
	// marshal and write the contents in a file
	beginUnlockingFile, err := json.MarshalIndent(beginUnlocking, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile("begin-unlocking"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", beginUnlockingFile, 0644)
	if err != nil {
		return err
	}

	return nil
}
