package types

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Delegator struct {
	DelegatorAddress string
	ValidatorAddress string
	Share            sdk.Dec
}

type Delegators []Delegator

func (d Delegators) SortStable() {
	// sort with slice stable
	sort.SliceStable(d, func(i, j int) bool { return d[i].DelegatorAddress < d[j].DelegatorAddress })
}
