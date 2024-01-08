package bigquery

import (
	"github.com/arhamchordia/chain-details/internal"
	"github.com/arhamchordia/chain-details/internal/export"
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

	err = export.ExportFile(outputFormat, filename, headers, rows)
	if err != nil {
		return err
	}

	return nil
}
