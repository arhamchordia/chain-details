package cmd

import (
	"github.com/arhamchordia/chain-details/cmd/bigquery"
	"github.com/arhamchordia/chain-details/cmd/bigquery/vaults"
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

func BigQueryRegisterTransactionsCmd(parentCmd *cobra.Command) {
	bigquery.RawQueryCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Address to query (required)")
	err := bigquery.RawQueryCmd.MarkFlagRequired("query")
	if err != nil {
		return
	}
	parentCmd.AddCommand(bigquery.TransactionsCmd)
}

func BigQueryRegisterVaultsCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(vaults.BondCmd)
	parentCmd.AddCommand(vaults.UnbondCmd)
	parentCmd.AddCommand(vaults.WithdrawCmd)
}
