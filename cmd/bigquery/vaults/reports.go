package vaults

import (
	"github.com/arhamchordia/chain-details/cmd/config"
	"github.com/arhamchordia/chain-details/internal/bigquery/vaults"
	"github.com/spf13/cobra"
)

var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "TODO",
	Long:  `TODO TODO`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.QueryDailyReport(config.AddressQuery, config.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}
