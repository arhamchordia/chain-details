package vaults

import (
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	"github.com/arhamchordia/chain-details/internal/export"
	bigquerytypes "github.com/arhamchordia/chain-details/types/bigquery"
	"log"
	"strconv"
	"strings"
)

func CLQueryDeposit(addressQuery string, outputFormat string) error {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependCLQueryVaultsDeposit

	headers, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryCLVaultsDeposit, addressQuery, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = export.ExportFile(outputFormat, filename, headers, rows)
	if err != nil {
		return err
	}

	return nil
}

func CLQueryWithdraw(addressQuery string, outputFormat string) error {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependCLQueryVaultsWithdraw

	headers, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryCLVaultsWithdraw, addressQuery, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = export.ExportFile(outputFormat, filename, headers, rows)
	if err != nil {
		return err
	}

	return nil
}

func CLQueryClaim(addressQuery string, outputFormat string) error {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependCLQueryVaultsClaim

	headers, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryCLVaultsClaim, addressQuery, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = export.ExportFile(outputFormat, filename, headers, rows)
	if err != nil {
		return err
	}

	return nil
}

// TODO Mock conversion rates for denoms to USD
var MOCK_CONVERSION_RATES = map[string]float64{
	"uosmo": 0.29,
	"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2": 7.35, // ATOM
	"ibc/57AA1A70A4BC9769C525EBF6386F7A21536E04A79D62E1981EFCEF9428EBB205": 0.60, // KAVA
	"ibc/4ABBEF4C8926DDDB320AE5188CFD63267ABBCEFC0583E4AE05D6E5AA2401DDAB": 1.00, // USDT/USDC
	"ibc/498A0751C798A0D9A389AA3691123DADA57DAA4FE165D5C75894505B876BA6E4": 1.00, // USDT/USDC
}

// Map to store the conversion exponent for each denom
var DENOM_EXPONENTS = map[string]float64{
	"uosmo": 1000000.0, // 6 decimal places for uosmo
	// Add other denoms as needed. For example:
	// "anotherDenom": 10000.0,  // 4 decimal places for anotherDenom
}

const DEFAULT_EXPONENT = 1000000.0 // Default to 6 decimal places

func aggregateByDenomAndConvert(amountList []string) map[string]float64 {
	// Initialize a map to aggregate amounts
	aggregated := make(map[string]float64)
	fmt.Println("amountList", amountList)

	for _, amount := range amountList {
		// Remove the square brackets and split by space to get individual amounts
		amounts := strings.Split(strings.Trim(amount, "[]"), " ")

		for _, amt := range amounts {
			for denom, _ := range MOCK_CONVERSION_RATES {
				if strings.Contains(amt, denom) {
					valueStr := strings.Split(amt, denom)[0]
					value, err := strconv.ParseFloat(valueStr, 64)
					if err != nil {
						// Handle error if needed, or continue
						continue
					}
					aggregated[denom] = aggregated[denom] + value
				}
			}
		}
	}
	fmt.Println("aggregated", aggregated)

	// Convert aggregated amounts to mock USD
	converted := make(map[string]float64)
	for denom, value := range aggregated {
		// Get the conversion exponent for the denom, or default to DEFAULT_EXPONENT
		exponent, exists := DENOM_EXPONENTS[denom]
		if !exists {
			exponent = DEFAULT_EXPONENT
		}

		// Convert micro units (or other units) to main units
		mainUnitValue := value / exponent

		conversionRate := MOCK_CONVERSION_RATES[denom]
		converted[denom] = mainUnitValue * conversionRate
	}

	return converted
}

func sumUSDConversion(amountList []string) float64 {
	// This function will use the logic from aggregateByDenomAndConvert to
	// compute the total USD value for a list of amounts with denoms.
	totalUSD := 0.0
	convertedAmounts := aggregateByDenomAndConvert(amountList)
	for _, usdValue := range convertedAmounts {
		totalUSD += usdValue
	}
	return totalUSD
}

func CLQueryDistributeRewards(addressQuery string, daysInterval int, outputFormat string) error {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependCLQueryVaultsDistributeRewards

	headers, rows, err := internal.ExecuteQueryAndFetchRows(fmt.Sprintf(bigquerytypes.QueryCLVaultsDistributeRewards, addressQuery, daysInterval), "", false)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Determine column indices for "amount_incentives" and "amount_spread_rewards"
	incentivesColIndex := -1
	spreadRewardsColIndex := -1
	for i, header := range headers {
		if header == "amount_incentives" {
			incentivesColIndex = i
		} else if header == "amount_spread_rewards" {
			spreadRewardsColIndex = i
		}
	}

	if incentivesColIndex == -1 || spreadRewardsColIndex == -1 {
		log.Fatalf("Expected columns not found in the data")
	}

	// Add columns for USD conversions
	headers = append(headers, "amount_incentives_usd", "amount_spread_rewards_usd")

	// Process rows
	for rowIndex, row := range rows {
		// Access the columns using the determined indices
		incentivesList := strings.Split(row[incentivesColIndex], ",")
		spreadRewardsList := strings.Split(row[spreadRewardsColIndex], ",")

		// Convert amount_incentives and amount_spread_rewards to USD
		incentivesUSD := sumUSDConversion(incentivesList)
		spreadRewardsUSD := sumUSDConversion(spreadRewardsList)

		// Append the USD values to the row
		rows[rowIndex] = append(row, fmt.Sprintf("%.2f", incentivesUSD), fmt.Sprintf("%.2f", spreadRewardsUSD))
	}

	err = export.ExportFile(outputFormat, filename, headers, rows)
	if err != nil {
		return err
	}

	return nil
}

func CLQueryAPR(addressQuery string, daysInterval int, outputFormat string) error {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependCLQueryVaultsAPR

	headers, rows, err := internal.ExecuteQueryAndFetchRows(fmt.Sprintf(bigquerytypes.QueryCLVaultsDistributeRewards, addressQuery, daysInterval), "", false)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Determine column indices for "amount_incentives" and "amount_spread_rewards"
	incentivesColIndex := -1
	spreadRewardsColIndex := -1
	for i, header := range headers {
		if header == "amount_incentives" {
			incentivesColIndex = i
		} else if header == "amount_spread_rewards" {
			spreadRewardsColIndex = i
		}
	}

	// Initialize aggregated amounts for incentives and spread rewards
	totalIncentivesUSD := 0.0
	totalSpreadRewardsUSD := 0.0

	// Process rows to aggregate total incentives and spread rewards in USD
	for _, row := range rows {
		// Access the columns using the determined indices
		incentivesList := strings.Split(row[incentivesColIndex], ",")
		spreadRewardsList := strings.Split(row[spreadRewardsColIndex], ",")

		totalIncentivesUSD += sumUSDConversion(incentivesList)
		totalSpreadRewardsUSD += sumUSDConversion(spreadRewardsList)
	}

	// Compute averages
	avgIncentivesUSD := totalIncentivesUSD / float64(len(rows))
	avgSpreadRewardsUSD := totalSpreadRewardsUSD / float64(len(rows))

	// Prepare headers and a single row with the aggregated information
	headers = []string{
		"average_amount_incentives_usd",
		"average_amount_spread_rewards_usd",
		"total_amount_incentives_usd",
		"total_amount_spread_rewards_usd",
		"number_of_transactions",
	}
	aggregatedRow := []string{
		fmt.Sprintf("%.2f", avgIncentivesUSD),
		fmt.Sprintf("%.2f", avgSpreadRewardsUSD),
		fmt.Sprintf("%.2f", totalIncentivesUSD),
		fmt.Sprintf("%.2f", totalSpreadRewardsUSD),
		strconv.Itoa(len(rows)),
	}

	err = export.ExportFile(outputFormat, filename, headers, [][]string{aggregatedRow})
	if err != nil {
		return err
	}

	return nil
}

// CLQueryReport returns a file with the last 24h statistics of a given vaultAddress
func CLQueryReport(addressQuery string, outputFormat string) error {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}

	//filename := fmt.Sprintf("%s_%s", bigquerytypes.PrefixBigQuery+bigquerytypes.PrependCLQueryReport, addressQuery)

	// Making queries as in the above methods
	deposit_headers, deposit_rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryCLVaultsDeposit, addressQuery, true)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println(">>> DEPOSIT", deposit_headers, deposit_rows)
	withdraw_headers, withdraw_rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryCLVaultsWithdraw, addressQuery, true)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println(">>> DEPOSIT", withdraw_headers, withdraw_rows)

	claim_headers, claim_rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryCLVaultsClaim, addressQuery, true)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println(">>> DEPOSIT", claim_headers, claim_rows)

	// TODO distribute_rewards query

	/*err = internal.ExportFile(outputFormat, filename, headers, rows)
	if err != nil {
		return err
	}*/

	return nil
}
