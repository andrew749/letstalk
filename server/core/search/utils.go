package search

import (
	"context"

	"github.com/olivere/elastic"
)

// Used for indexes that require completion suggestions
type SuggestInput struct {
	Input  []string `json:"input"`
	Weight *int     `json:"weight,omitempty"`
}

func NewEsClient(addr string) (*elastic.Client, error) {
	return elastic.NewClient(elastic.SetURL(addr))
}

// Search client to be used within the request context
type ClientWithContext struct {
	client  *elastic.Client
	context context.Context
}

func NewClientWithContext(client *elastic.Client, context context.Context) *ClientWithContext {
	return &ClientWithContext{client, context}
}

func (c *ClientWithContext) CreateEsIndexes() error {
	return c.createSimpleTraitIndex()
}
