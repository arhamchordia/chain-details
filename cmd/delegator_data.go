package cmd

import (
	"github.com/arhamchrodia/validator-status/internal"
	"github.com/spf13/cobra"
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
