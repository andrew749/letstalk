package recommendations

import "letstalk/server/data"

// Calculates a score for a single user
// RequiredObjects returns a list of user objects required to calculate the score
// (e.g. cohort, surveys).
// Calculate actually calculates the score for a user. Generally, only error if some unexpected
// value occurs, otherwise, use sane default values.
type UserScore interface {
	RequiredObjects() []string
	Calculate(user *data.User) (float32, error)
}

type UserScoreWithWeight struct {
	UserScore UserScore
	Weight    float32
}
