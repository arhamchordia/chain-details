package cmd

import (
	"github.com/arhamchordia/chain-details/cmd/grpc"
	"github.com/spf13/cobra"
)

func RegisterDelegatorsCommands(parentCmd *cobra.Command) {
	parentCmd.AddCommand(grpc.DelegatorsDataCmd)
}

func RegisterDepositorsCommands(parentCmd *cobra.Command) {
	parentCmd.AddCommand(grpc.DepositorsBondCmd)
	parentCmd.AddCommand(grpc.DepositorsUnbondCmd)
	parentCmd.AddCommand(grpc.DepositorsLockedTokensCmd)
	parentCmd.AddCommand(grpc.DepositorsMintsCmd)
	parentCmd.AddCommand(grpc.DepositorsCallbackInfoCmd)
	parentCmd.AddCommand(grpc.DepositorsBeginUnlockingCmd)
}

func RegisterGenesisCommands(parentCmd *cobra.Command) {
	parentCmd.AddCommand(grpc.GenesisVestingAccountsCmd)
}

func RegisterValidatorsCommands(parentCmd *cobra.Command) {
	parentCmd.AddCommand(grpc.ValidatorsDataCmd)
}
