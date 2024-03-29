package vaults

import (
	"github.com/arhamchordia/chain-details/cmd/config"
	"github.com/arhamchordia/chain-details/internal/bigquery/vaults"
	"github.com/spf13/cobra"
)

var BondCmd = &cobra.Command{
	Use:   "bond",
	Short: "Generates csv files with wasm vault bond data",
	Long:  `This command allows you to query all the bond actions to the Quasar Vault. Provide optionally the address to filter for with the --address flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.QueryBond(config.AddressQuery, config.ConfirmedQuery, config.PendingQuery, config.OutputFormat)
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
		err := vaults.QueryUnbond(config.AddressQuery, config.ConfirmedQuery, config.PendingQuery, config.OutputFormat)
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
		err := vaults.QueryClaim(config.AddressQuery, config.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}
