package main

import (
	"fmt"
	"github.com/arhamchordia/chain-details/cmd"
	"github.com/arhamchordia/chain-details/internal"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Genesis struct {
	AppState AppState `json:"app_state"`
}

type AppState struct {
	Bank Bank `json:"bank"`
	Auth Auth `json:"auth"`
}

type Bank struct {
	Balances []Balances `json:"balances"`
}

type Balances struct {
	Address string    `json:"address"`
	Coins   sdk.Coins `json:"coins"`
}

type Auth struct {
	Accounts []Accounts `json:"accounts"`
}

type Accounts struct {
	Type          string      `json:"@type"`
	Address       string      `json:"address"`
	PubKey        interface{} `json:"pub_key"`
	AccountNumber string      `json:"account_number"`
	Sequence      string      `json:"sequence"`
}

func main() {
	//jsonFile, err := os.Open("export.json")
	//// if we os.Open returns an error then handle it
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Successfully Opened export.json")
	//// defer the closing of our jsonFile so that we can parse it later on
	//defer jsonFile.Close()
	//
	//byteValue, _ := ioutil.ReadAll(jsonFile)
	//
	//var result Genesis
	//err = json.Unmarshal(byteValue, &result)
	//if err != nil {
	//	panic(err)
	//}
	//
	//var bank Bank
	//var auth Auth
	//
	//for _, res := range result.AppState.Bank.Balances {
	//	var tempBalance Balances
	//	tempBalance.Address = res.Address
	//
	//	found, uayy := res.Coins.Find("uayy")
	//	if found {
	//		tempBalance.Coins = tempBalance.Coins.Add(uayy)
	//	}
	//	found, uoro := res.Coins.Find("uoro")
	//	if found {
	//		tempBalance.Coins = tempBalance.Coins.Add(uoro)
	//	}
	//	found, uqsr := res.Coins.Find("uqsr")
	//	if found {
	//		tempBalance.Coins = tempBalance.Coins.Add(uqsr)
	//	}
	//
	//	bank.Balances = append(bank.Balances, tempBalance)
	//	auth.Accounts = append(auth.Accounts, Accounts{
	//		Type:          "/cosmos.auth.v1beta1.BaseAccount",
	//		Address:       res.Address,
	//		PubKey:        nil,
	//		AccountNumber: "0",
	//		Sequence:      "0",
	//	})
	//}
	//
	//count := 0
	//for _, res := range result.AppState.Auth.Accounts {
	//	if res.Sequence != "0" && res.Sequence != "1" {
	//		count++
	//	}
	//}
	//fmt.Println(count)
	//fmt.Println(len(result.AppState.Bank.Balances))
	//fmt.Println(len(result.AppState.Auth.Accounts))
	//
	//file, _ := json.MarshalIndent(bank, "", " ")
	//_ = ioutil.WriteFile("balances.json", file, 0644)
	//
	//file, _ = json.MarshalIndent(auth, "", " ")
	//_ = ioutil.WriteFile("accounts.json", file, 0644)
	err := internal.ReplayChain("https://quasar-rpc.lavenderfive.com:443")
	fmt.Println(err)

	cmd.Execute()
}
