package recommendations

import "github.com/jinzhu/gorm"

func Recommend(
	db *gorm.DB,
	fetcherOptions UserFetcherOptions,
	strat RecommendationStrategy,
) ([]UserMatch, error) {
	users, err := FetchUsers(db, fetcherOptions, strat.RequiredObjects())
	if err != nil {
		return nil, err
	}

	return strat.Matcher.Match(users, strat.Score)
}
