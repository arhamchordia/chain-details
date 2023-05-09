package internal

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var dir = "output/"

func WriteCSV(fileName string, header []string, rows [][]string) error {
	if len(rows) == 0 {
		return fmt.Errorf("No rows to write to the CSV")
	}

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	timestamp := time.Now().Unix()
	filePath := filepath.Join(dir, fileName+"_"+strconv.FormatInt(timestamp, 10)+".csv")

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
		for j := range rows[i] {
			csvRow = append(csvRow, rows[i][j])
		}
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}
