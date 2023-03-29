package internal

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tendermint/tendermint/rpc/client/http"
)

func ReplayChain(RPCAddress string) error {
	rpcClient, err := http.New(RPCAddress, "/websocket")
	if err != nil {
		return err
	}

	addressAndAmount := map[string]string{}
	for i := 18400; i < 19424; i++ {
		time.Sleep(time.Second)
		blockHeight := int64(i)

		blockResults, err := rpcClient.BlockResults(context.Background(), &blockHeight)
		if err != nil {
			return err
		}
		fmt.Println(blockHeight)
		for _, j := range blockResults.TxsResults {
			if strings.Contains(string(j.Data), "/cosmwasm.wasm.v1.MsgExecuteContract") {
				for o, k := range j.Events {
					if k.Type == "message" && string(k.Attributes[0].Value) == "/cosmwasm.wasm.v1.MsgExecuteContract" {
						if len(j.Events) >= o+3 {
							if j.Events[o+1].Type == "message" && string(j.Events[o+1].Attributes[0].Value) == "wasm" {
								if j.Events[o+2].Type == "coin_spent" {
									if j.Events[o+3].Type == "coin_received" && string(j.Events[o+3].Attributes[0].Value) == "quasar18a2u6az6dzw528rptepfg6n49ak6hdzkf8ewf0n5r0nwju7gtdgqamr7qu" {
										fmt.Println(string(j.Events[o+2].Attributes[0].Value), ":", string(j.Events[o+2].Attributes[1].Value))
										addressAndAmount[string(j.Events[o+2].Attributes[0].Value)] = string(j.Events[o+2].Attributes[1].Value)
									}
								}
							}
						} else {
							fmt.Println("couldn't find the next 2 events")
						}
					}
				}
			}
		}
	}
	fmt.Println(addressAndAmount)
	return nil
}
