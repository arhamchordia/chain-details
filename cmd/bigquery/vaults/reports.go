package vaults

import (
	"github.com/arhamchordia/chain-details/cmd/config"
	"github.com/arhamchordia/chain-details/internal/bigquery/vaults"
	"github.com/spf13/cobra"
)

var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generates a comprehensive daily report on user activity for a given vault.",
	Long:  `Generates a comprehensive daily report on user activity within the last 24 hours, as well as general activity since the start of the vault. The report includes information on new and old user bonds, unbonds, exits, total bonded and active users, pending unbond amounts, and averages. It also features a 'Wall of Fame' section highlighting the users with the biggest deposits and holdings.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.QueryDailyReport(config.BlockHeight, config.AddressQuery, config.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}
