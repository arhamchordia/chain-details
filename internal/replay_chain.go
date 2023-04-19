package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/arhamchordia/chain-details/types"
	"github.com/tendermint/tendermint/rpc/client/http"
)

// ReplayChainBond returns a file with the bond events in all the blocks given as startingHeight and endHeight
func ReplayChainBond(RPCAddress string, startingHeight, endHeight int64) error {
	// create an rpcClient with the given RPCAddress
	rpcClient, err := http.New(RPCAddress, types.Websocket)
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
			if strings.Contains(string(j.Data), types.IdentifierMsgExecuteContract) {
				for o, k := range j.Events {
					if k.Type == types.Message && string(k.Attributes[0].Value) == types.IdentifierMsgExecuteContract {
						if len(j.Events) >= o+3 {
							if j.Events[o+1].Type == types.Message && string(j.Events[o+1].Attributes[0].Value) == types.Wasm {
								if j.Events[o+2].Type == types.CoinSpent {
									if j.Events[o+3].Type == types.CoinReceived && string(j.Events[o+3].Attributes[0].Value) == types.VaultAddress {
										tempDepositorDetails = append(tempDepositorDetails, types.DepositorDetailsBond{
											Address:      string(j.Events[o+2].Attributes[0].Value),
											Amount:       string(j.Events[o+2].Attributes[1].Value),
											BlockHeight:  i,
											VaultAddress: types.VaultAddress,
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
					if q.Type == "wasm" && string(q.Attributes[0].Value) == types.VaultAddress && string(q.Attributes[1].Key) == "bond_id" {
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

	err = os.WriteFile("replay"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ReplayChainUnbond returns a file with the unbond events in all the blocks given as startingHeight and endHeight
func ReplayChainUnbond(RPCAddress string, startingHeight, endHeight int64) error {
	// create an rpcClient with the given RPCAddress
	rpcClient, err := http.New(RPCAddress, types.Websocket)
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
			if strings.Contains(string(j.Data), types.IdentifierMsgExecuteContract) {
				for _, k := range j.Events {
					if k.Type == types.Wasm {
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

	err = os.WriteFile("replay-unbond"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

// CheckLockedTokens returns a file with the locked tokens events in all the blocks given as startingHeight and endHeight
func CheckLockedTokens(RPCAddress string, startingHeight, endHeight int64) error {
	// create an rpcClient with the given RPCAddress
	rpcClient, err := http.New(RPCAddress, types.Websocket)
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
			if strings.Contains(string(j.Data), types.IdentifierMsgUpdateClient) && strings.Contains(string(j.Data), types.IdentifierMsgAcknowledgement) {
				for _, k := range j.Events {
					if k.Type == types.Wasm && len(k.Attributes) == 3 {
						if string(k.Attributes[0].Key) == types.ContractAddress && string(k.Attributes[1].Key) == types.LockID && string(k.Attributes[2].Key) == types.LockedTokens {
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
					if k.Type == types.Wasm && len(k.Attributes) == 5 {
						if string(k.Attributes[0].Key) == types.ContractAddress && string(k.Attributes[1].Key) == types.Action &&
							string(k.Attributes[2].Key) == types.CallbackInfo && string(k.Attributes[3].Key) == types.ReplyMsgID &&
							string(k.Attributes[4].Key) == types.ReplyResult {
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

	err = os.WriteFile("lock-details"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ParseMints returns a file with the mint tokens in incentive contract events in all the blocks given as startingHeight and endHeight
func ParseMints(RPCAddress string, startingHeight, endHeight int64) error {
	// create an rpcClient with the given RPCAddress
	rpcClient, err := http.New(RPCAddress, types.Websocket)
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
			if strings.Contains(j.String(), types.VaultTokenBalance) {
				for _, k := range j.Events {
					if k.Type == types.Wasm && len(k.Attributes) > 3 {
						if string(k.Attributes[0].Key) == types.ContractAddress && string(k.Attributes[1].Key) == types.Action &&
							string(k.Attributes[2].Key) == types.User && string(k.Attributes[3].Key) == types.VaultTokenBalance {
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

	err = os.WriteFile("minted-shares"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

// CallBackInfos returns a file with the callback info of primitives events in all the blocks given as startingHeight and endHeight
func CallBackInfos(RPCAddress string, startingHeight, endHeight int64) error {
	// create an rpcClient with the given RPCAddress
	rpcClient, err := http.New(RPCAddress, types.Websocket)
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
			if strings.Contains(j.String(), types.ReplyResult) && strings.Contains(j.String(), types.CallbackInfo) {
				for _, k := range j.Events {
					if k.Type == types.Wasm && len(k.Attributes) == 5 {
						if string(k.Attributes[0].Key) == types.ContractAddress && string(k.Attributes[1].Key) == types.Action &&
							string(k.Attributes[2].Key) == types.CallbackInfo && string(k.Attributes[3].Key) == types.ReplyMsgID &&
							string(k.Attributes[4].Key) == types.ReplyResult {
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

	err = os.WriteFile("callback-infos"+"-"+strconv.FormatInt(startingHeight, 10)+"-"+strconv.FormatInt(endHeight, 10)+".json", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

func CheckString(RPCAddress string, startingHeight, endHeight int64) error {
	rpcClient, err := http.New(RPCAddress, types.Websocket)
	if err != nil {
		return err
	}

	for i := startingHeight; i <= endHeight; i++ {
		//if i%1000 == 0 {
		//	fmt.Println(i)
		//}
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
				for _, l := range k.Attributes {
					if string(l.Key) == "start-unbond-status" && string(l.Value) == "starting-unbond" {
						fmt.Println(i)
					}
				}
			}
		}
	}
	return nil
}

//func ReplayChain(RPCAddress string, startingHeight, endHeight int64) error {
//	// create an rpcClient with the given RPCAddress
//	rpcClient, err := http.New(RPCAddress, types.Websocket)
//	if err != nil {
//		return err
//	}
//
//	// bonds
//	var depositorDetails []types.DepositorDetailsBond
//
//	for i := startingHeight; i <= endHeight; i++ {
//		if i%1000 == 0 {
//			fmt.Println(i)
//		}
//		time.Sleep(time.Millisecond * 10)
//
//		// get block results
//		blockResults, err := rpcClient.BlockResults(context.Background(), &i)
//		if err != nil {
//			return err
//		}
//
//		if blockResults.Height == 0 {
//			return fmt.Errorf("cannot read height %d", i)
//		}
//
//		// filter bonds
//		// filter unbonds
//		// locked tokens
//		// mints
//		// call back infos
//	}
//}
