package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "chain-details",
	Short: "Chain details is a simple CLI to generate specific details for any Cosmos chain in a CSV or JSON file",
	Long:  `Chain details is a command-line interface that provides various commands to query specific details from Cosmos chains. The CLI supports querying data from both gRPC endpoints and Google Cloud BigQuery.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Please provide a valid subcommand. Use 'chain-details --help' for more information.")
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Whoops. There was an error while executing your CLI '%s'", err)
	}
}
