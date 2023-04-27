package bigquery

import (
	"github.com/arhamchordia/chain-details/internal/bigquery"
	"github.com/spf13/cobra"
)

// TODO: Add flag for dynamic VaultAddress
var DepositorsBondCmd = &cobra.Command{
	Use:   "depositors-bond",
	Short: "Generates csv files with wasm vault bond data",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := bigquery.QueryDepositorsBond()
		if err != nil {
			return err
		}
		return nil
	},
}

// TODO: Add flag for dynamic VaultAddress
var DepositorsUnbondCmd = &cobra.Command{
	Use:   "depositors-unbond",
	Short: "Generates csv files with wasm vault unbond data",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := bigquery.QueryDepositorsUnbond()
		if err != nil {
			return err
		}
		return nil
	},
}
