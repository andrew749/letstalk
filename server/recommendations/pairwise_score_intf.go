package recommendations

import "letstalk/server/data"

// Calculates a score for two users
// RequiredObjects returns a list of objects required to calculate the score (e.g. cohort, surveys)
type PairwiseScore interface {
	RequiredObjects() []string
	Calculate(userOne *data.User, userTwo *data.User) (float32, error)
}
