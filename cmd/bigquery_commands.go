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
	vaults.BondCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Address to query (optional)")
	vaults.BondCmd.Flags().BoolVarP(&bigquery.ConfirmedQuery, "confirmed", "c", false, "Filter by confirmed bond actions (optional)")
	vaults.BondCmd.Flags().BoolVarP(&bigquery.PendingQuery, "pending", "p", false, "Filter by pending bond actions (optional)")
	parentCmd.AddCommand(vaults.BondCmd)

	// VaultUnbondQuery
	vaults.UnbondCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Address to query (optional)")
	vaults.UnbondCmd.Flags().BoolVarP(&bigquery.ConfirmedQuery, "confirmed", "c", false, "Filter by confirmed unbond actions (optional)")
	vaults.UnbondCmd.Flags().BoolVarP(&bigquery.PendingQuery, "pending", "p", false, "Filter by pending unbond actions (optional)")
	parentCmd.AddCommand(vaults.UnbondCmd)

	// VaultClaimQuery
	vaults.ClaimCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Address to query (optional)")
	parentCmd.AddCommand(vaults.ClaimCmd)
}

func BigQueryRegisterReportCmd(parentCmd *cobra.Command) {
	// ReportBondCmd
	vaults.ReportCmd.Flags().IntVarP(&bigquery.BlockHeight, "block", "b", 1, "Block height to query from")
	vaults.ReportCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Vault address to query")
	parentCmd.AddCommand(vaults.ReportCmd)
}
