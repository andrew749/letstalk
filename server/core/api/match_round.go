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
