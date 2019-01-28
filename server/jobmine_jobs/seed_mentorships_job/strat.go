package seed_mentorships_job

import (
	"time"

	"letstalk/server/data"
	"letstalk/server/recommendations"
)

// Gets RecommendationStrategy
func getRecommendationStrategy(
	maxLowerYearsPerUpperYear uint,
	maxUpperYearsPerLowerYear uint,
	youngestUpperYear uint,
) recommendations.RecommendationStrategy {
	return recommendations.RecommendationStrategy{
		Score: recommendations.CombinedPairwiseScore{
			UserScores: []recommendations.UserScoreWithWeight{},
			// Weights set up in such a way that:
			// - Same program match will always beat different program
			// - If different program, same gender will always beat different gender
			PairwiseScores: []recommendations.PairwiseScoreWithWeight{
				recommendations.PairwiseScoreWithWeight{
					PairwiseScore: recommendations.ProgramPairwiseScore{},
					Weight:        4.0,
				},
				recommendations.PairwiseScoreWithWeight{
					PairwiseScore: recommendations.GenderPairwiseScore{},
					Weight:        2.0,
				},
				recommendations.PairwiseScoreWithWeight{
					PairwiseScore: recommendations.SurveyPairwiseScore{},
					Weight:        1.0,
				},
			},
		},
		Matcher: recommendations.GreedyUpperLowerYearMatcher{
			MaxLowerYears:     maxLowerYearsPerUpperYear,
			MaxUpperYears:     maxUpperYearsPerLowerYear,
			YoungestUpperYear: youngestUpperYear,
		},
	}
}

// Same as above except downweights users that were made before the start of the term
func getRecommendationStrategyWithOlderDownrank(
	maxLowerYearsPerUpperYear uint,
	maxUpperYearsPerLowerYear uint,
	youngestUpperYear uint,
	termStartTime time.Time,
	blacklistUserIds []data.TUserID,
) recommendations.RecommendationStrategy {
	strat := getRecommendationStrategy(
		maxLowerYearsPerUpperYear, maxUpperYearsPerLowerYear, youngestUpperYear)
	combinedScore := strat.Score.(recommendations.CombinedPairwiseScore)

	blacklistUserIdSet := make(map[data.TUserID]interface{})
	for _, userId := range blacklistUserIds {
		blacklistUserIdSet[userId] = nil
	}

	combinedScore.UserScores = append(combinedScore.UserScores, recommendations.UserScoreWithWeight{
		UserScore: recommendations.UserScoreOlder{
			Before:           termStartTime,
			BlacklistUserIds: blacklistUserIdSet,
		},
		// Set it to this, since we probably want to prefer old but same program over new but different
		// program.
		Weight: -3.0,
	})
	strat.Score = combinedScore
	return strat
}
