package internal

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var dir = "output/"

func WriteCSV(fileName string, header []string, data [][]string) error {
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

	for i := range data {
		var csvRow []string
		for j := range data[i] {
			csvRow = append(csvRow, data[i][j])
		}
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}
