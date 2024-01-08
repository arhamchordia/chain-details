package grpc

import (
	"context"
	"crypto/tls"
	"github.com/arhamchordia/chain-details/internal/export"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"

	grpctypes "github.com/arhamchordia/chain-details/types/grpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func QueryDelegationAnalysisData(grpcUrl, address, denom string) error {
	// initialise config for grpc connection
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	// Create a connection to the gRPC server.
	grpcConn, err := grpc.Dial(
		grpcUrl,
		grpc.WithTransportCredentials(credentials.NewTLS(config)),
	)
	if err != nil {
		return err
	}
	defer grpcConn.Close()

	// send a query only when connection state is ready
	for {
		// wait for 4 milliseconds for grpc to connect
		time.Sleep(4 * time.Millisecond)

		if grpcConn.GetState().String() == "READY" {
			err = ParseDelegationAnalysis(grpcConn, address, denom)
			if err != nil {
				return err
			}
			break
		} else if grpcConn.GetState().String() == "TRANSIENT_FAILURE" {
			break
		}
	}

	return nil
}

// ParseDelegationAnalysis parses all the details of delegations of the account address
// on a given chain
func ParseDelegationAnalysis(grpcConn *grpc.ClientConn, address, denom string) error {
	// create a staking client in order to query validators list
	stakingClient := stakingtypes.NewQueryClient(grpcConn)

	// query the validators list using stakingClient
	delegationsResponse, err := stakingClient.DelegatorDelegations(
		context.Background(),
		&stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: address,
			Pagination:    &query.PageRequest{Limit: grpctypes.DelegatorsLimit},
		})
	if err != nil {
		return err
	}

	type delegationAnalysis struct {
		ValidatorAddress string
		DelegatedAmount  sdk.Coin
	}
	var list []delegationAnalysis
	total := sdk.NewCoin(denom, sdk.ZeroInt())
	for _, j := range delegationsResponse.DelegationResponses {
		list = append(list, delegationAnalysis{
			DelegatedAmount:  j.Balance,
			ValidatorAddress: j.Delegation.ValidatorAddress,
		})
		total = total.Add(j.Balance)
	}

	// generate a 2d string array for populating csv files.
	var data [][]string
	for _, i := range list {
		var temp []string
		temp = append(
			temp,
			address,
			i.ValidatorAddress,
			i.DelegatedAmount.String(),
			total.Amount.Quo(i.DelegatedAmount.Amount).String()+"%",
		)
		data = append(data, temp)
	}

	err = export.WriteCSV(
		grpctypes.PrefixGRPC+grpctypes.DelegationAnalysis,
		[]string{
			grpctypes.HeaderAddress,
			grpctypes.HeaderValidator,
			grpctypes.HeaderTotalDelegations,
			grpctypes.HeaderPercentageWeight,
		},
		data,
	)
	if err != nil {
		return err
	}

	return nil

}
