package vaults

import (
	"github.com/arhamchordia/chain-details/cmd/bigquery"
	"github.com/arhamchordia/chain-details/internal/bigquery/vaults"
	"github.com/spf13/cobra"
)

var BondCmd = &cobra.Command{
	Use:   "bond",
	Short: "Generates csv files with wasm vault bond data",
	Long:  `This command allows you to query all the bond actions to the Quasar Vault. Provide optionally the address to filter for with the --address flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.QueryBond(bigquery.AddressQuery, bigquery.ConfirmedQuery, bigquery.PendingQuery, bigquery.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}

var UnbondCmd = &cobra.Command{
	Use:   "unbond",
	Short: "Generates csv files with wasm vault unbond data",
	Long:  `This command allows you to query all the unbond actions to the Quasar Vault. Provide optionally the address to filter for with the --address flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.QueryUnbond(bigquery.AddressQuery, bigquery.ConfirmedQuery, bigquery.PendingQuery, bigquery.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}

var ClaimCmd = &cobra.Command{
	Use:   "claim",
	Short: "Generates csv files with wasm vault claim data",
	Long:  `This command allows you to query all the claim actions to the Quasar Vault. Provide optionally the address to filter for with the --address flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.QueryClaim(bigquery.AddressQuery, bigquery.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}
