package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func WriteCSV(fileName string, header []string, rows [][]string) error {
	if len(rows) == 0 {
		return fmt.Errorf("No rows to write to the CSV")
	}

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	timestamp := time.Now().Unix()
	filePath := filepath.Join(outputDir, fileName+"_"+strconv.FormatInt(timestamp, 10)+".csv")

	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		return err
	}

	for i := range rows {
		var csvRow []string
		csvRow = append(csvRow, rows[i]...)
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}
