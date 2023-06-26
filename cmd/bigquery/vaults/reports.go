package vaults

import (
	"github.com/arhamchordia/chain-details/cmd/config"
	"github.com/arhamchordia/chain-details/internal/bigquery/vaults"
	"github.com/spf13/cobra"
)

var ReportBondCmd = &cobra.Command{
	Use:   "report-bond",
	Short: "TODO",
	Long:  `TODO TODO`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.QueryDailyReportBond(config.AddressQuery, config.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}
