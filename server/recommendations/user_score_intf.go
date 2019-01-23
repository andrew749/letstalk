package recommendations

import "letstalk/server/data"

type Score float32
type Weight float32

func (s1 Score) Add(s2 Score) Score {
	return s1 + s2
}

func (s Score) Weighted(w Weight) Score {
	return Score(float32(s) * float32(w))
}

// Calculates a score for a single user
// RequiredObjects returns a list of user objects required to calculate the score
// (e.g. cohort, surveys).
// Calculate actually calculates the score for a user. Generally, only error if some unexpected
// value occurs, otherwise, use sane default values.
type UserScore interface {
	RequiredObjects() []string
	Calculate(user *data.User) (Score, error)
}

type UserScoreWithWeight struct {
	UserScore UserScore
	Weight    Weight
}
