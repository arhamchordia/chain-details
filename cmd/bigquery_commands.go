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

func BigQueryRegisterLPVaultsCmd(parentCmd *cobra.Command) {
	// LPVaultBondQuery
	vaults.LPBondCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Address to query (optional)")
	vaults.LPBondCmd.Flags().BoolVarP(&bigquery.ConfirmedQuery, "confirmed", "c", false, "Filter by confirmed bond actions (optional)")
	vaults.LPBondCmd.Flags().BoolVarP(&bigquery.PendingQuery, "pending", "p", false, "Filter by pending bond actions (optional)")
	parentCmd.AddCommand(vaults.LPBondCmd)

	// LPVaultUnbondQuery
	vaults.LPUnbondCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Address to query (optional)")
	vaults.LPUnbondCmd.Flags().BoolVarP(&bigquery.ConfirmedQuery, "confirmed", "c", false, "Filter by confirmed unbond actions (optional)")
	vaults.LPUnbondCmd.Flags().BoolVarP(&bigquery.PendingQuery, "pending", "p", false, "Filter by pending unbond actions (optional)")
	parentCmd.AddCommand(vaults.LPUnbondCmd)

	// LPVaultClaimQuery
	vaults.LPClaimCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Address to query (optional)")
	parentCmd.AddCommand(vaults.LPClaimCmd)

	// LPVaultReportCmd
	vaults.LPReportCmd.Flags().IntVarP(&bigquery.BlockHeight, "block", "b", 1, "Block height to query from")
	vaults.LPReportCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Vault address to query")
	parentCmd.AddCommand(vaults.LPReportCmd)
}

func BigQueryRegisterCLVaultsCmd(parentCmd *cobra.Command) {
	// CLVaultDepositCmd
	vaults.CLDepositCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Vault address to query")
	parentCmd.AddCommand(vaults.CLDepositCmd)

	// CLVaultWithdrawCmd
	vaults.CLWithdrawCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Vault address to query")
	parentCmd.AddCommand(vaults.CLWithdrawCmd)

	// CLVaultClaimCmd
	vaults.CLClaimCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Vault address to query")
	parentCmd.AddCommand(vaults.CLClaimCmd)

	// TODO: CLVaultDistributeRewardsCmd
	vaults.CLDistributeRewardsCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Vault address to query")
	vaults.CLDistributeRewardsCmd.Flags().IntVarP(&bigquery.DaysInterval, "days", "d", 1, "Days to take in account exporting the rewards (default 1)")
	parentCmd.AddCommand(vaults.CLDistributeRewardsCmd)

	// TODO: CLVaultAPRCmd
	vaults.CLAPRCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Vault address to query")
	vaults.CLAPRCmd.Flags().IntVarP(&bigquery.DaysInterval, "days", "d", 1, "Days to take in account computing the avg APR (default 1)")
	parentCmd.AddCommand(vaults.CLAPRCmd)

	// CLVaultReportCmd
	vaults.CLReportCmd.Flags().StringVarP(&bigquery.AddressQuery, "address", "a", "", "Vault address to query")
	parentCmd.AddCommand(vaults.CLReportCmd)
}
