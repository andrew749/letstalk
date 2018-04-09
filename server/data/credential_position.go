package data

import (
	"fmt"
)

type CredentialPositionType int

const (
	CREDENTIAL_POSITION_TYPE_COOP CredentialPositionType = iota
	CREDENTIAL_POSITION_TYPE_CLUB
	CREDENTIAL_POSITION_TYPE_SPORT
	CREDENTIAL_POSITION_TYPE_COHORT
)

// Right now, the entire contents of this table can be stored in memory since it's not super big,
// but we can transition to a DB table in the future if need be. Just doing this for performance
// reasons and to get this feature out quickly.
type CredentialPosition struct {
	Id   int                    `json:"id"`
	Name string                 `json:"name"`
	Type CredentialPositionType `json:"type"`
}

type idNamePair struct {
	id   int
	name string
}

// Ids below must be unique
var coopPositions = []idNamePair{
	idNamePair{1, "Software Engineer"},
	idNamePair{2, "Product Manager"},
	idNamePair{3, "Product Designer"},
	idNamePair{4, "Business Analyst"},
	idNamePair{5, "Trader"},
	idNamePair{6, "Researcher"},
}

var clubPositions = []idNamePair{
	idNamePair{10001, "President"},
	idNamePair{10002, "VP Marketing"},
	idNamePair{10003, "Member"},
}

var sportPositions = []idNamePair{
	idNamePair{20001, "Captain"},
	idNamePair{20002, "Point Guard"},
	idNamePair{20003, "Striker"},
}

var cohortPositions = []idNamePair{
	idNamePair{30001, "Student Rep"},
	idNamePair{30002, "Student"},
}

func buildPositionTypeSlice(
	tpe CredentialPositionType,
	orgList []idNamePair,
) []CredentialPosition {
	orgs := make([]CredentialPosition, len(orgList))
	for i, pair := range orgList {
		orgs[i] = CredentialPosition{pair.id, pair.name, tpe}
	}
	return orgs
}

// Mapping from CredentialPositionType to slice of CredentialPositions
func BuildPositions() []CredentialPosition {
	positions := []CredentialPosition{}

	positions = append(positions, buildPositionTypeSlice(
		CREDENTIAL_POSITION_TYPE_COOP,
		coopPositions,
	)...)
	positions = append(positions, buildPositionTypeSlice(
		CREDENTIAL_POSITION_TYPE_CLUB,
		clubPositions,
	)...)
	positions = append(positions, buildPositionTypeSlice(
		CREDENTIAL_POSITION_TYPE_SPORT,
		sportPositions,
	)...)
	positions = append(positions, buildPositionTypeSlice(
		CREDENTIAL_POSITION_TYPE_COHORT,
		cohortPositions,
	)...)

	return positions
}

// Mapping from position id to CredentialPositionType using slice created by BuildPositions.
func BuildInversePositionTypeMap(
	positionTypeMap map[CredentialPositionType][]CredentialPosition,
) map[int]CredentialPositionType {
	positionInverseTypeMap := map[int]CredentialPositionType{}

	for positionType, positions := range positionTypeMap {
		for _, position := range positions {
			if _, ok := positionInverseTypeMap[position.Id]; ok {
				panic(fmt.Sprintf("Duplicate position id %d\n", position.Id))
			}
			positionInverseTypeMap[position.Id] = positionType
		}
	}

	return positionInverseTypeMap
}
