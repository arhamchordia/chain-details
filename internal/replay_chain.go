package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tendermint/tendermint/rpc/client/http"
)

type DepositorDetailsBond struct {
	Address      string `json:"address"`
	BlockHeight  int64  `json:"block_height"`
	Amount       string `json:"amount"`
	VaultAddress string `json:"primitive_address"`
	BondID       int64  `json:"bond_id"`
}

type DepositorDetailsUnbond struct {
	Address      string `json:"address"`
	BlockHeight  int64  `json:"block_height"`
	VaultAddress string `json:"primitive_address"`
	BurntShares  string `json:"burnt_shares"`
	UnbondID     int64  `json:"unbond_id"`
}

type LockDetailsByHeight struct {
	Height          int64             `json:"height"`
	ContractDetails []ContractDetails `json:"contract_details"`
}

type ContractDetails struct {
	Address                 string `json:"address"`
	LockID                  int64  `json:"lock_id"`
	LockedTokensProtoString string `json:"locked_tokens_proto_string"`
	Action                  string `json:"action"`
	CallbackInfo            string `json:"callback_info"`
	ReplyMessageID          string `json:"reply_message_id"`
	ReplyResult             string `json:"reply_result"`
}

type Test struct {
	Address           string   `json:"address"`
	Shares            []string `json:"shares"`
	LastUpdatedHeight []int64  `json:"last_updated_height"`
}

type AddressSharesInIncentiveContract struct {
	Shares            []string `json:"shares"`
	LastUpdatedHeight []int64  `json:"last_updated_height"`
}

type CallBackInfoWithHeight struct {
	Height        int64          `json:"height"`
	CallBackInfos []CallBackInfo `json:"callBackInfos"`
}

type CallBackInfo struct {
	ContractAddress    string `json:"contract_address"`
	Action             string `json:"action"`
	CallBackInfoString string `json:"call_back_info"`
	ReplyMsgID         string `json:"reply_msg_id"`
	ReplyResult        string `json:"reply_result"`
}

func ReplayChainBond(RPCAddress string, startingHeight, endHeight int64) error {
	rpcClient, err := http.New(RPCAddress, "/websocket")
	if err != nil {
		return err
	}

	var depositorDetails []DepositorDetailsBond
	for i := startingHeight; i <= endHeight; i++ {
		if i%1000 == 0 {
			fmt.Println(i)
		}
		time.Sleep(time.Millisecond * 10)

		blockResults, err := rpcClient.BlockResults(context.Background(), &i)
		if err != nil {
			return err
		}

		if blockResults.Height == 0 {
			return fmt.Errorf("cannot read height %d", i)
		}

		for _, j := range blockResults.TxsResults {
			var tempDepositorDetails []DepositorDetailsBond
			var tempBondIDs []int64
			if strings.Contains(string(j.Data), "/cosmwasm.wasm.v1.MsgExecuteContract") {
				for o, k := range j.Events {
					if k.Type == "message" && string(k.Attributes[0].Value) == "/cosmwasm.wasm.v1.MsgExecuteContract" {
						if len(j.Events) >= o+3 {
							if j.Events[o+1].Type == "message" && string(j.Events[o+1].Attributes[0].Value) == "wasm" {
								if j.Events[o+2].Type == "coin_spent" {
									if j.Events[o+3].Type == "coin_received" && string(j.Events[o+3].Attributes[0].Value) == "quasar18a2u6az6dzw528rptepfg6n49ak6hdzkf8ewf0n5r0nwju7gtdgqamr7qu" {
										// fmt.Println(string(j.Events[o+2].Attributes[0].Value), ":", string(j.Events[o+2].Attributes[1].Value))
										tempDepositorDetails = append(tempDepositorDetails, DepositorDetailsBond{
											Address:      string(j.Events[o+2].Attributes[0].Value),
											Amount:       string(j.Events[o+2].Attributes[1].Value),
											BlockHeight:  i,
											VaultAddress: "quasar18a2u6az6dzw528rptepfg6n49ak6hdzkf8ewf0n5r0nwju7gtdgqamr7qu",
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
					if q.Type == "wasm" && string(q.Attributes[0].Value) == "quasar18a2u6az6dzw528rptepfg6n49ak6hdzkf8ewf0n5r0nwju7gtdgqamr7qu" && string(q.Attributes[1].Key) == "bond_id" {
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
					depositorDetails = append(depositorDetails, DepositorDetailsBond{
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

func ReplayChainUnbond(RPCAddress string, startingHeight, endHeight int64) error {
	rpcClient, err := http.New(RPCAddress, "/websocket")
	if err != nil {
		return err
	}

	var depositorDetailsUnbond []DepositorDetailsUnbond
	for i := startingHeight; i <= endHeight; i++ {
		if i%1000 == 0 {
			fmt.Println(i)
		}
		time.Sleep(time.Millisecond * 10)

		blockResults, err := rpcClient.BlockResults(context.Background(), &i)
		if err != nil {
			return err
		}

		if blockResults.Height == 0 {
			return fmt.Errorf("cannot read height %d", i)
		}

		for _, j := range blockResults.TxsResults {
			if strings.Contains(string(j.Data), "/cosmwasm.wasm.v1.MsgExecuteContract") {
				for _, k := range j.Events {
					if k.Type == "wasm" {
						if len(k.Attributes) == 5 {
							unbondID, err := strconv.ParseInt(string(k.Attributes[4].Value), 10, 64)
							if err != nil {
								return fmt.Errorf("incorrect unbond ID at height %d", i)
							}
							depositorDetailsUnbond = append(depositorDetailsUnbond, DepositorDetailsUnbond{
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

func CheckLockedTokens(RPCAddress string, startingHeight, endHeight int64) error {
	rpcClient, err := http.New(RPCAddress, "/websocket")
	if err != nil {
		return err
	}

	var lockDetailsByHeight []LockDetailsByHeight

	for i := startingHeight; i <= endHeight; i++ {
		time.Sleep(time.Millisecond * 50)

		blockResults, err := rpcClient.BlockResults(context.Background(), &i)
		if err != nil {
			return err
		}

		if blockResults.Height == 0 {
			return fmt.Errorf("cannot read height %d", i)
		}

		tempContractDetailsMap := make(map[string]ContractDetails)
		var tempContractDetails []ContractDetails
		for _, j := range blockResults.TxsResults {
			if strings.Contains(string(j.Data), "/ibc.core.client.v1.MsgUpdateClient") && strings.Contains(string(j.Data), "/ibc.core.channel.v1.MsgAcknowledgement") {
				for _, k := range j.Events {
					if k.Type == "wasm" && len(k.Attributes) == 3 {
						if string(k.Attributes[0].Key) == "_contract_address" && string(k.Attributes[1].Key) == "lock_id" && string(k.Attributes[2].Key) == "locked_tokens" {
							lockID, err := strconv.ParseInt(string(k.Attributes[1].Value), 10, 64)
							if err != nil {
								return fmt.Errorf("incorrect lock ID at height %d", i)
							}
							//fmt.Println(i, string(k.Attributes[0].Value), lockID, string(k.Attributes[2].Value))
							tempContractDetailsMap[string(k.Attributes[0].Value)] = ContractDetails{
								Address:                 string(k.Attributes[0].Value),
								LockID:                  lockID,
								LockedTokensProtoString: string(k.Attributes[2].Value),
							}
						}
					}
					if k.Type == "wasm" && len(k.Attributes) == 5 {
						if string(k.Attributes[0].Key) == "_contract_address" && string(k.Attributes[1].Key) == "action" &&
							string(k.Attributes[2].Key) == "callback-info" && string(k.Attributes[3].Key) == "reply-msg-id" &&
							string(k.Attributes[4].Key) == "reply-result" {
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
			lockDetailsByHeight = append(lockDetailsByHeight, LockDetailsByHeight{
				Height:          i,
				ContractDetails: tempContractDetails,
			})
		}
	}

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

func ParseMints(RPCAddress string, startingHeight, endHeight int64) error {
	rpcClient, err := http.New(RPCAddress, "/websocket")
	if err != nil {
		return err
	}

	addressToSharesMap := make(map[string]AddressSharesInIncentiveContract)

	for i := startingHeight; i <= endHeight; i++ {
		time.Sleep(time.Millisecond * 50)

		blockResults, err := rpcClient.BlockResults(context.Background(), &i)
		if err != nil {
			return err
		}

		if blockResults.Height == 0 {
			return fmt.Errorf("cannot read height %d", i)
		}

		for _, j := range blockResults.TxsResults {
			if strings.Contains(j.String(), "vault_token_balance") {
				for _, k := range j.Events {
					if k.Type == "wasm" && len(k.Attributes) > 3 {
						if string(k.Attributes[0].Key) == "_contract_address" && string(k.Attributes[1].Key) == "action" &&
							string(k.Attributes[2].Key) == "user" && string(k.Attributes[3].Key) == "vault_token_balance" {
							if len(k.Attributes) > 4 {
								fmt.Println("found a block with multiple mints at height :", i)
							}
							value, ok := addressToSharesMap[string(k.Attributes[2].Value)]
							if ok {
								value.Shares = append(value.Shares, string(k.Attributes[3].Value))
								value.LastUpdatedHeight = append(value.LastUpdatedHeight, i)
								addressToSharesMap[string(k.Attributes[2].Value)] = value
							} else {
								addressToSharesMap[string(k.Attributes[2].Value)] = AddressSharesInIncentiveContract{
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

func CallBackInfos(RPCAddress string, startingHeight, endHeight int64) error {
	rpcClient, err := http.New(RPCAddress, "/websocket")
	if err != nil {
		return err
	}

	var callBackInfoWithHeight []CallBackInfoWithHeight
	for i := startingHeight; i <= endHeight; i++ {
		time.Sleep(time.Millisecond * 20)

		blockResults, err := rpcClient.BlockResults(context.Background(), &i)
		if err != nil {
			return err
		}

		if blockResults.Height == 0 {
			return fmt.Errorf("cannot read height %d", i)
		}

		var tempCallBackInfos []CallBackInfo
		for _, j := range blockResults.TxsResults {
			if strings.Contains(j.String(), "reply-result") && strings.Contains(j.String(), "callback-info") {
				for _, k := range j.Events {
					if k.Type == "wasm" && len(k.Attributes) == 5 {
						if string(k.Attributes[0].Key) == "_contract_address" && string(k.Attributes[1].Key) == "action" &&
							string(k.Attributes[2].Key) == "callback-info" && string(k.Attributes[3].Key) == "reply-msg-id" &&
							string(k.Attributes[4].Key) == "reply-result" {
							tempCallBackInfos = append(tempCallBackInfos, CallBackInfo{
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
			callBackInfoWithHeight = append(callBackInfoWithHeight, CallBackInfoWithHeight{
				Height:        i,
				CallBackInfos: tempCallBackInfos,
			})
		}
	}

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
