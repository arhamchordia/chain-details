package internal

import (
	"context"
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

func (bq *BigQueryQuerier) ExecuteQuery(query string) (*bigquery.RowIterator, error) {
	q := bq.Client.Query(query)
	it, err := q.Read(bq.Ctx)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return nil, err
	}

	return it, nil
}
