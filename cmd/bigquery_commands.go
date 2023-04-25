package cmd

import (
	"github.com/arhamchordia/chain-details/cmd/bigquery"
	"github.com/spf13/cobra"
)

func RegisterSampleCommandsBigQuery(parentCmd *cobra.Command) {
	bigquery.SampleCmd.Flags().StringVarP(&bigquery.SampleQuery, "query", "q", "", "SQL query to execute against BigQuery (required)")
	err := bigquery.SampleCmd.MarkFlagRequired("query")
	if err != nil {
		return
	}
	parentCmd.AddCommand(bigquery.SampleCmd)
}

func RegisterDelegatorsCommandsBigQuery(parentCmd *cobra.Command) {
	parentCmd.AddCommand(bigquery.DelegatorsDataCmd)
}
