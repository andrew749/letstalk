package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"letstalk/server/data"

	"github.com/olivere/elastic"
)

const ROLE_INDEX = "roles"
const ROLE_TYPE = "role"
const ROLE_SUGGESTER = "role-suggester"

type Role struct {
	Id              data.TRoleID `json:"id"`
	Name            string       `json:"name"`
	IsUserGenerated bool         `json:"isUserGenerated"`
	Suggest         SuggestInput `json:"suggest"`
}

func (c *ClientWithContext) IndexRole(role Role) error {
	_, err := c.client.Index().
		Index(ROLE_INDEX).
		Type(ROLE_TYPE).
		Id(strconv.Itoa(int(role.Id))).
		BodyJson(role).
		Do(c.context)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientWithContext) CompletionSuggestionRoles(
	prefix string,
	size int,
) ([]Role, error) {
	fuzzyOptions := elastic.NewFuzzyCompletionSuggesterOptions().
		EditDistance("AUTO").
		Transpositions(true)

	suggester := elastic.NewCompletionSuggester(ROLE_SUGGESTER).
		Prefix(prefix).
		Size(size).
		Field("suggest").
		SkipDuplicates(true).
		FuzzyOptions(fuzzyOptions)

	searchResult, err := c.client.Search().
		Index(ROLE_INDEX).
		Suggester(suggester).
		Do(c.context)
	if err != nil {
		return nil, err
	}

	roles := make([]Role, 0)
	if searchSuggestions, ok := searchResult.Suggest[ROLE_SUGGESTER]; ok {
		if len(searchSuggestions) > 0 {
			opts := searchSuggestions[0].Options
			for _, opt := range opts {
				var role Role
				err = json.Unmarshal(*opt.Source, &role)
				if err != nil {
					return nil, err
				}
				roles = append(roles, role)
			}
		}
	} else {
		return nil, errors.New(
			fmt.Sprintf("Completion results doesn't have suggestion %s", ROLE_SUGGESTER),
		)
	}

	return roles, nil
}

// For use in backfill jobs
func (c *ClientWithContext) BulkIndexRoles(roles []Role) error {
	bulkRequest := c.client.Bulk()
	for _, role := range roles {
		req := elastic.NewBulkIndexRequest().
			Index(ROLE_INDEX).
			Type(ROLE_TYPE).
			Id(strconv.Itoa(int(role.Id))).
			Doc(role)
		bulkRequest = bulkRequest.Add(req)
	}
	res, err := bulkRequest.Do(c.context)
	if err != nil {
		return err
	}
	return consolidateBulkResponseErrors(res)
}

func (c *ClientWithContext) createRoleIndex() error {
	exists, err := c.client.IndexExists(ROLE_INDEX).Do(c.context)
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
    `, ROLE_TYPE)
		_, err := c.client.CreateIndex(ROLE_INDEX).BodyString(mapping).Do(c.context)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewRoleFromDataModel(dataRole data.Role) Role {
	return Role{
		dataRole.Id,
		dataRole.Name,
		dataRole.IsUserGenerated,
		SuggestInput{[]string{dataRole.Name}, nil},
	}
}
