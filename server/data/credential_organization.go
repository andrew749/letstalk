package data

import (
	"fmt"
)

type CredentialOrganizationType int

const (
	CREDENTIAL_ORGANIZATION_TYPE_COOP CredentialOrganizationType = iota
	CREDENTIAL_ORGANIZATION_TYPE_CLUB
	CREDENTIAL_ORGANIZATION_TYPE_SPORT
	CREDENTIAL_ORGANIZATION_TYPE_COHORT
)

type CredentialOrganizationId int

// Right now, the entire contents of this table can be stored in memory since it's not super big,
// but we can transition to a DB table in the future if need be. Just doing this for performance
// reasons and to get this feature out quickly.
type CredentialOrganization struct {
	Id   CredentialOrganizationId   `json:"id"`
	Name string                     `json:"name"`
	Type CredentialOrganizationType `json:"type"`
}

// Ids below must be unique
var coopOrganizations = []idNamePair{
	idNamePair{1, "Jane Street"},
	idNamePair{2, "Google"},
	idNamePair{3, "Facebook"},
	idNamePair{4, "Uber"},
	idNamePair{5, "TD Bank"},
	idNamePair{6, "University of Waterloo"},
}

var clubOrganizations = []idNamePair{
	idNamePair{10001, "Blueprint"},
	idNamePair{10002, "WATonomous"},
	idNamePair{10003, "CS Club"},
}

var sportOrganizations = []idNamePair{
	idNamePair{20001, "Mens' Basketball Team"},
	idNamePair{20002, "Womens' Soccer Team"},
	idNamePair{20003, "Womens' Badminton Team"},
}

var cohortOrganizations = []idNamePair{
	idNamePair{30001, "Software Engineering"},
	idNamePair{30002, "Computer Engineering"},
	idNamePair{30003, "Computer Science"},
	idNamePair{30004, "Science"},
}

func buildOrganizationTypeSlice(
	tpe CredentialOrganizationType,
	orgList []idNamePair,
) []CredentialOrganization {
	orgs := make([]CredentialOrganization, len(orgList))
	for i, pair := range orgList {
		id := CredentialOrganizationId(pair.id)
		orgs[i] = CredentialOrganization{id, pair.name, tpe}
	}
	return orgs
}

// Mapping from CredentialOrganizationType to slice of CredentialOrganizations
func BuildOrganizations() []CredentialOrganization {
	orgs := []CredentialOrganization{}

	orgs = append(orgs, buildOrganizationTypeSlice(
		CREDENTIAL_ORGANIZATION_TYPE_COOP,
		coopOrganizations,
	)...)
	orgs = append(orgs, buildOrganizationTypeSlice(
		CREDENTIAL_ORGANIZATION_TYPE_CLUB,
		clubOrganizations,
	)...)
	orgs = append(orgs, buildOrganizationTypeSlice(
		CREDENTIAL_ORGANIZATION_TYPE_SPORT,
		sportOrganizations,
	)...)
	orgs = append(orgs, buildOrganizationTypeSlice(
		CREDENTIAL_ORGANIZATION_TYPE_COHORT,
		cohortOrganizations,
	)...)

	return orgs
}

// Mapping from organization id to CredentialOrganizationType using slice created by
// BuildOrganizations.
func BuildOrganizationIdIndex() map[CredentialOrganizationId]CredentialOrganization {
	orgs := BuildOrganizations()
	orgMap := map[CredentialOrganizationId]CredentialOrganization{}

	for _, org := range orgs {
		if _, ok := orgMap[org.Id]; ok {
			panic(fmt.Sprintf("Duplicate org id %d\n", org.Id))
		}
		orgMap[org.Id] = org
	}

	return orgMap
}
