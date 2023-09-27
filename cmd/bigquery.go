package cmd

import (
	"fmt"
	"github.com/arhamchordia/chain-details/cmd/bigquery"
	"github.com/spf13/cobra"
)

// bigQueryCmd represents the bigquery command
var bigQueryCmd = &cobra.Command{
	Use:   "bigquery",
	Short: "Commands related to BigQuery queries",
	Long:  `This command is the parent command for all BigQuery related subcommands. Use this command with the appropriate subcommand to execute BigQuery queries.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if bigquery.OutputFormat != "csv" && bigquery.OutputFormat != "json" {
			return fmt.Errorf("invalid output format: %s. Please use 'csv' or 'json'", bigquery.OutputFormat)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(bigQueryCmd)

	bigQueryCmd.PersistentFlags().StringVarP(&bigquery.OutputFormat, "output", "o", "csv", "Output format for generated files (csv/json)")

	BigQueryRegisterRawQueryCmd(bigQueryCmd)
	BigQueryRegisterVaultsCmd(bigQueryCmd)
	BigQueryRegisterTransactionsCmd(bigQueryCmd)
	BigQueryRegisterReportCmd(bigQueryCmd)
}
