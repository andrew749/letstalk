package search

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"letstalk/server/data"

	"github.com/olivere/elastic"
)

const SIMPLE_TRAIT_INDEX = "simple_traits"

type SimpleTrait struct {
	Id              data.TSimpleTraitID  `json:"id"`
	Name            string               `json:"name"`
	Type            data.SimpleTraitType `json:"type"`
	IsSensitive     bool                 `json:"isSensitive"`
	IsUserGenerated bool                 `json:"isUserGenerated"`
}

func InsertSimpleTrait(es *elastic.Client, trait SimpleTrait) error {
	_, err := es.Index().
		Index(SIMPLE_TRAIT_INDEX).
		Type("doc").
		Id(strconv.Itoa(int(trait.Id))).
		BodyJson(trait).
		Refresh("wait_for").
		Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func PrintAllSimpleTraits(es *elastic.Client) error {
	termQuery := elastic.NewTermQuery("isUserGenerated", true)

	searchResult, err := es.Search().
		Index("tweets").            // search in index "tweets"
		Query(termQuery).           // specify the query
		Sort("user.keyword", true). // sort by "user" field, ascending
		From(0).Size(10).           // take documents 0-9
		Pretty(true).               // pretty print request and response JSON
		Do(context.Background())    // execute
	if err != nil {
		return err
	}

	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	var ttyp SimpleTrait
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(SimpleTrait); ok {
			fmt.Println(t)
		}
	}
	return nil
}
