package search

import (
	"net/http"

	"github.com/olivere/elastic"
)

type SuggestInput struct {
	Input  []string `json:"input"`
	Weight *int     `json:"weight,omitempty"`
}

func NewEsClient(addr string) (*elastic.Client, error) {
	return elastic.NewClient(elastic.SetURL(addr))
}

func CreateEsIndexes(client *elastic.Client) error {
	return createSimpleTraitIndex(client)
}

// Search client to be used within the request context
type RequestSearchClient struct {
	client  *elastic.Client
	request *http.Request
}

func NewSearchClient(client *elastic.Client, request *http.Request) *RequestSearchClient {
	return &RequestSearchClient{client, request}
}
