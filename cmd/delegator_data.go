package cmd

import (
	"github.com/spf13/cobra"

	"github.com/arhamchordia/chain-details/internal"
)

var parseDelegatorsCmd = &cobra.Command{
	Use:   "delegators-data [grpc-url]",
	Short: "Generates csv files with delegators data",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		grpcUrl := args[0]

		err := internal.QueryDelegatorsData(grpcUrl)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseDelegatorsCmd)
}
