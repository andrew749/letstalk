package recommendations

import (
	"letstalk/server/algo_helper"
	"letstalk/server/data"
)

// A score like w_1 * x_1 + w_2 * x_2 + ... + w_n * x_n for a pair of users.
// Users scores are applied to both users and summed.
// Pairwise scores are applied once to both users.
type CombinedPairwiseScore struct {
	UserScores     []UserScoreWithWeight
	PairwiseScores []PairwiseScoreWithWeight
}

func (s CombinedPairwiseScore) RequiredObjects() []string {
	options := make([]string, 0)
	for _, userScore := range s.UserScores {
		options = append(options, userScore.UserScore.RequiredObjects()...)
	}
	for _, pairwiseScore := range s.PairwiseScores {
		options = append(options, pairwiseScore.PairwiseScore.RequiredObjects()...)
	}
	return algo_helper.DedupStringList(options)
}

func (s CombinedPairwiseScore) Calculate(
	userOne *data.User,
	userTwo *data.User,
) (Score, error) {
	var score Score = 0.0
	for _, userScore := range s.UserScores {
		value, err := userScore.UserScore.Calculate(userOne)
		if err != nil {
			return 0.0, err
		}
		score = score.Add(value.Weighted(userScore.Weight))
		value, err = userScore.UserScore.Calculate(userTwo)
		if err != nil {
			return 0.0, err
		}
		score = score.Add(value.Weighted(userScore.Weight))
	}
	for _, pairwiseScore := range s.PairwiseScores {
		value, err := pairwiseScore.PairwiseScore.Calculate(userOne, userTwo)
		if err != nil {
			return 0.0, err
		}
		score = score.Add(value.Weighted(pairwiseScore.Weight))
	}
	return score, nil
}
