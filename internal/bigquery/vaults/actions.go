package vaults

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	bigquerytypes "github.com/arhamchordia/chain-details/types/bigquery"
	"log"
	"regexp"
	"strings"
)

// QueryBond returns a file with the bond events in all the blocks
func QueryBond(addressQuery string, confirmedQuery bool, pendingQuery bool, outputFormat string) error {
	if confirmedQuery && pendingQuery {
		return fmt.Errorf("--confirmed and --pending flags cannot be used together")
	}

	addressFilterString := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependQueryVaultsBond
	if len(addressQuery) > 0 {
		addressFilterString = fmt.Sprintf(bigquerytypes.QueryVaultsBondAddressFilter, addressQuery)
		filename = fmt.Sprintf("%s_%s", filename, addressQuery)
	}
	if confirmedQuery {
		filename = fmt.Sprintf("%s_%s", filename, "confirmed")
	}
	if pendingQuery {
		filename = fmt.Sprintf("%s_%s", filename, "pending")
	}

	headers, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryVaultsBond, addressFilterString, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if confirmedQuery || pendingQuery {
		_, bondResponses, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryVaultsBondResponseFilter, "", false)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// creating map storing the bond_id and share amounts from the second query
		bondResponsesIDs := make(map[string]int)
		for _, row := range bondResponses {
			bondID := row[0]
			shareAmounts := strings.Split(row[1], ", ") // assuming share_amounts is at index 1, consider changing this if the query changes
			bondResponsesIDs[bondID] = len(shareAmounts)
		}

		_, bondShareAmountsTxIds, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryVaultsBondShareAmountsTxIds, "", false)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Convert txIDs into a slice of strings for the query
		var txIDs []string
		for _, row := range bondShareAmountsTxIds {
			txIDs = append(txIDs, fmt.Sprintf("'%s'", row[0]))
		}
		txIDsStr := strings.Join(txIDs, ", ")

		// Running the second query SELECT message FROM numia-data.quasar.quasar_tx_messages WHERE tx_id IN (%s)
		_, messageRows, err := internal.ExecuteQueryAndFetchRows(fmt.Sprintf("SELECT message FROM numia-data.quasar.quasar_tx_messages WHERE tx_id IN (%s)", txIDsStr), "", false)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Create a map to store bond_ids that are included in the second query
		includedInSecondQuery := make(map[string]bool)

		// Create an array of bond_ids from the JSON struct retrieved, and decoding the "msg" from base64 if it is encoded, otherwise work with JSON already
		for _, row := range messageRows {
			var jsonData map[string]interface{}
			var jsonMsgData map[string]interface{}
			var data []byte

			message := row[0]

			if err := json.Unmarshal([]byte(message), &jsonData); err != nil {
				continue
			}

			msg, ok := jsonData["msg"]
			if !ok {
				continue
			}

			switch v := msg.(type) {
			case string:
				decodedMessage, err := base64.StdEncoding.DecodeString(v)
				if err == nil {
					data = decodedMessage
				} else {
					data = []byte(v)
				}
			case map[string]interface{}:
				jsonMsgData, err := json.Marshal(v)
				if err != nil {
					continue
				}
				data = jsonMsgData
			default:
				continue
			}

			if err := json.Unmarshal(data, &jsonMsgData); err != nil {
				continue
			}

			callbacks, ok := jsonMsgData["callbacks"].([]interface{})
			if ok {
				for _, callback := range callbacks {
					callbackMap, ok := callback.(map[string]interface{})
					if ok {
						bondID, ok := callbackMap["bond_id"].(string)
						if ok {
							includedInSecondQuery[bondID] = true
						}
					}
				}
			}
		}

		filteredRows := [][]string{}
		bondIDRegex := regexp.MustCompile(`bond_id (\d+)`)
		for _, row := range rows {
			column3 := row[2]
			match := bondIDRegex.FindStringSubmatch(column3)
			if len(match) > 1 && bondIDRegex.MatchString(column3) {
				bondID := match[1]
				shareAmounts, exists := bondResponsesIDs[bondID]
				_, included := includedInSecondQuery[bondID]
				if confirmedQuery {
					if (exists && shareAmounts >= 3) || included {
						filteredRows = append(filteredRows, row)
					}
				} else if pendingQuery {
					if (!exists && !included) || (exists && shareAmounts < 3 && !included) {
						filteredRows = append(filteredRows, row)
					}
				}
			}
		}
		rows = filteredRows
	}

	if outputFormat == "csv" {
		err = internal.WriteCSV(filename, headers, rows)
	} else {
		err = internal.WriteJSON(filename, headers, rows)
	}
	if err != nil {
		log.Printf("Warning: %v", err)
		return err
	}

	return nil
}

// QueryUnbond returns a file with the unbond events in all the blocks
func QueryUnbond(addressQuery string, confirmedQuery bool, pendingQuery bool, outputFormat string) error {
	addressFilterString := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependQueryVaultsUnbond
	if len(addressQuery) > 0 {
		addressFilterString = fmt.Sprintf(bigquerytypes.QueryVaultsUnbondAddressFilter, addressQuery)
		filename = fmt.Sprintf("%s_%s", filename, addressQuery)
	}
	if confirmedQuery {
		filename = fmt.Sprintf("%s_%s", filename, "confirmed")
	}
	if pendingQuery {
		filename = fmt.Sprintf("%s_%s", filename, "pending")
	}

	headers, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryVaultsUnbond, addressFilterString, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if confirmedQuery || pendingQuery {
		_, confirmedRows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryVaultsUnbondConfirmedFilter, "", false)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// creating map storing the bond_id from the second query
		confirmedUnbondIDs := make(map[string]bool)
		for _, row := range confirmedRows {
			bondID := row[0]
			confirmedUnbondIDs[bondID] = true
		}

		// filtering rows from the first query by checking if bond_id exists
		filteredRows := [][]string{}
		bondIDRegex := regexp.MustCompile(`bond_id (\d+)`)
		for _, row := range rows {
			column3 := row[2]
			match := bondIDRegex.FindStringSubmatch(column3)
			if len(match) > 1 && bondIDRegex.MatchString(column3) {
				bondID := match[1]
				if confirmedQuery {
					if _, exists := confirmedUnbondIDs[bondID]; exists {
						filteredRows = append(filteredRows, row)
					}
				} else if pendingQuery {
					if _, exists := confirmedUnbondIDs[bondID]; !exists {
						filteredRows = append(filteredRows, row)
					}
				}
			}
		}
		rows = filteredRows
	}

	if outputFormat == "csv" {
		err = internal.WriteCSV(filename, headers, rows)
	} else {
		err = internal.WriteJSON(filename, headers, rows)
	}
	if err != nil {
		log.Printf("Warning: %v", err)
		return err
	}

	return nil
}

// QueryClaim returns a file with the claim events in all the blocks
func QueryClaim(addressQuery string, outputFormat string) error {
	addressFilterString := ""
	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependQueryVaultsClaim
	if len(addressQuery) > 0 {
		addressFilterString = fmt.Sprintf(bigquerytypes.QueryVaultsClaimAddressFilter, addressQuery)
		filename = fmt.Sprintf("%s_%s", filename, addressQuery)
	}

	headers, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryVaultsClaim, addressFilterString, true)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if outputFormat == "csv" {
		err = internal.WriteCSV(filename, headers, rows)
	} else {
		err = internal.WriteJSON(filename, headers, rows)
	}
	if err != nil {
		log.Printf("Warning: %v", err)
		return err
	}

	return nil
}
