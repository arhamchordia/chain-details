package internal

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arhamchordia/chain-details/types"
)

// ParseVestingAccounts parses all the vesting accounts in genesis file and
// returns and error if anything fails
func ParseVestingAccounts(vestingAccounts []types.Account, denom string) error {
	// initialise all types of vesting account arrays
	var (
		delayedVestingAccounts    []types.DelayedVestingAccount
		continuousVestingAccounts []types.ContinuousVestingAccount
		permanentLockedAccounts   []types.PermanentLockedAccount
		periodicVestingAccounts   []types.PeriodicVestingAccount
	)

	// iterate over all types of vesting accounts and performs
	// corresponding operations
	for _, account := range vestingAccounts {
		switch account.GetType() {
		// only end time and tokens are stored in delayed vesting account
		// as the tokens are locked till the end time of the vesting
		case types.IdentifierDelayedVestingAccount:
			// get end time
			endTime, err := account.GetEndTime()
			if err != nil {
				return err
			}

			delayedVestingAccounts = append(
				delayedVestingAccounts,
				types.DelayedVestingAccount{
					Address: account.GetAddress(),
					Tokens:  account.GetOriginalVesting(),
					EndTime: endTime,
				})

		// in continuous vesting account, the tokens are continuously being freed up every block
		// tokens every block = tokens / number of block between the period
		case types.IdentifierContinuousVestingAccount:
			// convert start time and end time in int
			startTimeUNIX, err := account.GetStartTimeUNIX()
			if err != nil {
				return err
			}
			endTimeUNIX, err := account.GetEndTimeUNIX()
			if err != nil {
				return err
			}

			// calculate number of blocks generated during that period and number of days in between the period
			numberOfBlockInBetween := int64((endTimeUNIX - startTimeUNIX) / types.AverageBlockTime)
			numberOfDayaInBetween := int64((endTimeUNIX - startTimeUNIX) / types.SecondsInADay)

			// calculate tokens freed up every block and everyday
			tokensFreeEveryBlock := account.GetOriginalVesting().AmountOf(denom).Quo(sdk.NewInt(numberOfBlockInBetween))
			tokensFreeEveryDay := account.GetOriginalVesting().AmountOf(denom).Quo(sdk.NewInt(numberOfDayaInBetween))

			// get start time
			startTime, err := account.GetStartTime()
			if err != nil {
				return err
			}

			// get end time
			endTime, err := account.GetEndTime()
			if err != nil {
				return err
			}

			continuousVestingAccounts = append(
				continuousVestingAccounts,
				types.ContinuousVestingAccount{
					Address:              account.GetAddress(),
					StartTime:            startTime,
					EndTime:              endTime,
					Tokens:               account.GetOriginalVesting(),
					TokensFreeEveryBlock: tokensFreeEveryBlock,
					TokensFreeEveryDay:   tokensFreeEveryDay,
				})

		// in permanent locked account, the tokens are locked forever and can
		// can only be used in case of delegating and participating in governance
		case types.IdentifierPermanentLockedAccount:
			permanentLockedAccounts = append(
				permanentLockedAccounts,
				types.PermanentLockedAccount{
					Address: account.GetAddress(),
					Tokens:  account.GetOriginalVesting(),
				})

		// in periodic vesting accounts, the tokens are freed up during the mentioned
		// periods in the vesting_period array.
		case types.IdentifierPeriodicVestingAccount:
			// calculate start time in UNIX
			startTimeUNIX, err := account.GetStartTimeUNIX()
			if err != nil {
				return err
			}

			// iterate over all the vesting period of the given account
			for i := range account.VestingPeriods {
				// convert period string into int
				period, err := strconv.Atoi(account.VestingPeriods[i].Length)
				if err != nil {
					return err
				}

				// calculate start time of that specific period by adding period*i to startTimeUNIX
				startTimeOfThisPeriod := startTimeUNIX + (period * (i))

				// calculate end of that specific period by adding period to the startTimeOfThisPeriod
				endTimeOfThisPeriod := startTimeOfThisPeriod + period

				periodicVestingAccounts = append(
					periodicVestingAccounts,
					types.PeriodicVestingAccount{
						Address:   account.GetAddress(),
						Tokens:    account.VestingPeriods[i].Amount,
						StartTime: types.GetTimeFromUNIXTimeStamp(startTimeUNIX),
						EndTime:   types.GetTimeFromUNIXTimeStamp(endTimeOfThisPeriod),
					},
				)
			}
		}
	}

	var data [][]string
	for _, dva := range delayedVestingAccounts {
		data = append(
			data,
			[]string{
				dva.Address,
				dva.Tokens.String(),
				dva.EndTime.String(),
			})
	}
	for _, cva := range continuousVestingAccounts {
		data = append(
			data,
			[]string{
				cva.Address,
				cva.EndTime.String(),
				cva.Tokens.String(),
				cva.StartTime.String(),
				cva.TokensFreeEveryBlock.String(),
				cva.TokensFreeEveryDay.String(),
			})
	}
	for _, pla := range permanentLockedAccounts {
		data = append(
			data,
			[]string{
				pla.Address,
				pla.Tokens.String(),
			})
	}

	for _, pva := range periodicVestingAccounts {
		data = append(
			data,
			[]string{
				pva.Address,
				pva.Tokens.String(),
				pva.EndTime.String(),
				pva.StartTime.String(),
			},
		)
	}

	err := WriteCSV(
		types.GenesisAccountAnalysisFileName,
		[]string{
			types.HeaderAddress,
			types.HeaderOriginalVesting,
			types.HeaderVestingEndTime,
			types.HeaderVestingStartTime,
			types.HeaderTokensFreeEveryBlock,
			types.HeaderTokensFreeEveryDay,
		},
		data,
	)
	if err != nil {
		return err
	}

	return nil
}
