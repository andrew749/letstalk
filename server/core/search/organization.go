package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"letstalk/server/data"

	"github.com/olivere/elastic"
)

const ORGANIZATION_INDEX = "organizations"
const ORGANIZATION_TYPE = "organization"
const ORGANIZATION_SUGGESTER = "organization-suggester"

type Organization struct {
	Id              data.TOrganizationID  `json:"id"`
	Name            string                `json:"name"`
	Type            data.OrganizationType `json:"type"`
	IsUserGenerated bool                  `json:"isUserGenerated"`
	Suggest         SuggestInput          `json:"suggest"`
}

func (c *ClientWithContext) IndexOrganization(organization Organization) error {
	_, err := c.client.Index().
		Index(ORGANIZATION_INDEX).
		Type(ORGANIZATION_TYPE).
		Id(strconv.Itoa(int(organization.Id))).
		BodyJson(organization).
		Do(c.context)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientWithContext) CompletionSuggestionOrganizations(
	prefix string,
	size int,
) ([]Organization, error) {
	fuzzyOptions := elastic.NewFuzzyCompletionSuggesterOptions().
		EditDistance("AUTO").
		Transpositions(true)

	suggester := elastic.NewCompletionSuggester(ORGANIZATION_SUGGESTER).
		Prefix(prefix).
		Size(size).
		Field("suggest").
		SkipDuplicates(true).
		FuzzyOptions(fuzzyOptions)

	searchResult, err := c.client.Search().
		Index(ORGANIZATION_INDEX).
		Suggester(suggester).
		Do(c.context)
	if err != nil {
		return nil, err
	}

	organizations := make([]Organization, 0)
	if searchSuggestions, ok := searchResult.Suggest[ORGANIZATION_SUGGESTER]; ok {
		if len(searchSuggestions) > 0 {
			opts := searchSuggestions[0].Options
			for _, opt := range opts {
				var organization Organization
				err = json.Unmarshal(*opt.Source, &organization)
				if err != nil {
					return nil, err
				}
				organizations = append(organizations, organization)
			}
		}
	} else {
		return nil, errors.New(
			fmt.Sprintf("Completion results doesn't have suggestion %s", ORGANIZATION_SUGGESTER),
		)
	}

	return organizations, nil
}

// For use in backfill jobs
func (c *ClientWithContext) BulkIndexOrganizations(organizations []Organization) error {
	bulkRequest := c.client.Bulk()
	for _, organization := range organizations {
		req := elastic.NewBulkIndexRequest().
			Index(ORGANIZATION_INDEX).
			Type(ORGANIZATION_TYPE).
			Id(strconv.Itoa(int(organization.Id))).
			Doc(organization)
		bulkRequest = bulkRequest.Add(req)
	}
	res, err := bulkRequest.Do(c.context)
	if err != nil {
		return err
	}
	return consolidateBulkResponseErrors(res)
}

func (c *ClientWithContext) createOrganizationIndex() error {
	exists, err := c.client.IndexExists(ORGANIZATION_INDEX).Do(c.context)
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
    `, ORGANIZATION_TYPE)
		_, err := c.client.CreateIndex(ORGANIZATION_INDEX).BodyString(mapping).Do(c.context)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewOrganizationFromDataModel(dataOrganization data.Organization) Organization {
	return Organization{
		dataOrganization.Id,
		dataOrganization.Name,
		dataOrganization.Type,
		dataOrganization.IsUserGenerated,
		SuggestInput{[]string{dataOrganization.Name}, nil},
	}
}
