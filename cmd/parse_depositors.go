package cmd

import (
	"github.com/arhamchordia/chain-details/internal"
	"github.com/spf13/cobra"
	"strconv"
)

var parseDepositorsBondCmd = &cobra.Command{
	Use:   "parse-depositors-bond [rpc-url] start-height end-height",
	Short: "Queries data for the people depositing to contracts",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		rpcURL := args[0]
		startingHeight, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return err
		}
		endHeight, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			return err
		}

		err = internal.ReplayChainBond(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}

var parseDepositorsUnbondCmd = &cobra.Command{
	Use:   "parse-depositors-unbond [rpc-url] start-height end-height",
	Short: "Queries data for the people depositing to contracts",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		rpcURL := args[0]
		startingHeight, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return err
		}
		endHeight, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			return err
		}

		err = internal.ReplayChainUnbond(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseDepositorsBondCmd)
	rootCmd.AddCommand(parseDepositorsUnbondCmd)
}
