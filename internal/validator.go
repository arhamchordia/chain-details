package internal

import (
	"context"
	"github.com/arhamchrodia/validator-status/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"google.golang.org/grpc"
)

// ParseValidators parses all the requested information
// returns an error if any of the steps fail
func ParseValidators(grpcConn *grpc.ClientConn, accountPrefix string) error {
	// create a staking client in order to query validators list
	stakingClient := stakingtypes.NewQueryClient(grpcConn)

	// query the validators list using stakingClient
	stakingResponse, err := stakingClient.Validators(
		context.Background(),
		&stakingtypes.QueryValidatorsRequest{
			Pagination: &query.PageRequest{
				Limit: 500,
			},
		})
	if err != nil {
		return err
	}

	// convert response to internal validators for sorting
	validators := types.ConvertToInternalValidators(stakingResponse.Validators)
	validators.SortStable()

	// list of validator names
	monikerList := validators.GetListOfMoniker()

	// list of percentage weight
	percentageWeight := validators.GetListOfDecPercentage(validators.GetTotalShares())

	// total delegations for each validator
	totalDelegations := validators.GetTotalDelegations()

	// get self delegations for each validator
	selfDelegations, err := GetSelfDelegations(stakingClient, validators, accountPrefix)
	if err != nil {
		return err
	}

	// generate a 2d string array for populating csv files.
	var data [][]string
	for i := range validators {
		var temp []string
		temp = append(
			temp,
			monikerList[i],
			percentageWeight[i].String(),
			selfDelegations[i].String(),
			totalDelegations[i].String(),
		)
		data = append(data, temp)
	}

	err = WriteCSV(
		types.ValidatorsInfoFileName,
		[]string{
			types.HeaderMoniker,
			types.HeaderPercentageWeight,
			types.HeaderSelfDelegation,
			types.HeaderTotalDelegations,
		},
		data,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetSelfDelegations returns a list of coin and error after querying self delegations of the
// validators provided in the input.
func GetSelfDelegations(stakingClient stakingtypes.QueryClient, validators types.Validators, accountPrefix string) ([]sdk.Coin, error) {
	var selfDelegations []sdk.Coin

	// iterate through all the validator address and find self delegations
	for _, val := range validators {
		// find account address of the given validator address
		accAddress, err := types.GetAccAddress(val.OperatorAddress, accountPrefix)
		if err != nil {
			return nil, err
		}

		// make a delegation query with the account address and the validator address
		// if there is an error : it means that delegators address has no self delegations
		delegatorValidator, err := stakingClient.Delegation(
			context.Background(),
			&stakingtypes.QueryDelegationRequest{
				ValidatorAddr: val.OperatorAddress,
				DelegatorAddr: accAddress,
			})
		if err != nil {
			selfDelegations = append(selfDelegations, sdk.Coin{})
		} else {
			selfDelegations = append(selfDelegations, delegatorValidator.DelegationResponse.Balance)
		}
	}

	return selfDelegations, nil
}
