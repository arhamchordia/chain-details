package vaults

import (
	"fmt"
	"github.com/arhamchordia/chain-details/cmd/config"
	"github.com/arhamchordia/chain-details/internal/bigquery/vaults"
	"github.com/spf13/cobra"
	"time"
)

var ReportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generates a comprehensive daily report on user activity for a given vault.",
	Long:  `Generates a comprehensive daily report on user activity within the last 24 hours, as well as general activity since the start of the vault. The report includes information on new and old user bonds, unbonds, exits, total bonded and active users, pending unbond amounts, and averages. It also features a 'Wall of Fame' section highlighting the users with the biggest deposits and holdings.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var endDate *time.Time

		// Parse end date if provided, else set default to current time
		endDateFlag, _ := cmd.Flags().GetString("end-date")
		if endDateFlag != "" {
			parsedEndDate, err := time.Parse("2006-01-02", endDateFlag)
			if err != nil {
				return fmt.Errorf("invalid end date format: %w", err)
			}
			endDate = &parsedEndDate
		} else {
			defaultEndDate := time.Now()
			endDate = &defaultEndDate
		}

		// Pass the date pointers to the QueryDailyReport function
		err = vaults.QueryDailyReport(config.BlockHeight, config.AddressQuery, endDate, config.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}
