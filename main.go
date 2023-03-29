package main

import (
	"fmt"
	
	"github.com/arhamchordia/chain-details/cmd"
	"github.com/arhamchordia/chain-details/internal"
)

func main() {
	err := internal.ReplayChain("https://quasar-rpc.lavenderfive.com:443")
	fmt.Println(err)

	cmd.Execute()
}
