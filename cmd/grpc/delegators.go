package grpc

import (
	"github.com/arhamchordia/chain-details/internal/grpc"
	"github.com/spf13/cobra"
)

var DelegatorsDataCmd = &cobra.Command{
	Use:   "delegators-data [grpc-url]",
	Short: "Generates csv files with delegators data",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		grpcUrl := args[0]

		err := grpc.QueryDelegatorsData(grpcUrl)
		if err != nil {
			return err
		}
		return nil
	},
}
