package internal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func WriteJSON(fileName string, headers []string, rows [][]string) error {
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		return err
	}

	timestamp := time.Now().Unix()
	filePath := filepath.Join(outputDir, fileName+"_"+strconv.FormatInt(timestamp, 10)+".json")

	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	data := make([]map[string]string, len(rows))

	for i, row := range rows {
		rowMap := make(map[string]string)
		for j, cell := range row {
			rowMap[headers[j]] = cell
		}
		data[i] = rowMap
	}

	encoder := json.NewEncoder(outputFile)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}
