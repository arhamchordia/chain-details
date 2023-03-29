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

type DepositorDetails struct {
	Address      string `json:"address"`
	BlockHeight  int64  `json:"block_height"`
	Amount       string `json:"amount"`
	VaultAddress string `json:"primitive_address"`
	//Action       string `json:"action"`
	//BurntShares  string `json:"burnt_shares"` //TODO incorporate unbonding
	BondID string `json:"bond_id"`
}

func ReplayChain(RPCAddress string, startingHeight, endHeight int64) error {
	rpcClient, err := http.New(RPCAddress, "/websocket")
	if err != nil {
		return err
	}

	var depositorDetails []DepositorDetails
	for i := startingHeight; i <= endHeight; i++ {
		time.Sleep(time.Millisecond * 10)

		blockResults, err := rpcClient.BlockResults(context.Background(), &i)
		if err != nil {
			return err
		}

		if blockResults.Height == 0 {
			return fmt.Errorf("cannot read height %s", i)
		}

		for _, j := range blockResults.TxsResults {
			var tempDepositorDetails []DepositorDetails
			var tempBondIDs []string
			if strings.Contains(string(j.Data), "/cosmwasm.wasm.v1.MsgExecuteContract") {
				for o, k := range j.Events {
					if k.Type == "message" && string(k.Attributes[0].Value) == "/cosmwasm.wasm.v1.MsgExecuteContract" {
						if len(j.Events) >= o+3 {
							if j.Events[o+1].Type == "message" && string(j.Events[o+1].Attributes[0].Value) == "wasm" {
								if j.Events[o+2].Type == "coin_spent" {
									if j.Events[o+3].Type == "coin_received" && string(j.Events[o+3].Attributes[0].Value) == "quasar18a2u6az6dzw528rptepfg6n49ak6hdzkf8ewf0n5r0nwju7gtdgqamr7qu" {
										fmt.Println(string(j.Events[o+2].Attributes[0].Value), ":", string(j.Events[o+2].Attributes[1].Value))
										tempDepositorDetails = append(tempDepositorDetails, DepositorDetails{
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
						tempBondIDs = append(tempBondIDs, string(q.Attributes[1].Value))
					}
				}

				if len(tempBondIDs) != len(tempDepositorDetails) {
					fmt.Println(tempBondIDs)
					fmt.Println(tempDepositorDetails)
					return fmt.Errorf("mismatch in the counting of bond IDs %s", i)
				}

				for p := range tempBondIDs {
					depositorDetails = append(depositorDetails, DepositorDetails{
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
