package bigquery

import (
	"fmt"
	"google.golang.org/api/iterator"
	"log"

	"cloud.google.com/go/bigquery"
	"github.com/arhamchordia/chain-details/internal"
	"github.com/spf13/cobra"
)

var RawQuery string

// SampleCmd represents the bigquery command
var RawQueryCmd = &cobra.Command{
	Use:   "bigquery",
	Short: "Execute a BigQuery SQL query",
	Long:  `This command allows you to execute a SQL query against Google Cloud BigQuery. Provide the SQL query with the --query flag.`,
	Run: func(cmd *cobra.Command, args []string) {
		bqQuerier, _ := internal.NewBigQueryQuerier()

		it, err := bqQuerier.ExecuteQuery(RawQuery)
		if err != nil {
			log.Fatalf("Failed to execute BigQuery query: %v", err)
		}
		defer bqQuerier.Close()

		for {
			var row []bigquery.Value
			err := it.Next(&row)
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate results: %v", err)
			}

			for _, val := range row {
				fmt.Printf("%v ", val)
			}
			fmt.Println()
		}
	},
}
