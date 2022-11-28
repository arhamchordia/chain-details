package internal

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"

	"github.com/arhamchordia/chain-details/types"
)

// ParseDelegators parses all the requested information
// returns an error if any of the steps fail
func ParseDelegators(grpcConn *grpc.ClientConn) error {
	// create a staking client in order to query validators list
	stakingClient := stakingtypes.NewQueryClient(grpcConn)

	// query the validators list using stakingClient
	stakingResponse, err := stakingClient.Validators(
		context.Background(),
		&stakingtypes.QueryValidatorsRequest{
			Pagination: &query.PageRequest{Limit: types.ValidatorsLimit},
		})
	if err != nil {
		return err
	}

	// convert response to internal validators
	validators := types.ConvertToInternalValidators(stakingResponse.Validators)
	validators.SortStable()

	// define delegators array
	var delegatorsSlice types.Delegators

	// iterate through all the validators to find out delegations made to them
	for _, val := range validators {
		validatorDelegations, err := stakingClient.ValidatorDelegations(
			context.Background(),
			&stakingtypes.QueryValidatorDelegationsRequest{
				ValidatorAddr: val.OperatorAddress,
				Pagination:    &query.PageRequest{Limit: types.DelegatorsLimit},
			},
		)
		if err != nil {
			return err
		}

		// iterate through all the delegations made to a particular validator
		for _, delRes := range validatorDelegations.DelegationResponses {
			delegatorsSlice = append(
				delegatorsSlice,
				types.Delegator{
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

	err = WriteCSV(
		types.DelegatorDelegationEntriesFileName,
		[]string{
			types.HeaderDelegator,
			types.HeaderValidator,
			types.HeaderShares,
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

	err = WriteCSV(
		types.DelegatorSharesFileName,
		[]string{
			types.HeaderDelegator,
			types.HeaderShares,
		},
		delegatorShares,
	)
	if err != nil {
		return err
	}

	return nil
}
