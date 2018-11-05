package search

import (
	"encoding/json"
	"fmt"

	"letstalk/server/data"

	"github.com/getsentry/raven-go"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
)

type MultiTraitID int
type MultiTraitType string

const MULTI_TRAIT_INDEX = "multi_traits"
const MULTI_TRAIT_TYPE = "multi_trait"

const (
	MULTI_TRAIT_TYPE_COHORT       MultiTraitType = "COHORT"
	MULTI_TRAIT_TYPE_POSITION     MultiTraitType = "POSITION"
	MULTI_TRAIT_TYPE_SIMPLE_TRAIT MultiTraitType = "SIMPLE_TRAIT"
	MULTI_TRAIT_TYPE_GROUP        MultiTraitType = "GROUP"
)

type MultiTrait struct {
	TraitName string         `json:"traitName"`
	TraitType MultiTraitType `json:"traitType"`
}

type CohortMultiTrait struct {
	CohortId     data.TCohortID `json:"cohortId"`
	ProgramId    string         `json:"programId"`
	ProgramName  string         `json:"programName"`
	GradYear     uint           `json:"gradYear"`
	IsCoop       bool           `json:"isCoop"`
	SequenceId   *string        `json:"sequenceId,omitempty"`
	SequenceName *string        `json:"sequenceName,omitempty"`
	MultiTrait
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

type GroupMultiTrait struct {
	GroupId   data.TGroupID `json:"groupId"`
	GroupName string        `json:"groupName"`
	MultiTrait
}

func (c *ClientWithContext) indexMultiTrait(id string, trait interface{}) error {
	if !isMultiTrait(trait) {
		return errors.WithStack(errors.New(fmt.Sprintf("Invalid type of trait %T", trait)))
	}

	_, err := c.client.Index().
		Index(MULTI_TRAIT_INDEX).
		Type(MULTI_TRAIT_TYPE).
		Id(id).
		BodyJson(trait).
		Do(c.context)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientWithContext) IndexCohortMultiTrait(userCohort *data.UserCohort) error {
	id, trait := NewMultiTraitFromUserCohort(userCohort)
	return c.indexMultiTrait(id, trait)
}

func (c *ClientWithContext) IndexPositionMultiTrait(userPosition *data.UserPosition) error {
	id, trait := NewMultiTraitFromUserPosition(userPosition)
	return c.indexMultiTrait(id, trait)
}

func (c *ClientWithContext) IndexSimpleTraitMultiTrait(userSimpleTrait *data.UserSimpleTrait) error {
	id, trait := NewMultiTraitFromUserSimpleTrait(userSimpleTrait)
	return c.indexMultiTrait(id, trait)
}

func (c *ClientWithContext) IndexGroupMultiTrait(userGroup *data.UserGroup) error {
	id, trait := NewMultiTraitFromUserGroup(userGroup)
	return c.indexMultiTrait(id, trait)
}

// Returns id for the document and the CohortMultiTrait struct
func NewMultiTraitFromUserCohort(userCohort *data.UserCohort) (string, *CohortMultiTrait) {
	cohort := userCohort.Cohort
	id := fmt.Sprintf("%s-%d", MULTI_TRAIT_TYPE_COHORT, cohort.CohortId)

	var traitName string
	if cohort.SequenceName == nil || *cohort.SequenceId == "OTHER" {
		traitName = fmt.Sprintf("%s %d", cohort.ProgramName, cohort.GradYear)
	} else {
		traitName = fmt.Sprintf("%s %d %s", cohort.ProgramName, cohort.GradYear, *cohort.SequenceName)
	}

	cohortMultiTrait := &CohortMultiTrait{
		CohortId:     cohort.CohortId,
		ProgramId:    cohort.ProgramId,
		ProgramName:  cohort.ProgramName,
		GradYear:     cohort.GradYear,
		IsCoop:       cohort.IsCoop,
		SequenceId:   cohort.SequenceId,
		SequenceName: cohort.SequenceName,
		MultiTrait: MultiTrait{
			TraitName: traitName,
			TraitType: MULTI_TRAIT_TYPE_COHORT,
		},
	}
	return id, cohortMultiTrait
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

// Returns id for the document and the GroupMultiTrait struct
func NewMultiTraitFromUserGroup(trait *data.UserGroup) (string, *GroupMultiTrait) {
	id := fmt.Sprintf("%s-%s", MULTI_TRAIT_TYPE_GROUP, trait.GroupId)
	groupMultiTrait := &GroupMultiTrait{
		GroupId:   trait.GroupId,
		GroupName: trait.GroupName,
		MultiTrait: MultiTrait{
			TraitName: trait.GroupName,
			TraitType: MULTI_TRAIT_TYPE_GROUP,
		},
	}
	return id, groupMultiTrait
}

func parseMultiTraitHit(hit *elastic.SearchHit) (interface{}, error) {
	var (
		fieldMap     map[string]interface{}
		typeField    interface{}
		typeFieldStr string
		ok           bool
		res          interface{}
	)
	err := json.Unmarshal(*hit.Source, &fieldMap)
	if err != nil {
		return nil, err
	}
	if typeField, ok = fieldMap["traitType"]; !ok {
		return nil, errors.New("Missing trait type in hit")
	}
	if typeFieldStr, ok = typeField.(string); !ok {
		return nil, errors.New("Trait type is not a string")
	}

	switch MultiTraitType(typeFieldStr) {
	case MULTI_TRAIT_TYPE_POSITION:
		var pos PositionMultiTrait
		err = json.Unmarshal(*hit.Source, &pos)
		if err != nil {
			return nil, err
		}
		res = &pos
	case MULTI_TRAIT_TYPE_SIMPLE_TRAIT:
		var trait SimpleTraitMultiTrait
		err = json.Unmarshal(*hit.Source, &trait)
		if err != nil {
			return nil, err
		}
		res = &trait
	case MULTI_TRAIT_TYPE_COHORT:
		var cohort CohortMultiTrait
		err = json.Unmarshal(*hit.Source, &cohort)
		if err != nil {
			return nil, err
		}
		res = &cohort
	case MULTI_TRAIT_TYPE_GROUP:
		var group GroupMultiTrait
		err = json.Unmarshal(*hit.Source, &group)
		if err != nil {
			return nil, err
		}
		res = &group
	default:
		return nil, errors.New("Trait type is not a string")
	}
	return res, nil
}

// Queries multi traits by name, using the autocomplete analyzer on the traitName field.
func (c *ClientWithContext) QueryMultiTraitsByName(prefix string, size int) ([]interface{}, error) {
	matchQuery := elastic.NewMatchQuery("traitName", prefix)

	res, err := c.client.Search().
		Index(MULTI_TRAIT_INDEX).
		Type(MULTI_TRAIT_TYPE).
		Query(matchQuery).
		Size(size).
		Do(c.context)
	if err != nil {
		return nil, err
	}

	traits := make([]interface{}, 0)
	for _, hit := range res.Hits.Hits {
		trait, err := parseMultiTraitHit(hit)
		if err != nil {
			raven.CaptureError(err, nil)
			// Ignore bad results, but record so that we can fix them.
			continue
		}
		traits = append(traits, trait)
	}

	return traits, nil
}

// Checks whether type of object is one for the valid multi trait types
// Must be a pointer
func isMultiTrait(obj interface{}) bool {
	_, ok1 := obj.(*CohortMultiTrait)
	_, ok2 := obj.(*PositionMultiTrait)
	_, ok3 := obj.(*SimpleTraitMultiTrait)
	_, ok4 := obj.(*GroupMultiTrait)
	return ok1 || ok2 || ok3 || ok4
}

// For use in backfill jobs
func (c *ClientWithContext) BulkIndexMultiTraits(traits map[string]interface{}) error {
	for _, trait := range traits {
		if !isMultiTrait(trait) {
			return errors.New(fmt.Sprintf("Invalid type of trait %T", trait))
		}
	}

	bulkRequest := c.client.Bulk()
	for id, trait := range traits {
		req := elastic.NewBulkIndexRequest().
			Index(MULTI_TRAIT_INDEX).
			Type(MULTI_TRAIT_TYPE).
			Id(id).
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
