package bigquery

import (
	"github.com/arhamchordia/chain-details/internal/bigquery"
	"github.com/spf13/cobra"
)

var TransactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "Generates csv files with address transactions list",
	Long:  `This command allows you to execute a SQL query against Google Cloud BigQuery. Provide the address to query with the --address flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := bigquery.TransactionsQuery(AddressQuery, OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}
