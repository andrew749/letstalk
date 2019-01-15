package recommendations

type RecommendationStrategy struct {
	UserFetcherOptions UserFetcherOptions
	UserScores         []UserScoreWithWeight
	PairwiseScores     []PairwiseScoreWithWeight
	Matcher            Matcher
}

func dedupStringList(ss []string) []string {
	smap := make(map[string]interface{})
	for _, s := range ss {
		smap[s] = nil
	}
	ssNew := make([]string, 0, len(smap))
	for s := range smap {
		ssNew = append(ssNew, s)
	}
	return ssNew
}

func (r *RecommendationStrategy) RequiredObjects() []string {
	options := make([]string, 0)
	options = append(options, r.Matcher.RequiredObjects()...)
	for _, userScore := range r.UserScores {
		options = append(options, userScore.UserScore.RequiredObjects()...)
	}
	for _, pairwiseScore := range r.PairwiseScores {
		options = append(options, pairwiseScore.PairwiseScore.RequiredObjects()...)
	}
	return dedupStringList(options)
}
