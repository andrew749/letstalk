package search

import (
	"context"
	"net/http"

	"github.com/olivere/elastic"
)

func NewEsClient(addr string) (*elastic.Client, error) {
	return elastic.NewClient(elastic.SetURL(addr))
}

func CreateEsIndexes(client *elastic.Client) error {
	_, err := client.CreateIndex("simple_traits").Do(context.Background())
	return err
}

// Search client to be used within the request context
type RequestSearchClient struct {
	client  *elastic.Client
	request *http.Request
}

func NewSearchClient(client *elastic.Client, request *http.Request) *RequestSearchClient {
	return &RequestSearchClient{client, request}
}
