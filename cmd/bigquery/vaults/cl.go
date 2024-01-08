package vaults

import (
	"github.com/arhamchordia/chain-details/cmd/bigquery"
	"github.com/arhamchordia/chain-details/internal/bigquery/vaults"
	"github.com/spf13/cobra"
)

// TODO update descriptions

var CLDepositCmd = &cobra.Command{
	Use:   "cl-deposit",
	Short: "Generates csv files with wasm vault deposit data",
	Long:  `This command allows you to query all the deposit actions to the Quasar Vault. Provide optionally the address to filter for with the --address flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.CLQueryDeposit(bigquery.AddressQuery, bigquery.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}

var CLWithdrawCmd = &cobra.Command{
	Use:   "cl-withdraw",
	Short: "Generates csv files with wasm vault withdraw data",
	Long:  `This command allows you to query all the withdraw actions to the Quasar Vault. Provide optionally the address to filter for with the --address flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.CLQueryWithdraw(bigquery.AddressQuery, bigquery.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}

var CLClaimCmd = &cobra.Command{
	Use:   "cl-claim",
	Short: "Generates csv files with wasm vault claim data",
	Long:  `This command allows you to query all the claim actions to the Quasar Vault. Provide optionally the address to filter for with the --address flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.CLQueryClaim(bigquery.AddressQuery, bigquery.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}

var CLDistributeRewardsCmd = &cobra.Command{
	Use:   "cl-distribute-rewards",
	Short: "Generates csv files with wasm vault distribute rewards data",
	Long:  `This command allows you to query all the distribute rewards actions to the Quasar Vault. Provide optionally the address to filter for with the --address flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.CLQueryDistributeRewards(bigquery.AddressQuery, bigquery.DaysInterval, bigquery.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}

var CLAPRCmd = &cobra.Command{
	Use:   "cl-apr",
	Short: "Generates csv files with wasm vault average APR data",
	Long:  `This command allows you to compute the APR for a specific Quasar CL Vault. Provide optionally the days interval with the --days flag.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.CLQueryAPR(bigquery.AddressQuery, bigquery.DaysInterval, bigquery.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}

var CLReportCmd = &cobra.Command{
	Use:   "cl-report",
	Short: "Generates a comprehensive report on user activity for a given vault.",
	Long:  `Generates a comprehensive report on user activity within the last 24 hours, as well as general activity since the start of the vault. The report includes information on new and old user bonds, unbonds, exits, total bonded and active users, pending unbond amounts, and averages. It also features a 'Wall of Fame' section highlighting the users with the biggest deposits and holdings.`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := vaults.CLQueryReport(bigquery.AddressQuery, bigquery.OutputFormat)
		if err != nil {
			return err
		}
		return nil
	},
}
