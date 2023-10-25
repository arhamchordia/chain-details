package grpc

import (
	"github.com/arhamchordia/chain-details/internal/grpc"
	"github.com/spf13/cobra"
	"strconv"
)

var DepositorsBondCmd = &cobra.Command{
	Use:   "depositors-bond [rpc-url] [start-height] [end-height]",
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

		err = grpc.QueryDepositorsBond(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}

var DepositorsUnbondCmd = &cobra.Command{
	Use:   "depositors-unbond [rpc-url] [start-height] [end-height]",
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

		err = grpc.QueryDepositorsUnbond(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}

var DepositorsLockedTokensCmd = &cobra.Command{
	Use:   "depositors-locked-tokens [rpc-url] [start-height] [end-height]",
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

		err = grpc.QueryDepositorsLockedTokens(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}

var DepositorsMintsCmd = &cobra.Command{
	Use:   "depositors-mints [rpc-url] [start-height] [end-height]",
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

		err = grpc.QueryDepositorsMints(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}

var DepositorsCallbackInfoCmd = &cobra.Command{
	Use:   "depositors-callback-info [rpc-url] [start-height] [end-height]",
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

		err = grpc.QueryDepositorsCallbackInfo(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}

var DepositorsBeginUnlockingCmd = &cobra.Command{
	Use:   "depositors-begin-unlocking [rpc-url] [start-height] [end-height]",
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

		err = grpc.QueryDepositorsBeginUnlocking(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}

// DepositorsReplayChainCmd TODO: This looks out of context here. If it is not related to depositors please move it away
var DepositorsReplayChainCmd = &cobra.Command{
	Use:   "depositors-replay-chain [rpc-url] start-height end-height",
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

		err = grpc.QueryDepositorsReplayChain(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}

var BlockSignatureOfValidators = &cobra.Command{
	Use:   "signer-counter [rpc-url] [start-height] [end-height]",
	Short: "Queries blocks signed by the validators",
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

		err = grpc.QueryBlocksSignerCounter(rpcURL, startingHeight, endHeight)
		if err != nil {
			return err
		}
		return nil
	},
}
