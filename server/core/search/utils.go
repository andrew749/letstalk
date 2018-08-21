package search

import (
	"context"

	"github.com/olivere/elastic"
)

func NewEsClient(addr string) (*elastic.Client, error) {
	return elastic.NewClient(elastic.SetURL(addr))
}

func CreateEsIndexes(client *elastic.Client) error {
	_, err := client.CreateIndex("simple_traits").Do(context.Background())
	return err
}
