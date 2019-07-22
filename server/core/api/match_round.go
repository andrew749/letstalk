package api

import "letstalk/server/data"

type MatchRoundParameters struct {
	MaxLowerYearsPerUpperYear uint `json:"maxLowerYearsPerUpperYear" binding:"required"`
	MaxUpperYearsPerLowerYear uint `json:"maxUpperYearsPerLowerYear" binding:"required"`
	YoungestUpperGradYear     uint `json:"youngestUpperGradYear" binding:"required"`
}

type CreateMatchRoundRequest struct {
	// Parameters used to match
	Parameters MatchRoundParameters `json:"parameters" binding:"required"`

	// TODO(match-api): Figure out if this is the right ID after Andrew's changes
	// Group id associated with the match round
	GroupId data.TGroupID `json:"groupId" binding:"required"`

	// Users selected to take part in match round by admin
	// This endpoint checks that the users are in the given group
	UserIds []data.TUserID `json:"userIds" binding:"required"`
}

// Doesn't enforce that cohort exists so that we are more liberal with which type of users are
// accepted for a match.
type MatchUser struct {
	User   UserPersonalInfo `json:"user" binding:"required"`
	Cohort *CohortV2        `json:"cohort"`
}

type MatchRoundMatch struct {
	Mentee MatchUser `json:"mentee" binding:"required"`
	Mentor MatchUser `json:"mentor" binding:"required"`
	Score  float32   `json:"score" binding:"required"`
}

type MatchRound struct {
	MatchRoundId data.TMatchRoundID `json:"matchRoundId" binding:"required"`
	Name         string             `json:"name" binding:"required"`
	Matches      []MatchRoundMatch  `json:"matches" binding:"required"`
}
