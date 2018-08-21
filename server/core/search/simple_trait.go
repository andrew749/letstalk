package search

import (
	"fmt"
	"reflect"
	"strconv"

	"letstalk/server/data"

	"github.com/olivere/elastic"
	"github.com/romana/rlog"
)

const SIMPLE_TRAIT_INDEX = "simple_traits"

type SimpleTrait struct {
	Id              data.TSimpleTraitID  `json:"id"`
	Name            string               `json:"name"`
	Type            data.SimpleTraitType `json:"type"`
	IsSensitive     bool                 `json:"isSensitive"`
	IsUserGenerated bool                 `json:"isUserGenerated"`
}

func (c *RequestSearchClient) IndexSimpleTrait(trait SimpleTrait) error {
	_, err := c.client.Index().
		Index(SIMPLE_TRAIT_INDEX).
		Type("doc").
		Id(strconv.Itoa(int(trait.Id))).
		BodyJson(trait).
		Do(c.request.Context())
	if err != nil {
		return err
	}
	return nil
}

func (c *RequestSearchClient) PrintAllSimpleTraits() error {
	termQuery := elastic.NewTermQuery("isUserGenerated", true)

	searchResult, err := c.client.Search().
		Index(SIMPLE_TRAIT_INDEX). // search in index "tweets"
		Query(termQuery).          // specify the query
		Sort("id", true).          // sort by "user" field, ascending
		Size(100).                 // take documents 0-9
		Pretty(true).              // pretty print request and response JSON
		Do(c.request.Context())    // execute
	if err != nil {
		return err
	}

	rlog.Info(fmt.Sprintf("Query took %d milliseconds\n", searchResult.TookInMillis))
	var ttyp SimpleTrait
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(SimpleTrait); ok {
			rlog.Info(t)
		}
	}
	return nil
}
