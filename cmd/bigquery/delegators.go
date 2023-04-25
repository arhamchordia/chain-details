package bigquery

import (
	internalbigquery "github.com/arhamchordia/chain-details/internal/bigquery"
	"github.com/spf13/cobra"
)

var DelegatorsDataCmd = &cobra.Command{
	Use:   "delegators-data",
	Short: "Generates csv files with delegators data",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := internalbigquery.QueryDelegatorsData()
		if err != nil {
			return err
		}
		return nil
	},
}
