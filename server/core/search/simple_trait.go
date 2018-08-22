package search

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"letstalk/server/data"

	"github.com/olivere/elastic"
	"github.com/romana/rlog"
)

const SIMPLE_TRAIT_INDEX = "simple_traits"
const SIMPLE_TRAIT_TYPE = "simple_trait"

type SimpleTrait struct {
	Id              data.TSimpleTraitID  `json:"id"`
	Name            string               `json:"name"`
	Type            data.SimpleTraitType `json:"type"`
	IsSensitive     bool                 `json:"isSensitive"`
	IsUserGenerated bool                 `json:"isUserGenerated"`
	Suggest         SuggestInput         `json:"suggest"`
}

func (c *RequestSearchClient) IndexSimpleTrait(trait SimpleTrait) error {
	_, err := c.client.Index().
		Index(SIMPLE_TRAIT_INDEX).
		Type(SIMPLE_TRAIT_TYPE).
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
		Index(SIMPLE_TRAIT_INDEX).
		Query(termQuery).
		Sort("id", true).
		Size(100).
		Do(c.request.Context())
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

func (c *RequestSearchClient) CompletionSuggestionSimpleTraits(
	prefix string,
	size int,
) ([]SimpleTrait, error) {
	fuzzyOptions := elastic.NewFuzzyCompletionSuggesterOptions().
		EditDistance("AUTO").
		Transpositions(true)

	suggesterName := "simple-traits-suggester"

	rlog.Info(prefix, size)

	suggester := elastic.NewCompletionSuggester(suggesterName).
		Text(prefix).
		Size(size).
		Field("suggest").
		SkipDuplicates(true).
		FuzzyOptions(fuzzyOptions)

	searchResult, err := c.client.Search().
		Index(SIMPLE_TRAIT_INDEX).
		Suggester(suggester).
		Do(c.request.Context())
	if err != nil {
		return nil, err
	}

	traits := make([]SimpleTrait, 0)
	if _, ok := searchResult.Suggest[suggesterName]; ok {
		searchSuggestions := searchResult.Suggest[suggesterName]
		if len(searchSuggestions) > 0 {
			opts := searchSuggestions[0].Options
			for _, opt := range opts {
				var trait SimpleTrait
				err = json.Unmarshal(*opt.Source, &trait)
				if err != nil {
					return nil, err
				}
				traits = append(traits, trait)
			}
		}
	}

	return traits, nil
}

// For use in backfill jobs
func BulkIndexSimpleTraits(es *elastic.Client, traits []SimpleTrait) error {
	bulkRequest := es.Bulk()
	for _, trait := range traits {
		req := elastic.NewBulkIndexRequest().
			Index(SIMPLE_TRAIT_INDEX).
			Type(SIMPLE_TRAIT_TYPE).
			Id(strconv.Itoa(int(trait.Id))).
			Doc(trait)
		bulkRequest = bulkRequest.Add(req)
	}
	res, err := bulkRequest.Do(context.Background())
	if err != nil {
		return err
	} else if len(res.Failed()) > 0 {
		return errors.New(fmt.Sprintf("More than 0 operations failed: %d, example: %s\n",
			len(res.Failed()),
			res.Failed()[0].Error.Reason))
	}
	return nil
}

func createSimpleTraitIndex(es *elastic.Client) error {
	exists, err := es.IndexExists(SIMPLE_TRAIT_INDEX).Do(context.Background())
	if err != nil {
		return err
	}

	if !exists {

		mapping := `
      {
        "mappings": {
          "simple_trait" : {
            "properties" : {
              "suggest" : {
                "type" : "completion"
              }
            }
          }
        }
      }
		`
		_, err := es.CreateIndex(SIMPLE_TRAIT_INDEX).BodyString(mapping).Do(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

func NewSimpleTraitFromDataModel(dataTrait data.SimpleTrait) SimpleTrait {
	return SimpleTrait{
		dataTrait.Id,
		dataTrait.Name,
		dataTrait.Type,
		dataTrait.IsSensitive,
		dataTrait.IsUserGenerated,
		SuggestInput{[]string{dataTrait.Name}, nil},
	}
}
