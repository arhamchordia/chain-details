package bigquery

import (
	"github.com/arhamchordia/chain-details/internal"
	bigquerytypes "github.com/arhamchordia/chain-details/types/bigquery"
	"log"
)

func RawQuery(RawQuery string) error {
	headers, rows, err := internal.ExecuteQueryAndFetchRows(RawQuery, "", false)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if len(rows) == 0 {
		log.Fatalf("No rows returned by query")
	}

	filename := bigquerytypes.PrefixBigQuery + "raw"
	err = internal.WriteCSV(filename, headers, rows)
	if err != nil {
		return err
	}

	return nil
}
