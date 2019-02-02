package recommendations

import "letstalk/server/data"

// Creates a bunch of user matches
// RequiredObjects returns a list of user objects required to do matching (e.g. cohort, surveys).
// Match calculates a bunch of
type Matcher interface {
	RequiredObjects() []string
	Match(users []data.User, score PairwiseScore) ([]UserMatch, error)
}
