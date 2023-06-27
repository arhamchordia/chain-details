package vaults

import (
	"fmt"
	"github.com/arhamchordia/chain-details/internal"
	bigquerytypes "github.com/arhamchordia/chain-details/types/bigquery"
	"log"
	"regexp"
	"time"
)

//TODO:
// - Remove distinct from BondBefore query in order to be able ask for block_timestamp and OSMO bonding amount per-bond
// - Merge all the functions in one to post-compute between different actions as bond, unbond and claim
// - - Query all the bonds before 24h
// - - Query all the bonds in the last 24h
// - - Query all the unbonds before 24h
// - - Query all the unbonds in the last 24
// - - Query all the claims before 24h
// - - Query all the claims in the last 24h
// - - - Choose if we want to convert OSMO to Shares during bonding, or Shares to OSMO during unbond/claim
// - - - Use the converted value to post compute between bonds and unbonds to know which user exited completely and which one not
// - - - Use the claims to check if in the effective claim timestamp the user had pending unbonds to know if the claim was about all the unbonds or not

type UpdateIndex struct {
	VaultTokenBalance  string
	IngestionTimestamp time.Time
}

// QueryDailyReport returns a file with the last 24h statistics of a given vaultAddress
func QueryDailyReport(addressQuery string, outputFormat string) error {
	if len(addressQuery) == 0 {
		log.Fatal("Vault address to query is mandatory")
	}

	filename := fmt.Sprintf("%s_%s", bigquerytypes.PrefixBigQuery+bigquerytypes.PrependQueryDailyReport, addressQuery)

	rewardsUpdateUser, err := queryDailyReportRewardsUpdateUser(addressQuery)
	if err != nil {
		return err
	}

	//fmt.Println(rewardsUpdateUser)
	fmt.Println("RESULT LENGTH", len(rewardsUpdateUser))

	// TODO: Headers and rows
	headers := []string{
		// Bond
		"bond_new_users_count", "bond_new_users_amount", "bond_old_users_count", "bond_old_users_amount",
		// Unbond
		"unbond_new_users_count", "unbond_new_users_amount", "unbond_old_users_count", "unbond_old_users_amount",
		// General
		"bond_total_users",
	}
	rows := [][]string{
		//{
		//	// Bond
		//	strconv.Itoa(bondNewUsersCount), strconv.Itoa(bondNewUsersAmount), strconv.Itoa(bondOldUsersCount), strconv.Itoa(bondOldUsersAmount),
		//	// Unbond
		//	strconv.Itoa(unbondNewUsersCount), strconv.Itoa(unbondNewUsersAmount), strconv.Itoa(unbondOldUsersCount), strconv.Itoa(unbondOldUsersAmount),
		//	// General
		//	strconv.Itoa(bondTotalUsers),
		//},
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

//func queryDailyReportRewardsUpdateUser(addressQuery string) (map[string]UpdateIndex, error) {
//	_, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryDailyReportRewardsUpdateUser, "", false)
//	if err != nil {
//		log.Fatalf("%v", err)
//	}
//
//	Addresses := make(map[string]UpdateIndex)
//
//	for _, row := range rows {
//		attributeValue := row[2]
//		ingestionTimestampStr := row[3]
//
//		fmt.Println("ATTR PRE EXTRACT >>>>>>>>>.", attributeValue)
//		attributes, err := extractAttributes(attributeValue)
//		if err != nil {
//			log.Fatalf("%v", err)
//		}
//		fmt.Println("ATTR POST EXTRACT >>>>>>>>>.", attributes)
//
//		contractAddress, caOk := attributes["_contract_address"]
//		user, uOk := attributes["user"]
//		vaultTokenBalance, vtbOk := attributes["vault_token_balance"]
//
//		if caOk && uOk && vtbOk && contractAddress == addressQuery {
//			ingestionTimestamp, err := time.Parse("2006-01-02 15:04:05 -0700 MST", ingestionTimestampStr)
//			if err != nil {
//				//fmt.Println(err)
//			}
//			Addresses[user] = UpdateIndex{
//				VaultTokenBalance:  vaultTokenBalance,
//				IngestionTimestamp: ingestionTimestamp,
//			}
//		}
//	}
//
//	return Addresses, nil
//}

func queryDailyReportRewardsUpdateUser(addressQuery string) (map[string][]UpdateIndex, error) {
	_, rows, err := internal.ExecuteQueryAndFetchRows(bigquerytypes.QueryDailyReportRewardsUpdateUser, "", false)
	if err != nil {
		log.Fatalf("%v", err)
	}

	Addresses := make(map[string][]UpdateIndex)

	for _, row := range rows {
		attributeValue := row[2]
		ingestionTimestampStr := row[3]

		//fmt.Println("ATTR PRE EXTRACT >>>>>>>>>.", attributeValue)
		attributes, err := extractAttributes(attributeValue)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// iter over all values of contract_address
		for i, contractAddress := range attributes["_contract_address"] {
			// if contract address matches the addressQuery (basic vault address) and usr and vault_token_balance attributes for this index exist
			// TODO check, this is wrong as we are not filtering out Events that are not action "wasm" and have an Attribute key "update_token_index"
			if contractAddress == addressQuery && i < len(attributes["user"]) && i < len(attributes["vault_token_balance"]) {
				user := attributes["user"][i]
				vaultTokenBalance := attributes["vault_token_balance"][i]

				ingestionTimestamp, err := time.Parse("2006-01-02 15:04:05 -0700 MST", ingestionTimestampStr)
				if err != nil {
					// Handle error
				}

				Addresses[user] = append(Addresses[user], UpdateIndex{
					VaultTokenBalance:  vaultTokenBalance,
					IngestionTimestamp: ingestionTimestamp,
				})
			}
		}
	}

	//fmt.Println(Addresses)

	return Addresses, nil
}

//func extractAttributes(input string) (map[string][]string, error) {
//	fmt.Println("INPUT >>>", input)
//
//	re := regexp.MustCompile(`Attribute\s\{\skey:\s"([^"]+)",\svalue:\s"([^"]+)"\s\}`)
//	matches := re.FindAllStringSubmatch(input, -1)
//
//	attributes := make(map[string][]string)
//
//	//fmt.Println("MATCHES >>>", matches)
//	for _, match := range matches {
//		//fmt.Println("MATCH >>>", match)
//		if len(match) != 3 {
//			return nil, fmt.Errorf("invalid key value pair: %s", match)
//		}
//		key := match[1]
//		value := match[2]
//		attributes[key] = append(attributes[key], value)
//	}
//
//	return attributes, nil
//}

func extractAttributes(input string) (map[string][]string, error) {
	// filter out all the events that are not related to ty:"wasm" && key:"action" value:"update_user_index" TODO not working properly
	re := regexp.MustCompile(`Event\s\{\sty:\s"wasm",\sattributes:\s\[(?s:.*?)Attribute\s\{\skey:\s"action",\svalue:\s"update_user_index"(?s:.*?)]\s\}`)

	matches := re.FindAllStringSubmatch(input, -1)
	fmt.Println("EVENTS >>> ", matches) // TODO: here there are still undesired events

	attributes := make(map[string][]string)

	// get key value pairs for each match
	for _, match := range matches {
		subRe := regexp.MustCompile(`Attribute\s\{\skey:\s"([^"]+)",\svalue:\s"([^"]+)"\s\}`)
		subMatches := subRe.FindAllStringSubmatch(match[0], -1)

		for _, subMatch := range subMatches {
			if len(subMatch) != 3 {
				return nil, fmt.Errorf("invalid key value pair: %s", subMatch)
			}
			key := subMatch[1]
			value := subMatch[2]
			attributes[key] = append(attributes[key], value)
		}
	}

	return attributes, nil
}
