package search

import (
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
const SIMPLE_TRAIT_SUGGESTER = "simple-trait-suggester"

type SimpleTrait struct {
	Id              data.TSimpleTraitID  `json:"id"`
	Name            string               `json:"name"`
	Type            data.SimpleTraitType `json:"type"`
	IsSensitive     bool                 `json:"isSensitive"`
	IsUserGenerated bool                 `json:"isUserGenerated"`
	Suggest         SuggestInput         `json:"suggest"`
}

func (c *ClientWithContext) IndexSimpleTrait(trait SimpleTrait) error {
	_, err := c.client.Index().
		Index(SIMPLE_TRAIT_INDEX).
		Type(SIMPLE_TRAIT_TYPE).
		Id(strconv.Itoa(int(trait.Id))).
		BodyJson(trait).
		Do(c.context)
	if err != nil {
		return err
	}
	return nil
}

// Only here as an example of how to do a search.
// Can remove this if we end up writing a search.
func (c *ClientWithContext) PrintAllSimpleTraits() error {
	termQuery := elastic.NewTermQuery("isUserGenerated", true)

	searchResult, err := c.client.Search().
		Index(SIMPLE_TRAIT_INDEX).
		Query(termQuery).
		Sort("id", true).
		Size(100).
		Do(c.context)
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

func (c *ClientWithContext) CompletionSuggestionSimpleTraits(
	prefix string,
	size int,
) ([]SimpleTrait, error) {
	fuzzyOptions := elastic.NewFuzzyCompletionSuggesterOptions().
		EditDistance("AUTO").
		Transpositions(true)

	suggester := elastic.NewCompletionSuggester(SIMPLE_TRAIT_SUGGESTER).
		Prefix(prefix).
		Size(size).
		Field("suggest").
		SkipDuplicates(true).
		FuzzyOptions(fuzzyOptions)

	searchResult, err := c.client.Search().
		Index(SIMPLE_TRAIT_INDEX).
		Suggester(suggester).
		Do(c.context)
	if err != nil {
		return nil, err
	}

	rlog.Debugf("%#v", searchResult.Suggest)
	traits := make([]SimpleTrait, 0)
	if searchSuggestions, ok := searchResult.Suggest[SIMPLE_TRAIT_SUGGESTER]; ok {
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
	} else {
		return nil, errors.New(
			fmt.Sprintf("Completion results doesn't have suggestion %s", SIMPLE_TRAIT_SUGGESTER),
		)
	}

	return traits, nil
}

// For use in backfill jobs
func (c *ClientWithContext) BulkIndexSimpleTraits(traits []SimpleTrait) error {
	bulkRequest := c.client.Bulk()
	for _, trait := range traits {
		req := elastic.NewBulkIndexRequest().
			Index(SIMPLE_TRAIT_INDEX).
			Type(SIMPLE_TRAIT_TYPE).
			Id(strconv.Itoa(int(trait.Id))).
			Doc(trait)
		bulkRequest = bulkRequest.Add(req)
	}
	res, err := bulkRequest.Do(c.context)
	if err != nil {
		return err
	}
	return consolidateBulkResponseErrors(res)
}

func (c *ClientWithContext) createSimpleTraitIndex() error {
	exists, err := c.client.IndexExists(SIMPLE_TRAIT_INDEX).Do(c.context)
	if err != nil {
		return err
	}

	if !exists {
		mapping := fmt.Sprintf(`
      {
        "settings": {
          "number_of_shards": 1
        },
        "mappings": {
          "%s" : {
            "properties" : {
              "suggest" : {
                "type" : "completion"
              }
            }
          }
        }
      }
    `, SIMPLE_TRAIT_TYPE)
		_, err := c.client.CreateIndex(SIMPLE_TRAIT_INDEX).BodyString(mapping).Do(c.context)
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
