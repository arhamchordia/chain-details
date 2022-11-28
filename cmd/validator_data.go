package cmd

import (
	"github.com/spf13/cobra"

	"github.com/arhamchordia/chain-details/internal"
)

var parseValidatorsCmd = &cobra.Command{
	Use:   "validators-data [grpc-url] [account-address-prefix]",
	Short: "Generates csv file with validators data",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		grpcUrl := args[0]
		accountPrefix := args[1]

		err := internal.QueryValidatorsData(grpcUrl, accountPrefix)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseValidatorsCmd)
}
