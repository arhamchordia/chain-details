package cmd

import (
	"github.com/arhamchordia/chain-details/cmd/grpc"
	"github.com/spf13/cobra"
)

func GRPCRegisterDelegatorsCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(grpc.DelegatorsDataCmd)
}

func GRPCRegisterDepositorsCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(grpc.DepositorsBondCmd)
	parentCmd.AddCommand(grpc.DepositorsUnbondCmd)
	parentCmd.AddCommand(grpc.DepositorsLockedTokensCmd)
	parentCmd.AddCommand(grpc.DepositorsMintsCmd)
	parentCmd.AddCommand(grpc.DepositorsCallbackInfoCmd)
	parentCmd.AddCommand(grpc.DepositorsBeginUnlockingCmd)
	parentCmd.AddCommand(grpc.DepositorsReplayChainCmd)
}

func GRPCRegisterGenesisCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(grpc.GenesisVestingAccountsCmd)
}

func GRPCRegisterValidatorsCmd(parentCmd *cobra.Command) {
	parentCmd.AddCommand(grpc.ValidatorsDataCmd)
	parentCmd.AddCommand(grpc.GenesisAndPostGenesisValidatorsDataCmd)
}
