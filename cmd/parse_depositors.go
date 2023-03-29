package cmd

import (
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	"github.com/spf13/cobra"
	"strconv"
)

var parseDepositorsCmd = &cobra.Command{
	Use:   "parse-depositors [rpc-url] start-height end-height",
	Short: "Queries data for the people depositing to contracts",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		rpcURL := args[0]
		startingHeight, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return err
		}
		endHeight, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			return err
		}

		depositorDetails, err := internal.ReplayChain(rpcURL, startingHeight, endHeight)
		fmt.Println(depositorDetails)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseDepositorsCmd)
}
