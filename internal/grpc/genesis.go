package grpc

import (
	"encoding/json"
	"github.com/arhamchordia/chain-details/internal/export"
	grpctypes "github.com/arhamchordia/chain-details/types/grpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io"
	"net/http"
	"strconv"
)

func QueryGenesisJSON(jsonURL, denom string) error {
	res, err := http.Get(jsonURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var response grpctypes.Genesis
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	err = parseVestingAccounts(response.AppState.Auth.Accounts, denom)
	if err != nil {
		return err
	}

	return nil
}

// parseVestingAccounts parses all the vesting accounts in genesis file and
// returns and error if anything fails
func parseVestingAccounts(vestingAccounts []grpctypes.Account, denom string) error {
	// initialise all types of vesting account arrays
	var (
		delayedVestingAccounts    []grpctypes.DelayedVestingAccount
		continuousVestingAccounts []grpctypes.ContinuousVestingAccount
		permanentLockedAccounts   []grpctypes.PermanentLockedAccount
		periodicVestingAccounts   []grpctypes.PeriodicVestingAccount
	)

	// iterate over all types of vesting accounts and performs
	// corresponding operations
	for _, account := range vestingAccounts {
		switch account.GetType() {
		// only end time and tokens are stored in delayed vesting account
		// as the tokens are locked till the end time of the vesting
		case grpctypes.IdentifierDelayedVestingAccount:
			// get end time
			endTime, err := account.GetEndTime()
			if err != nil {
				return err
			}

			delayedVestingAccounts = append(
				delayedVestingAccounts,
				grpctypes.DelayedVestingAccount{
					Address: account.GetAddress(),
					Tokens:  account.GetOriginalVesting(),
					EndTime: endTime,
				})

		// in continuous vesting account, the tokens are continuously being freed up every block
		// tokens every block = tokens / number of block between the period
		case grpctypes.IdentifierContinuousVestingAccount:
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
			numberOfBlockInBetween := int64((endTimeUNIX - startTimeUNIX) / grpctypes.AverageBlockTime)
			numberOfDayaInBetween := int64((endTimeUNIX - startTimeUNIX) / grpctypes.SecondsInADay)

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
				grpctypes.ContinuousVestingAccount{
					Address:              account.GetAddress(),
					StartTime:            startTime,
					EndTime:              endTime,
					Tokens:               account.GetOriginalVesting(),
					TokensFreeEveryBlock: tokensFreeEveryBlock,
					TokensFreeEveryDay:   tokensFreeEveryDay,
				})

		// in permanent locked account, the tokens are locked forever and can
		// can only be used in case of delegating and participating in governance
		case grpctypes.IdentifierPermanentLockedAccount:
			permanentLockedAccounts = append(
				permanentLockedAccounts,
				grpctypes.PermanentLockedAccount{
					Address: account.GetAddress(),
					Tokens:  account.GetOriginalVesting(),
				})

		// in periodic vesting accounts, the tokens are freed up during the mentioned
		// periods in the vesting_period array.
		case grpctypes.IdentifierPeriodicVestingAccount:
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
					grpctypes.PeriodicVestingAccount{
						Address:   account.GetAddress(),
						Tokens:    account.VestingPeriods[i].Amount,
						StartTime: grpctypes.GetTimeFromUNIXTimeStamp(startTimeUNIX),
						EndTime:   grpctypes.GetTimeFromUNIXTimeStamp(endTimeOfThisPeriod),
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

	err := export.WriteCSV(
		grpctypes.PrefixGRPC+grpctypes.GenesisAccountAnalysisFileName,
		[]string{
			grpctypes.HeaderAddress,
			grpctypes.HeaderOriginalVesting,
			grpctypes.HeaderVestingEndTime,
			grpctypes.HeaderVestingStartTime,
			grpctypes.HeaderTokensFreeEveryBlock,
			grpctypes.HeaderTokensFreeEveryDay,
		},
		data,
	)
	if err != nil {
		return err
	}

	return nil
}
