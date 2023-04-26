package cmd

import (
	"github.com/arhamchordia/chain-details/cmd/bigquery"
	"github.com/spf13/cobra"
)

func BigQueryRegisterRawQueryCmd(parentCmd *cobra.Command) {
	bigquery.RawQueryCmd.Flags().StringVarP(&bigquery.RawQuery, "query", "q", "", "SQL query to execute against BigQuery (required)")
	err := bigquery.RawQueryCmd.MarkFlagRequired("query")
	if err != nil {
		return
	}
	parentCmd.AddCommand(bigquery.RawQueryCmd)
}

func BigQueryRegisterDelegatorsCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(bigquery.DelegatorsDataCmd)
}

func BigQueryRegisterDepositorsCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(bigquery.DepositorsBondCmd)
}
