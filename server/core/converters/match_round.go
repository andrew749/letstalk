package converters

import (
	"letstalk/server/core/api"
	"letstalk/server/data"
)

func ApiMatchUserFromDataUser(user *data.User) api.MatchUser {
	var cohort *api.CohortV2 = nil
	if user.Cohort != nil && user.Cohort.Cohort != nil {
		cohort = ApiCohortV2FromDataCohort(user.Cohort.Cohort)
	}

	return api.MatchUser{
		User:   ApiUserPersonalInfoFromDataUser(user),
		Email:  user.Email,
		Cohort: cohort,
	}
}

// Assumes the users are preloaded onto the match round, otherwise crashes.
// Doesn't assume cohorts exist for the users.
func apiMatchRoundMatchFromData(matchRound *data.MatchRoundMatch) api.MatchRoundMatch {
	return api.MatchRoundMatch{
		Mentee: ApiMatchUserFromDataUser(matchRound.MenteeUser),
		Mentor: ApiMatchUserFromDataUser(matchRound.MentorUser),
		Score:  matchRound.Score,
	}
}

func ApiMatchRoundFromDataEntities(
	matchRound *data.MatchRound,
	state api.MatchRoundState,
) api.MatchRound {
	apiMatches := make([]api.MatchRoundMatch, 0, len(matchRound.Matches))
	for _, match := range matchRound.Matches {
		apiMatches = append(apiMatches, apiMatchRoundMatchFromData(&match))
	}

	return api.MatchRound{
		MatchRoundId: matchRound.Id,
		Name:         matchRound.Name,
		Matches:      apiMatches,
		State:        state,
	}
}
