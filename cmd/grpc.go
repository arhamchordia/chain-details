package cmd

import (
	"fmt"
	"github.com/arhamchordia/chain-details/cmd/config"
	"github.com/spf13/cobra"
)

// GrpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Commands related to gRPC queries",
	Long:  `This command is the parent command for all gRPC related subcommands. Use this command with the appropriate subcommand to execute gRPC queries.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if config.OutputFormat != "csv" && config.OutputFormat != "json" {
			return fmt.Errorf("invalid output format: %s. Please use 'csv' or 'json'", config.OutputFormat)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)

	grpcCmd.PersistentFlags().StringVarP(&config.OutputFormat, "output", "o", "csv", "Output format for generated files (csv/json)")

	GRPCRegisterDelegatorsCmd(grpcCmd)
	GRPCRegisterDepositorsCmd(grpcCmd)
	GRPCRegisterGenesisCmd(grpcCmd)
	GRPCRegisterValidatorsCmd(grpcCmd)
}
