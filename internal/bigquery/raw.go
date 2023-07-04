package bigquery

import (
	"github.com/arhamchordia/chain-details/internal"
	bigquerytypes "github.com/arhamchordia/chain-details/types/bigquery"
	"log"
)

func RawQuery(RawQuery string, outputFormat string) error {
	headers, rows, err := internal.ExecuteQueryAndFetchRows(RawQuery, "", false)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if len(rows) == 0 {
		log.Fatalf("No rows returned by query")
	}

	filename := bigquerytypes.PrefixBigQuery + bigquerytypes.PrependQueryRaw
	if outputFormat == "csv" {
		err = internal.WriteCSV(filename, headers, rows)
	} else {
		err = internal.WriteJSON(filename, headers, rows)
	}
	if err != nil {
		return err
	}

	return nil
}
