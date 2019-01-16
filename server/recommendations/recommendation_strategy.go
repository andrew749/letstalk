package recommendations

type RecommendationStrategy struct {
	Score   PairwiseScore
	Matcher Matcher
}

func (r *RecommendationStrategy) RequiredObjects() []string {
	options := make([]string, 0)
	options = append(options, r.Matcher.RequiredObjects()...)
	options = append(options, r.Score.RequiredObjects()...)
	return dedupStringList(options)
}
