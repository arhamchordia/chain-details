package cmd

import (
	"github.com/arhamchordia/chain-details/internal"
	"github.com/spf13/cobra"
)

var parseDepositorsCmd = &cobra.Command{
	Use:   "parse-depositors [rpc-url]",
	Short: "Queries data for the people depositing to contracts",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rpcURL := args[0]

		err := internal.ReplayChain(rpcURL)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseDepositorsCmd)
}
