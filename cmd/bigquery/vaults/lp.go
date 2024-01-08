package vaults

import (
	"github.com/arhamchordia/chain-details/cmd/bigquery"
	"github.com/arhamchordia/chain-details/internal/bigquery/vaults"
	"github.com/spf13/cobra"
)

var LPBondCmd = &cobra.Command{
	Use:   "lp-bond",
	Short: "Generates csv files with wasm vault bond data",
	Long:  `This command allows you to query all the bond actions to the Quasar Vault. Provide optionally the address to filter for with the --address flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.LPQueryBond(bigquery.AddressQuery, bigquery.ConfirmedQuery, bigquery.PendingQuery, bigquery.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}

var LPUnbondCmd = &cobra.Command{
	Use:   "lp-unbond",
	Short: "Generates csv files with wasm vault unbond data",
	Long:  `This command allows you to query all the unbond actions to the Quasar Vault. Provide optionally the address to filter for with the --address flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.LPQueryUnbond(bigquery.AddressQuery, bigquery.ConfirmedQuery, bigquery.PendingQuery, bigquery.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}

var LPClaimCmd = &cobra.Command{
	Use:   "lp-claim",
	Short: "Generates csv files with wasm vault claim data",
	Long:  `This command allows you to query all the claim actions to the Quasar Vault. Provide optionally the address to filter for with the --address flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.LPQueryClaim(bigquery.AddressQuery, bigquery.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}

var LPReportCmd = &cobra.Command{
	Use:   "lp-report",
	Short: "Generates a comprehensive report on user activity for a given vault.",
	Long:  `Generates a comprehensive report on user activity within the last 24 hours, as well as general activity since the start of the vault. The report includes information on new and old user bonds, unbonds, exits, total bonded and active users, pending unbond amounts, and averages. It also features a 'Wall of Fame' section highlighting the users with the biggest deposits and holdings.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.LPQueryReport(bigquery.BlockHeight, bigquery.AddressQuery, bigquery.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}
