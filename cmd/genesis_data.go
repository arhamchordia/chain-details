package cmd

import (
	"github.com/spf13/cobra"

	"github.com/arhamchordia/chain-details/internal"
)

var parseVestingAccountsJSONCmd = &cobra.Command{
	Use:   "vesting-accounts [json-url] [denom]",
	Short: "Generates csv file with accounts analysis in genesis, information on all the vesting accounts.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonURL := args[0]
		denom := args[1]

		err := internal.QueryGenesisJSON(jsonURL, denom)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseVestingAccountsJSONCmd)
}
