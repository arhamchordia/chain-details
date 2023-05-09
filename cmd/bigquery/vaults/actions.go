package vaults

import (
	"github.com/arhamchordia/chain-details/internal/bigquery/vaults"
	"github.com/spf13/cobra"
)

var AddressQuery string
var ConfirmedQuery bool
var PendingQuery bool

var BondCmd = &cobra.Command{
	Use:   "bond",
	Short: "Generates csv files with wasm vault bond data",
	Long:  `This command allows you to query all the bond actions to the Quasar Vault. Provide optionally the address to filter for with the --address flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.QueryBond(AddressQuery, ConfirmedQuery, PendingQuery)
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
		err := vaults.QueryUnbond(AddressQuery)
		if err != nil {
			return err
		}
		return nil
	},
}

var WithdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "Generates csv files with wasm vault withdraw data",
	Long:  `This command allows you to query all the withdraw actions to the Quasar Vault. Provide optionally the address to filter for with the --address flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.QueryWithdraw(AddressQuery)
		if err != nil {
			return err
		}
		return nil
	},
}
