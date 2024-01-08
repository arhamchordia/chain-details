package export

import (
	"fmt"
)

func WriteGoogleSpreadsheet(fileName string, header []string, rows [][]string) error {
	if len(rows) == 0 {
		return fmt.Errorf("No rows to write to the Google Spreadsheet")
	}

	// TODO

	return nil
}
