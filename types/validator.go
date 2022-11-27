package types

import (
	"encoding/hex"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"sort"
	"strings"
)

type Validators []stakingtypes.Validator

func (v Validators) SortStable() {
	sort.SliceStable(v, func(i, j int) bool { return v[i].DelegatorShares.LT(v[j].DelegatorShares) })
}

func (v Validators) GetListOfMoniker() []string {
	var monikerList []string
	for _, val := range v {
		monikerList = append(monikerList, val.Description.Moniker)
	}
	return monikerList
}

func (v Validators) GetTotalShares() sdk.Dec {
	totalShares := sdk.ZeroDec()
	for _, val := range v {
		totalShares = totalShares.Add(val.DelegatorShares)
	}
	return totalShares
}

func (v Validators) GetListOfDecPercentage(totalShares sdk.Dec) []sdk.Dec {
	var decPercentage []sdk.Dec
	for _, val := range v {
		valWeight := val.DelegatorShares.Mul(sdk.NewDec(100)).Quo(totalShares)
		decPercentage = append(decPercentage, valWeight)
	}
	return decPercentage
}

func (v Validators) GetTotalDelegations() []sdk.Int {
	var selfDelegations []sdk.Int
	for _, val := range v {
		selfDelegations = append(selfDelegations, val.Tokens)
	}
	return selfDelegations
}

func (v Validators) GetAccountAddressesList() ([]string, error) {
	var addressesList []string

	for _, val := range v {
		valAddress, err := sdk.ValAddressFromBech32(val.OperatorAddress)
		if err != nil {
			return []string{}, err
		}
		hexValAddress := hex.EncodeToString(valAddress)

		accAddress1, err := sdk.AccAddressFromHex(hexValAddress)
		if err != nil {
			return []string{}, err
		}

		addressesList = append(addressesList, accAddress1.String())
	}
	return addressesList, nil
}

func ConvertToInternalValidators(validatorList stakingtypes.Validators) Validators {
	var validators Validators
	for _, vals := range validatorList {
		validators = append(validators, vals)
	}
	return validators
}

func GetAccAddress(address, accountPrefix string) (string, error) {
	valAddressPrefix := accountPrefix + "valoper"
	valAddress, err := ValAddressFromBech32(address, valAddressPrefix)
	if err != nil {
		return "", err
	}

	hexValAddress := hex.EncodeToString(valAddress)

	accAddress, err := sdk.AccAddressFromHex(hexValAddress)
	if err != nil {
		return "", err
	}

	prefixAddress, err := Bech32ifyAddressBytes(accountPrefix, accAddress)
	if err != nil {
		return "", err
	}

	return prefixAddress, nil
}

// ValAddressFromBech32 creates a ValAddress from a Bech32 string.
func ValAddressFromBech32(address, prefix string) (valAddr sdk.ValAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdk.ValAddress{}, errors.New("empty address string is not allowed")
	}

	bz, err := sdk.GetFromBech32(address, prefix)
	if err != nil {
		return nil, err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// Bech32ifyAddressBytes returns a bech32 representation of address bytes.
// Returns an empty sting if the byte slice is 0-length. Returns an error if the bech32 conversion
// fails or the prefix is empty.
func Bech32ifyAddressBytes(prefix string, address sdk.AccAddress) (string, error) {
	if address.Empty() {
		return "", nil
	}
	if len(address.Bytes()) == 0 {
		return "", nil
	}
	if len(prefix) == 0 {
		return "", errors.New("prefix cannot be empty")
	}
	return bech32.ConvertAndEncode(prefix, address.Bytes())
}
