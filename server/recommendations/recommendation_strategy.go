package recommendations

import "letstalk/server/algo_helper"

type RecommendationStrategy struct {
	Score   PairwiseScore
	Matcher Matcher
}

func (r *RecommendationStrategy) RequiredObjects() []string {
	options := make([]string, 0)
	options = append(options, r.Matcher.RequiredObjects()...)
	options = append(options, r.Score.RequiredObjects()...)
	return algo_helper.DedupStringList(options)
}
