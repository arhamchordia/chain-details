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

var parseLockedTokensCmd = &cobra.Command{
	Use:   "parse-locked-tokens [rpc-url] start-height end-height",
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

		err = internal.CheckLockedTokens(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}

var parseMintsCmd = &cobra.Command{
	Use:   "parse-mints [rpc-url] start-height end-height",
	Short: "Queries data for the people received minted shares",
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

		err = internal.ParseMints(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}

var callBackInfosCmd = &cobra.Command{
	Use:   "callback-infos [rpc-url] start-height end-height",
	Short: "Queries data for the callback infos",
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

		err = internal.CallBackInfos(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}

var beginUnlockingCmd = &cobra.Command{
	Use:   "begin-unlocking [rpc-url] start-height end-height",
	Short: "Queries data for begin unlocking",
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

		err = internal.BeginUnlocking(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}

var parseAllDataCmd = &cobra.Command{
	Use:   "parse-all [rpc-url] start-height end-height",
	Short: "Queries data for all kinds",
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

		err = internal.ReplayChain(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(parseDepositorsBondCmd)
	rootCmd.AddCommand(parseDepositorsUnbondCmd)
	rootCmd.AddCommand(parseLockedTokensCmd)
	rootCmd.AddCommand(parseMintsCmd)
	rootCmd.AddCommand(callBackInfosCmd)
	rootCmd.AddCommand(beginUnlockingCmd)
	rootCmd.AddCommand(parseAllDataCmd)
}
