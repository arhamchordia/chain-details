package internal

import (
	"encoding/csv"
	"os"
)

func WriteCSV(fileName string, header []string, data [][]string) error {
	outputFile, err := os.Create(fileName + ".csv")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		panic(err)
	}

	for i := range data {
		var csvRow []string
		for j := range data[i] {
			csvRow = append(csvRow, data[i][j])
		}
		if err := writer.Write(csvRow); err != nil {
			panic(err)
		}
	}
	return nil
}
