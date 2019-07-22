package converters

import (
	"letstalk/server/core/api"
	"letstalk/server/data"
)

func apiMatchUserFromDataUser(user *data.User) api.MatchUser {
	var cohort *api.CohortV2 = nil
	if user.Cohort != nil && user.Cohort.Cohort != nil {
		cohort = ApiCohortV2FromDataCohort(user.Cohort.Cohort)
	}

	return api.MatchUser{
		User:   ApiUserPersonalInfoFromDataUser(user),
		Cohort: cohort,
	}
}

// Assumes the users are preloaded onto the match round, otherwise crashes.
// Doesn't assume cohorts exist for the users.
func apiMatchRoundMatchFromData(matchRound *data.MatchRoundMatch) api.MatchRoundMatch {
	return api.MatchRoundMatch{
		Mentee: apiMatchUserFromDataUser(matchRound.MenteeUser),
		Mentor: apiMatchUserFromDataUser(matchRound.MentorUser),
		Score:  matchRound.Score,
	}
}

func ApiMatchRoundFromDataEntities(
	matchRound *data.MatchRound,
	matches []data.MatchRoundMatch,
) api.MatchRound {
	apiMatches := make([]api.MatchRoundMatch, 0, len(matches))
	for _, match := range matches {
		apiMatches = append(apiMatches, apiMatchRoundMatchFromData(&match))
	}

	return api.MatchRound{
		MatchRoundId: matchRound.Id,
		Name:         matchRound.Name,
		Matches:      apiMatches,
	}
}
