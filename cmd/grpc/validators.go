package grpc

import (
	"github.com/arhamchordia/chain-details/internal/grpc"
	"github.com/spf13/cobra"
)

var ValidatorsDataCmd = &cobra.Command{
	Use:   "validators-data [grpc-url] [account-address-prefix]",
	Short: "Generates csv file with validators data",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		grpcUrl := args[0]
		accountPrefix := args[1]

		err := grpc.QueryValidatorsData(grpcUrl, accountPrefix)
		if err != nil {
			return err
		}
		return nil
	},
}
