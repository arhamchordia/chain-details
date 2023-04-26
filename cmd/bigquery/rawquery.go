package bigquery

import (
	"github.com/arhamchordia/chain-details/internal/bigquery"
	"github.com/spf13/cobra"
)

var RawQuery string

var RawQueryCmd = &cobra.Command{
	Use:   "bigquery",
	Short: "Execute a BigQuery SQL query",
	Long:  `This command allows you to execute a SQL query against Google Cloud BigQuery. Provide the SQL query with the --query flag.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := bigquery.RawQuery(RawQuery)
		if err != nil {
			return err
		}
		return nil
	},
}
