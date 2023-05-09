package cmd

import (
	"github.com/arhamchordia/chain-details/cmd/bigquery"
	"github.com/arhamchordia/chain-details/cmd/bigquery/vaults"
	"github.com/spf13/cobra"
)

func BigQueryRegisterRawQueryCmd(parentCmd *cobra.Command) {
	// RawQuery
	bigquery.RawQueryCmd.Flags().StringVarP(&bigquery.RawQuery, "query", "q", "", "SQL query to execute against BigQuery (required)")
	err := bigquery.RawQueryCmd.MarkFlagRequired("query")
	if err != nil {
		return
	}
	parentCmd.AddCommand(bigquery.RawQueryCmd)
}

func BigQueryRegisterTransactionsCmd(parentCmd *cobra.Command) {
	// TransactionsQuery
	bigquery.TransactionsCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Address to query (required)")
	err := bigquery.TransactionsCmd.MarkFlagRequired("address")
	if err != nil {
		return
	}
	parentCmd.AddCommand(bigquery.TransactionsCmd)
}

func BigQueryRegisterVaultsCmd(parentCmd *cobra.Command) {
	// VaultBondQuery
	vaults.BondCmd.Flags().StringVarP(&vaults.AddressQuery, "address", "a", "", "Address to query (optional)")
	vaults.BondCmd.Flags().BoolVarP(&vaults.ConfirmedQuery, "confirmed", "c", false, "Filter by confirmed bond actions (optional)")
	parentCmd.AddCommand(vaults.BondCmd)
	// VaultUnbondQuery
	vaults.UnbondCmd.Flags().StringVarP(&vaults.AddressQuery, "address", "a", "", "Address to query (optional)")
	parentCmd.AddCommand(vaults.UnbondCmd)
	// VaultWithdrawQuery
	vaults.WithdrawCmd.Flags().StringVarP(&vaults.AddressQuery, "address", "a", "", "Address to query (optional)")
	parentCmd.AddCommand(vaults.WithdrawCmd)
}
