package types

import (
	"strconv"
	"time"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

type Genesis struct {
	AppState AppState `json:"app_state"`
}

type AppState struct {
	Auth GenesisState `json:"auth"`
}

type GenesisState struct {
	Params   Params    `json:"params"`
	Accounts []Account `json:"accounts,omitempty"`
}

type Params struct {
	MaxMemoCharacters      string `json:"max_memo_characters"`
	TxSigLimit             string `json:"tx_sig_limit"`
	TxSizeCostPerByte      string `json:"tx_size_cost_per_byte"`
	SigVerifyCostED25519   string `json:"sig_verify_cost_ed25519"`
	SigVerifyCostSecp256k1 string `json:"sig_verify_cost_secp256k1"`
}

type Account struct {
	Type               string             `json:"@type"`
	BaseVestingAccount BaseVestingAccount `json:"base_vesting_account"`
	StartTime          string             `json:"start_time"`
	VestingPeriods     []VestingPeriods   `json:"vesting_periods"`
}

type BaseVestingAccount struct {
	BaseAccount      BaseAccount    `protobuf:"bytes,1,opt,name=base_account,json=baseAccount,proto3,embedded=base_account" json:"base_account,omitempty"`
	OriginalVesting  sdktypes.Coins `protobuf:"bytes,2,rep,name=original_vesting,json=originalVesting,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"original_vesting" yaml:"original_vesting"`
	DelegatedFree    sdktypes.Coins `protobuf:"bytes,3,rep,name=delegated_free,json=delegatedFree,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"delegated_free" yaml:"delegated_free"`
	DelegatedVesting sdktypes.Coins `protobuf:"bytes,4,rep,name=delegated_vesting,json=delegatedVesting,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"delegated_vesting" yaml:"delegated_vesting"`
	EndTime          string         `protobuf:"varint,5,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty" yaml:"end_time"`
}

type BaseAccount struct {
	Address       string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	PubKey        string `protobuf:"bytes,2,opt,name=pub_key,json=pubKey,proto3" json:"public_key,omitempty" yaml:"public_key"`
	AccountNumber string `protobuf:"varint,3,opt,name=account_number,json=accountNumber,proto3" json:"account_number,omitempty" yaml:"account_number"`
	Sequence      string `protobuf:"varint,4,opt,name=sequence,proto3" json:"sequence,omitempty"`
}

type VestingPeriods struct {
	Length string         `json:"length"`
	Amount sdktypes.Coins `json:"amount"`
}

type DelayedVestingAccount struct {
	Address string
	EndTime time.Time
	Tokens  sdktypes.Coins
}

type ContinuousVestingAccount struct {
	Address              string
	StartTime            time.Time
	EndTime              time.Time
	Tokens               sdktypes.Coins
	TokensFreeEveryBlock sdktypes.Int
	TokensFreeEveryDay   sdktypes.Int
}

type PermanentLockedAccount struct {
	Address string
	Tokens  sdktypes.Coins
}

type PeriodicVestingAccount struct {
	Address   string
	StartTime time.Time
	EndTime   time.Time
	Tokens    sdktypes.Coins
}

func (a Account) GetType() string {
	return a.Type
}

func (a Account) GetAddress() string {
	return a.BaseVestingAccount.BaseAccount.Address
}

func (a Account) GetOriginalVesting() sdktypes.Coins {
	return a.BaseVestingAccount.OriginalVesting
}

func (a Account) GetEndTime() (time.Time, error) {
	seconds, err := strconv.ParseInt(a.BaseVestingAccount.EndTime, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	tm := time.Unix(seconds, 0)
	return tm, err
}

func (a Account) GetStartTime() (time.Time, error) {
	seconds, err := strconv.ParseInt(a.StartTime, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	tm := time.Unix(seconds, 0)
	return tm, err
}

func (a Account) GetEndTimeUNIX() (int, error) {
	endTimeUNIX, err := strconv.Atoi(a.BaseVestingAccount.EndTime)
	if err != nil {
		return 0, err
	}
	return endTimeUNIX, nil
}

func (a Account) GetStartTimeUNIX() (int, error) {
	startTimeUNIX, err := strconv.Atoi(a.StartTime)
	if err != nil {
		return 0, err
	}
	return startTimeUNIX, nil
}
