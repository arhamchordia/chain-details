package grpc

import (
	"context"
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"strconv"
	"time"

	"github.com/arhamchordia/chain-details/internal"
	grpctypes "github.com/arhamchordia/chain-details/types/grpc"
	cryptotypes1 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/query"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func QueryValidatorsData(grpcUrl, accountPrefix string) error {
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

		if grpcConn.GetState().String() == "READY" {
			err = ParseValidators(grpcConn, accountPrefix)
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
				Limit: grpctypes.ValidatorsLimit,
			},
		})
	if err != nil {
		return err
	}

	// convert response to internal validators for sorting
	validators := grpctypes.ConvertToInternalValidators(stakingResponse.Validators)
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

	err = internal.WriteCSV(
		grpctypes.PrefixGRPC+grpctypes.ValidatorsInfoFileName,
		[]string{
			grpctypes.HeaderMoniker,
			grpctypes.HeaderPercentageWeight,
			grpctypes.HeaderSelfDelegation,
			grpctypes.HeaderTotalDelegations,
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
func GetSelfDelegations(stakingClient stakingtypes.QueryClient, validators grpctypes.Validators, accountPrefix string) ([]sdk.Coin, error) {
	var selfDelegations []sdk.Coin

	// iterate through all the validator address and find self delegations
	for _, val := range validators {
		// find account address of the given validator address
		accAddress, err := grpctypes.GetAccAddressFromValAdderss(val.OperatorAddress, accountPrefix)
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

func ParseGenesisPostGenesisValidatorsData(grpcUrl string) error {
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

		if grpcConn.GetState().String() == "READY" {
			err = ParseGenesisAndPostGenesisValidators(grpcConn)
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

func ParseGenesisAndPostGenesisValidators(grpcConn *grpc.ClientConn) error {
	// create a staking client in order to query validators list
	stakingClient := stakingtypes.NewQueryClient(grpcConn)

	// query the validators list using stakingClient
	stakingResponse, err := stakingClient.Validators(
		context.Background(),
		&stakingtypes.QueryValidatorsRequest{
			Pagination: &query.PageRequest{
				Limit: grpctypes.ValidatorsLimit,
			},
		})
	if err != nil {
		return err
	}

	slashingClient := slashingtypes.NewQueryClient(grpcConn)

	var genesisValidators []stakingtypes.Validator
	var postGenesisValidators []stakingtypes.Validator
	for _, val := range stakingResponse.Validators {
		validator, err := stakingClient.Validator(context.Background(), &stakingtypes.QueryValidatorRequest{ValidatorAddr: val.OperatorAddress})
		if err != nil {
			return err
		}

		var pubKey cryptotypes1.PubKey
		err = pubKey.Unmarshal(validator.Validator.ConsensusPubkey.Value)
		if err != nil {
			return err
		}

		accAddress := sdk.AccAddress(pubKey.Address())
		bech32Addr, err := bech32.ConvertAndEncode("quasarvalcons", accAddress)
		if err != nil {
			return err
		}

		info, err := slashingClient.SigningInfo(context.Background(), &slashingtypes.QuerySigningInfoRequest{ConsAddress: bech32Addr})
		if err != nil {
			return err
		}

		if info.ValSigningInfo.StartHeight != 0 {
			postGenesisValidators = append(postGenesisValidators, validator.Validator)
		} else {
			genesisValidators = append(genesisValidators, validator.Validator)
		}
	}

	// generate a 2d string array for populating csv files.
	var data [][]string
	for _, i := range genesisValidators {
		var temp []string
		temp = append(
			temp,
			"genesis",
			i.OperatorAddress,
			i.ConsensusPubkey.String(),
			strconv.FormatBool(i.Jailed),
			i.Status.String(),
			i.Tokens.String(),
			i.DelegatorShares.String(),
			i.Description.String(),
			strconv.FormatInt(i.UnbondingHeight, 10),
			i.UnbondingTime.String(),
			i.Commission.String(),
			i.MinSelfDelegation.String(),
		)
		data = append(data, temp)
	}

	for _, i := range postGenesisValidators {
		var temp []string
		temp = append(
			temp,
			"post-genesis",
			i.OperatorAddress,
			i.ConsensusPubkey.String(),
			strconv.FormatBool(i.Jailed),
			i.Status.String(),
			i.Tokens.String(),
			i.DelegatorShares.String(),
			i.Description.String(),
			strconv.FormatInt(i.UnbondingHeight, 10),
			i.UnbondingTime.String(),
			i.Commission.String(),
			i.MinSelfDelegation.String(),
		)
		data = append(data, temp)
	}

	err = internal.WriteCSV(
		grpctypes.PrefixGRPC+grpctypes.GenesisPostGenesisValidators,
		[]string{
			grpctypes.HeaderGenesisType,
			grpctypes.HeaderOperatorAddress,
			grpctypes.HeaderConsensusPubkey,
			grpctypes.HeaderJailed,
			grpctypes.HeaderStatus,
			grpctypes.HeaderTokens,
			grpctypes.HeaderDelegatorShares,
			grpctypes.HeaderDescription,
			grpctypes.HeaderUnbondingHeight,
			grpctypes.HeaderUnbondingTime,
			grpctypes.HeaderCommission,
			grpctypes.HeaderMinSelfDelegation,
		},
		data,
	)
	if err != nil {
		return err
	}

	return nil
}
