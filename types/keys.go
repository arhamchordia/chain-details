package types

import (
	"errors"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

const (
	ValidatorsInfoFileName             = "validators_info"
	DelegatorDelegationEntriesFileName = "delegator_delegation_entries"
	DelegatorSharesFileName            = "delegator_shares"
	GenesisAccountAnalysisFileName     = "genesis_accounts"

	HeaderMoniker              = "Moniker"
	HeaderPercentageWeight     = "Percentage Weight"
	HeaderSelfDelegation       = "Self Delegation"
	HeaderTotalDelegations     = "Total Delegations"
	HeaderDelegator            = "Delegator"
	HeaderValidator            = "Validator"
	HeaderShares               = "Shares"
	HeaderAddress              = "Address"
	HeaderVestingEndTime       = "Vesting End Time"
	HeaderOriginalVesting      = "Original Vesting"
	HeaderVestingStartTime     = "Vesting Start Time"
	HeaderTokensFreeEveryBlock = "Tokens Free Every Block"
	HeaderTokensFreeEveryDay   = "Tokens Free Every Day"

	Message           = "message"
	Wasm              = "wasm"
	BondID            = "bond_id"
	CoinSpent         = "coin_spent"
	CoinReceived      = "coin_received"
	ContractAddress   = "_contract_address"
	LockID            = "lock_id"
	LockedTokens      = "locked_tokens"
	Action            = "action"
	CallbackInfo      = "callback-info"
	ReplyMsgID        = "reply-msg-id"
	ReplyResult       = "reply-result"
	User              = "user"
	VaultTokenBalance = "vault_token_balance"
	Websocket         = "/websocket"

	VaultAddress      = "quasar18a2u6az6dzw528rptepfg6n49ak6hdzkf8ewf0n5r0nwju7gtdgqamr7qu"
	PrimitiveAddress1 = "quasar1kj8q8g2pmhnagmfepp9jh9g2mda7gzd0m5zdq0s08ulvac8ck4dq9ykfps"
	PrimitiveAddress2 = "quasar1ma0g752dl0yujasnfs9yrk6uew7d0a2zrgvg62cfnlfftu2y0egqx8e7fv"
	PrimitiveAddress3 = "quasar1ery8l6jquynn9a4cz2pff6khg8c68f7urt33l5n9dng2cwzz4c4qxhm6a2"

	IdentifierDelayedVestingAccount    = "/cosmos.vesting.v1beta1.DelayedVestingAccount"
	IdentifierContinuousVestingAccount = "/cosmos.vesting.v1beta1.ContinuousVestingAccount"
	IdentifierPermanentLockedAccount   = "/cosmos.vesting.v1beta1.PermanentLockedAccount"
	IdentifierPeriodicVestingAccount   = "/cosmos.vesting.v1beta1.PeriodicVestingAccount"
	IdentifierMsgExecuteContract       = "/cosmwasm.wasm.v1.MsgExecuteContract"
	IdentifierMsgUpdateClient          = "/ibc.core.client.v1.MsgUpdateClient"
	IdentifierMsgAcknowledgement       = "/ibc.core.channel.v1.MsgAcknowledgement"

	ValidatorsLimit  = 50000
	DelegatorsLimit  = 1000000000
	AverageBlockTime = 5
	SecondsInADay    = 86400
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

func GetTimeFromUNIXTimeStamp(unix int) time.Time {
	tm := time.Unix(int64(unix), 0)
	return tm
}
