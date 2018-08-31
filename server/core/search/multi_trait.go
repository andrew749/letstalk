package search

import (
	"errors"
	"fmt"

	"letstalk/server/data"

	"github.com/olivere/elastic"
)

type MultiTraitID int
type MultiTraitType string

const MULTI_TRAIT_INDEX = "multi_traits"
const MULTI_TRAIT_TYPE = "multi_trait"

const (
	MULTI_TRAIT_TYPE_COHORT       MultiTraitType = "COHORT"
	MULTI_TRAIT_TYPE_POSITION     MultiTraitType = "POSITION"
	MULTI_TRAIT_TYPE_SIMPLE_TRAIT MultiTraitType = "SIMPLE_TRAIT"
)

type MultiTrait struct {
	TraitName string         `json:"traitName"`
	TraitType MultiTraitType `json:"traitType"`
}

type PositionMultiTrait struct {
	RoleId           data.TRoleID          `json:"roleId"`
	RoleName         string                `json:"roleName"`
	OrganizationId   data.TOrganizationID  `json:"organizationId"`
	OrganizationName string                `json:"organizationName"`
	OrganizationType data.OrganizationType `json:"organizationType"`
	MultiTrait
}

type SimpleTraitMultiTrait struct {
	SimpleTraitId          data.TSimpleTraitID  `json:"simpleTraitId"`
	SimpleTraitName        string               `json:"simpleTraitName"`
	SimpleTraitType        data.SimpleTraitType `json:"simpleTraitType"`
	SimpleTraitIsSensitive bool                 `json:"simpleTraitIsSensitive"`
	MultiTrait
}

// Returns id for the document and the PositionMultiTrait struct
func NewMultiTraitFromUserPosition(pos *data.UserPosition) (string, *PositionMultiTrait) {
	id := fmt.Sprintf("%s-%d-%d", MULTI_TRAIT_TYPE_POSITION, pos.RoleId, pos.OrganizationId)
	posMultiTrait := &PositionMultiTrait{
		RoleId:           pos.RoleId,
		RoleName:         pos.RoleName,
		OrganizationId:   pos.OrganizationId,
		OrganizationName: pos.OrganizationName,
		OrganizationType: pos.OrganizationType,
		MultiTrait: MultiTrait{
			TraitName: fmt.Sprintf("%s at %s", pos.RoleName, pos.OrganizationName),
			TraitType: MULTI_TRAIT_TYPE_POSITION,
		},
	}
	return id, posMultiTrait
}

// Returns id for the document and the SimpleTraitMultiTrait struct
func NewMultiTraitFromUserSimpleTrait(trait *data.UserSimpleTrait) (string, *SimpleTraitMultiTrait) {
	id := fmt.Sprintf("%s-%d", MULTI_TRAIT_TYPE_SIMPLE_TRAIT, trait.SimpleTraitId)
	traitMultiTrait := &SimpleTraitMultiTrait{
		SimpleTraitId:          trait.SimpleTraitId,
		SimpleTraitName:        trait.SimpleTraitName,
		SimpleTraitType:        trait.SimpleTraitType,
		SimpleTraitIsSensitive: trait.SimpleTraitIsSensitive,
		MultiTrait: MultiTrait{
			TraitName: trait.SimpleTraitName,
			TraitType: MULTI_TRAIT_TYPE_SIMPLE_TRAIT,
		},
	}
	return id, traitMultiTrait
}

// Checks whether type of object is one for the valid multi trait types
// Must be a pointer
func isMultiTrait(obj interface{}) bool {
	_, ok1 := obj.(*PositionMultiTrait)
	_, ok2 := obj.(*SimpleTraitMultiTrait)
	return ok1 || ok2
}

// For use in backfill jobs
func (c *ClientWithContext) BulkIndexMultiTraits(ids []string, traits []interface{}) error {
	if len(ids) != len(traits) {
		return errors.New(fmt.Sprintf(
			"Length of ids (%d) and traits (%d) must be equal",
			len(ids),
			len(traits),
		))
	}
	for _, trait := range traits {
		if !isMultiTrait(trait) {
			return errors.New(fmt.Sprintf("Invalid type of trait %T", trait))
		}
	}

	bulkRequest := c.client.Bulk()
	for i, trait := range traits {
		req := elastic.NewBulkIndexRequest().
			Index(MULTI_TRAIT_INDEX).
			Type(MULTI_TRAIT_TYPE).
			Id(ids[i]).
			Doc(trait)
		bulkRequest = bulkRequest.Add(req)
	}
	res, err := bulkRequest.Do(c.context)
	if err != nil {
		return err
	}
	return consolidateBulkResponseErrors(res)
}

func (c *ClientWithContext) createMultiTraitIndex() error {
	exists, err := c.client.IndexExists(MULTI_TRAIT_INDEX).Do(c.context)
	if err != nil {
		return err
	}

	if !exists {
		mapping := fmt.Sprintf(`
      {
        "settings": {
          "number_of_shards": 1,
          "analysis": {
            "filter": {
              "autocomplete_filter": {
                "type":     "edge_ngram",
                "min_gram": 1,
                "max_gram": 20
              }
            },
            "analyzer": {
              "autocomplete": {
                "type":      "custom",
                "tokenizer": "standard",
                "filter": [
                    "lowercase",
                    "autocomplete_filter"
                ]
              }
            }
          }
        },
        "mappings": {
          "%s" : {
            "properties" : {
              "traitName": {
                "type":     "text",
                "analyzer": "autocomplete"
              }
            }
          }
        }
      }
    `, MULTI_TRAIT_TYPE)
		_, err := c.client.CreateIndex(MULTI_TRAIT_INDEX).BodyString(mapping).Do(c.context)
		if err != nil {
			return err
		}
	}
	return nil
}
