package vaults

import (
	"github.com/arhamchordia/chain-details/internal/bigquery/vaults"
	"github.com/spf13/cobra"
)

// TODO: Add flag for dynamic VaultAddress
var BondCmd = &cobra.Command{
	Use:   "bond",
	Short: "Generates csv files with wasm vault bond data",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.QueryBond()
		if err != nil {
			return err
		}
		return nil
	},
}

// TODO: Add flag for dynamic VaultAddress
var UnbondCmd = &cobra.Command{
	Use:   "unbond",
	Short: "Generates csv files with wasm vault unbond data",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.QueryUnbond()
		if err != nil {
			return err
		}
		return nil
	},
}

// TODO: Add flag for dynamic VaultAddress
var WithdrawCmd = &cobra.Command{
	Use:   "withdraw",
	Short: "Generates csv files with wasm vault withdraw data",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.QueryUnbond()
		if err != nil {
			return err
		}
		return nil
	},
}
