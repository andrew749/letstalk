package api

import (
	"letstalk/server/data"
)

type ValidCredentialPair struct {
	PositionType     data.CredentialPositionType     `json:"positionType"`
	OrganizationType data.CredentialOrganizationType `json:"organizationType"`
}

type CredentialOptions struct {
	ValidPairs    []ValidCredentialPair         `json:"validPairs"`
	Organizations []data.CredentialOrganization `json:"organizations"`
	Positions     []data.CredentialPosition     `json:"positions"`
}

var validPairs = []ValidCredentialPair{
	ValidCredentialPair{
		data.CREDENTIAL_POSITION_TYPE_COOP,
		data.CREDENTIAL_ORGANIZATION_TYPE_COOP,
	},
	ValidCredentialPair{
		data.CREDENTIAL_POSITION_TYPE_CLUB,
		data.CREDENTIAL_ORGANIZATION_TYPE_CLUB,
	},
	ValidCredentialPair{
		data.CREDENTIAL_POSITION_TYPE_SPORT,
		data.CREDENTIAL_ORGANIZATION_TYPE_SPORT,
	},
	ValidCredentialPair{
		data.CREDENTIAL_POSITION_TYPE_COHORT,
		data.CREDENTIAL_ORGANIZATION_TYPE_COHORT,
	},
}

// Returns a struct contain all info required to generate all possible credential options, where
// a credential consists of a position and an organization.
func GetCredentialOptions() CredentialOptions {
	// TODO: Could cache the results of BuildOrganizations and BuildPositions
	return CredentialOptions{
		ValidPairs:    validPairs,
		Organizations: data.BuildOrganizations(),
		Positions:     data.BuildPositions(),
	}
}
