package export

import (
	"log"
)

var outputDir = "output/"

func ExportFile(outputFormat string, filename string, headers []string, rows [][]string) error {
	var err error

	if outputFormat == "csv" {
		err = WriteCSV(filename, headers, rows)
	} else {
		err = WriteJSON(filename, headers, rows)
	}
	if err != nil {
		log.Printf("Warning: %v", err)
		return err
	}

	return nil
}
