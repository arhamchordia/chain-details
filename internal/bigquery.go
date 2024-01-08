package internal

import (
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
)

type BigQueryQuerier struct {
	Client *bigquery.Client
	Ctx    context.Context
}

func NewBigQueryQuerier() (*BigQueryQuerier, error) {
	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT_ID")
	if projectID == "" {
		log.Fatal("Environment variable GOOGLE_CLOUD_PROJECT_ID must be set")
	}

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create BigQuery client: %v", err)
	}

	return &BigQueryQuerier{
		Client: client,
		Ctx:    ctx,
	}, nil
}

func (bq *BigQueryQuerier) Close() {
	err := bq.Client.Close()
	if err != nil {
		return
	}
}

func ExecuteQueryAndFetchRows(query, addressFilter string, applyAddressFilter bool) ([]string, [][]string, error) {
	bqQuerier, _ := NewBigQueryQuerier()

	if applyAddressFilter {
		query = fmt.Sprintf(query, addressFilter)
	}

	it, err := bqQuerier.executeQuery(query)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to execute BigQuery query: %v", err)
	}
	defer bqQuerier.Close()

	var headers []string
	var rows [][]string

	for {
		var row []bigquery.Value
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to iterate results: %v", err)
		}

		// Extract the headers if they haven't been extracted yet
		if headers == nil {
			headers = make([]string, len(row))
			for i, schemaField := range it.Schema {
				headers[i] = schemaField.Name
			}
		}

		var csvRow []string
		for _, val := range row {
			csvRow = append(csvRow, fmt.Sprintf("%v", val))
		}
		rows = append(rows, csvRow)
	}

	if len(rows) == 0 {
		return nil, nil, fmt.Errorf("No rows returned by query")
	}

	return headers, rows, nil
}

func (bq *BigQueryQuerier) executeQuery(query string) (*bigquery.RowIterator, error) {
	q := bq.Client.Query(query)
	it, err := q.Read(bq.Ctx)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return nil, err
	}

	return it, nil
}
