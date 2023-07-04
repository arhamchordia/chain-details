package grpc

import (
	"github.com/arhamchordia/chain-details/internal/grpc"
	"github.com/spf13/cobra"
)

var GenesisVestingAccountsCmd = &cobra.Command{
	Use:   "genesis-vesting-accounts [json-url] [denom]",
	Short: "Generates csv file with accounts analysis in genesis, information on all the vesting accounts.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		jsonURL := args[0]
		denom := args[1]

		err := grpc.QueryGenesisJSON(jsonURL, denom)
		if err != nil {
			return err
		}
		return nil
	},
}
