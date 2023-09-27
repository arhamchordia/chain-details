package grpc

import (
	"context"
	"crypto/tls"
	"github.com/arhamchordia/chain-details/internal"
	"google.golang.org/grpc/credentials"
	"time"

	grpctypes "github.com/arhamchordia/chain-details/types/grpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
)

func QueryDelegatorsData(grpcUrl string) error {
	// initialise config for grpc connection
	config := &tls.Config{
		InsecureSkipVerify: false,
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

		// trigger action on the basis of state of the connection
		if grpcConn.GetState().String() == "READY" {
			err = ParseDelegators(grpcConn)
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

// ParseDelegators parses all the requested information
// returns an error if any of the steps fail
func ParseDelegators(grpcConn *grpc.ClientConn) error {
	// create a staking client in order to query validators list
	stakingClient := stakingtypes.NewQueryClient(grpcConn)

	// query the validators list using stakingClient
	stakingResponse, err := stakingClient.Validators(
		context.Background(),
		&stakingtypes.QueryValidatorsRequest{
			Pagination: &query.PageRequest{Limit: grpctypes.ValidatorsLimit},
		})
	if err != nil {
		return err
	}

	// convert response to internal validators
	validators := grpctypes.ConvertToInternalValidators(stakingResponse.Validators)
	validators.SortStable()

	// define delegators array
	var delegatorsSlice grpctypes.Delegators

	// iterate through all the validators to find out delegations made to them
	for _, val := range validators {
		validatorDelegations, err := stakingClient.ValidatorDelegations(
			context.Background(),
			&stakingtypes.QueryValidatorDelegationsRequest{
				ValidatorAddr: val.OperatorAddress,
				Pagination:    &query.PageRequest{Limit: grpctypes.DelegatorsLimit},
			},
		)
		if err != nil {
			return err
		}

		// iterate through all the delegations made to a particular validator
		for _, delRes := range validatorDelegations.DelegationResponses {
			delegatorsSlice = append(
				delegatorsSlice,
				grpctypes.Delegator{
					DelegatorAddress: delRes.Delegation.DelegatorAddress,
					ValidatorAddress: delRes.Delegation.ValidatorAddress,
					Share:            delRes.Delegation.Shares,
				})
		}
	}

	// stably sort the delegator delegations list
	delegatorsSlice.SortStable()

	// iterate delegator slice for creating data for csv file and populate delegators map
	delegatorsMap := make(map[string]sdk.Dec)
	var delegatorDelegationEntries [][]string
	for _, delegator := range delegatorsSlice {
		// append entries for delegator delegations
		delegatorDelegationEntries = append(
			delegatorDelegationEntries,
			[]string{delegator.DelegatorAddress,
				delegator.ValidatorAddress,
				delegator.Share.String(),
			})

		// insert delegators map with delegator  overall share
		value, ok := delegatorsMap[delegator.DelegatorAddress]
		if ok {
			delegatorsMap[delegator.DelegatorAddress] = value.Add(delegator.Share)
		} else {
			delegatorsMap[delegator.DelegatorAddress] = delegator.Share
		}
	}

	err = internal.WriteCSV(
		grpctypes.PrefixGRPC+grpctypes.DelegatorDelegationEntriesFileName,
		[]string{
			grpctypes.HeaderDelegator,
			grpctypes.HeaderValidator,
			grpctypes.HeaderShares,
		},
		delegatorDelegationEntries,
	)
	if err != nil {
		return err
	}

	var delegatorShares [][]string
	for key, value := range delegatorsMap {
		delegatorShares = append(delegatorShares, []string{key, value.String()})
	}

	err = internal.WriteCSV(
		grpctypes.PrefixGRPC+grpctypes.DelegatorSharesFileName,
		[]string{
			grpctypes.HeaderDelegator,
			grpctypes.HeaderShares,
		},
		delegatorShares,
	)
	if err != nil {
		return err
	}

	return nil
}
