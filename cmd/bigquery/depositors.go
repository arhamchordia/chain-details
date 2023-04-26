package bigquery

import (
	"github.com/arhamchordia/chain-details/internal/bigquery"
	"github.com/spf13/cobra"
)

var DepositorsBondCmd = &cobra.Command{
	Use:   "depositors-bond",
	Short: "Generates csv files with delegators data",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := bigquery.QueryDepositorsBond()
		if err != nil {
			return err
		}
		return nil
	},
}
