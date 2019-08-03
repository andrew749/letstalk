package recommendations

// Returns recommendation strategy for the most common mentor/mentee matching
// Gives priority to attributs in the following order:
// Program > Gender > Survey answers
func MentorMenteeStrat(
	maxLowerYearsPerUpperYear uint,
	maxUpperYearsPerLowerYear uint,
	youngestUpperYear uint,
) RecommendationStrategy {
	return RecommendationStrategy{
		Score: CombinedPairwiseScore{
			UserScores: []UserScoreWithWeight{},
			// Weights set up in such a way that:
			// - Same program match will always beat different program
			// - If different program, same gender will always beat different gender
			PairwiseScores: []PairwiseScoreWithWeight{
				PairwiseScoreWithWeight{
					PairwiseScore: ProgramPairwiseScore{},
					Weight:        4.0,
				},
				PairwiseScoreWithWeight{
					PairwiseScore: GenderPairwiseScore{},
					Weight:        2.0,
				},
				PairwiseScoreWithWeight{
					PairwiseScore: SurveyPairwiseScore{},
					Weight:        1.0,
				},
			},
		},
		Matcher: GreedyUpperLowerYearMatcher{
			MaxLowerYearsPerUpperYear: maxLowerYearsPerUpperYear,
			MaxUpperYearsPerLowerYear: maxUpperYearsPerLowerYear,
			YoungestUpperYear:         youngestUpperYear,
		},
	}
}
