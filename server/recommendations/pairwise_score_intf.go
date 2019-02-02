package recommendations

import "letstalk/server/data"

// Calculates a score for two users
// RequiredObjects returns a list of user objects required to calculate the score
// (e.g. cohort, surveys).
// Calculate actually calculates the pairwise score for two users. Generally, only error if some
// unexpected value occurs, otherwise, use sane default values.
type PairwiseScore interface {
	RequiredObjects() []string
	Calculate(userOne *data.User, userTwo *data.User) (Score, error)
}

type PairwiseScoreWithWeight struct {
	PairwiseScore PairwiseScore
	Weight        Weight
}
