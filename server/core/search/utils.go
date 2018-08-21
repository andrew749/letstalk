package search

import (
	"context"

	"github.com/olivere/elastic"
)

func NewEsClient(addr string) (*elastic.Client, error) {
	return elastic.NewClient(elastic.SetURL(addr))
}

func CreateEsIndexes(client *elastic.Client) error {
	_, err := client.CreateIndex(SIMPLE_TRAIT_INDEX).Do(context.Background())
	return err
}
