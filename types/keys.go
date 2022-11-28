package types

import (
	"errors"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

const (
	ValidatorsInfoFileName             = "validators_info"
	DelegatorDelegationEntriesFileName = "delegator_delegation_entries"
	DelegatorSharesFileName            = "delegator_shares"

	HeaderMoniker          = "Moniker"
	HeaderPercentageWeight = "Percentage Weight"
	HeaderSelfDelegation   = "Self Delegation"
	HeaderTotalDelegations = "Total Delegations"
	HeaderDelegator        = "Delegator"
	HeaderValidator        = "Validator"
	HeaderShares           = "Shares"

	ValidatorsLimit = 50000
	DelegatorsLimit = 1000000000
)

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
