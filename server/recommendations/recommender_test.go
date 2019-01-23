package recommendations

import (
	"sort"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/test_helpers"
)

func getUpperLowerYearStrat(
	maxLowerYears uint,
	maxUpperYears uint,
	youngestUpperYear uint,
) RecommendationStrategy {
	return RecommendationStrategy{
		Score: CombinedPairwiseScore{
			UserScores: []UserScoreWithWeight{
				UserScoreWithWeight{
					UserScore: UserScoreOlder{
						Before:           time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
						BlacklistUserIds: make(map[data.TUserID]interface{}),
					},
					Weight: -1.0,
				},
			},
			PairwiseScores: []PairwiseScoreWithWeight{
				PairwiseScoreWithWeight{
					PairwiseScore: ProgramPairwiseScore{},
					Weight:        2.0,
				},
			},
		},
		Matcher: GreedyUpperLowerYearMatcher{
			MaxLowerYears:     maxLowerYears,
			MaxUpperYears:     maxUpperYears,
			YoungestUpperYear: youngestUpperYear,
		},
	}
}

func seedUpperLowerYearUsers(t *testing.T, db *gorm.DB) []data.User {
	// Setup
	// New user, in arts, upper year
	user1, err := test_helpers.CreateTestUser(db, 1)
	assert.NoError(t, err)
	user1.CreatedAt = time.Date(2019, time.January, 2, 0, 0, 0, 0, time.UTC)
	err = db.Save(user1).Error
	assert.NoError(t, err)
	err = test_helpers.CreateCohortForUser(db, user1, "ARTS", "Arts", 2021, true, nil, nil)
	assert.NoError(t, err)

	// New user, in math, upper year
	user2, err := test_helpers.CreateTestUser(db, 2)
	assert.NoError(t, err)
	user2.CreatedAt = time.Date(2019, time.January, 2, 0, 0, 0, 0, time.UTC)
	err = db.Save(user2).Error
	assert.NoError(t, err)
	err = test_helpers.CreateCohortForUser(db, user2, "MATH", "Math", 2020, true, nil, nil)
	assert.NoError(t, err)

	// Old user, in arts, upper year
	user3, err := test_helpers.CreateTestUser(db, 3)
	assert.NoError(t, err)
	user3.CreatedAt = time.Date(2018, time.December, 19, 0, 0, 0, 0, time.UTC)
	err = db.Save(user3).Error
	assert.NoError(t, err)
	err = test_helpers.CreateCohortForUser(db, user3, "ARTS", "Arts", 2019, true, nil, nil)
	assert.NoError(t, err)

	// Old user, in math, upper year
	user4, err := test_helpers.CreateTestUser(db, 4)
	assert.NoError(t, err)
	user4.CreatedAt = time.Date(2018, time.December, 19, 0, 0, 0, 0, time.UTC)
	err = db.Save(user4).Error
	assert.NoError(t, err)
	err = test_helpers.CreateCohortForUser(db, user4, "MATH", "Math", 2020, true, nil, nil)
	assert.NoError(t, err)

	// New user, in arts, lower year
	user5, err := test_helpers.CreateTestUser(db, 5)
	assert.NoError(t, err)
	user5.CreatedAt = time.Date(2019, time.January, 2, 0, 0, 0, 0, time.UTC)
	err = db.Save(user5).Error
	assert.NoError(t, err)
	err = test_helpers.CreateCohortForUser(db, user5, "ARTS", "Arts", 2022, true, nil, nil)
	assert.NoError(t, err)

	// New user, in math, lower year
	user6, err := test_helpers.CreateTestUser(db, 6)
	assert.NoError(t, err)
	user6.CreatedAt = time.Date(2019, time.January, 2, 0, 0, 0, 0, time.UTC)
	err = db.Save(user6).Error
	assert.NoError(t, err)
	err = test_helpers.CreateCohortForUser(db, user6, "MATH", "Math", 2023, true, nil, nil)
	assert.NoError(t, err)

	// Old user, in arts, lower year
	user7, err := test_helpers.CreateTestUser(db, 7)
	assert.NoError(t, err)
	user7.CreatedAt = time.Date(2018, time.December, 19, 0, 0, 0, 0, time.UTC)
	err = db.Save(user7).Error
	assert.NoError(t, err)
	err = test_helpers.CreateCohortForUser(db, user7, "ARTS", "Arts", 2022, true, nil, nil)
	assert.NoError(t, err)

	// Old user, in math, lower year
	user8, err := test_helpers.CreateTestUser(db, 8)
	assert.NoError(t, err)
	user8.CreatedAt = time.Date(2018, time.December, 19, 0, 0, 0, 0, time.UTC)
	err = db.Save(user8).Error
	assert.NoError(t, err)
	err = test_helpers.CreateCohortForUser(db, user8, "MATH", "Math", 2023, true, nil, nil)
	assert.NoError(t, err)

	return []data.User{
		*user1, *user2, *user3, *user4, *user5, *user6, *user7, *user8,
	}
}

type byUserOneId []UserMatch

func (a byUserOneId) Len() int {
	return len(a)
}

func (a byUserOneId) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byUserOneId) Less(i, j int) bool {
	// Sort by decreasing
	if a[i].UserOneId == a[j].UserOneId {
		return a[i].UserTwoId < a[j].UserTwoId
	} else {
		return a[i].UserOneId < a[j].UserOneId
	}
}

func TestGreedyUpperLowerYearRecommendationUpperYearBounded(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			users := seedUpperLowerYearUsers(t, db)
			// Actual test
			opts := UserFetcherOptions{}
			strat := getUpperLowerYearStrat(1, 100, 2021)
			res, err := Recommend(db, opts, strat)
			assert.NoError(t, err)

			sort.Sort(byUserOneId(res))
			assert.Equal(t, 4, len(res))
			assert.Equal(t, UserMatch{users[4].UserId, users[0].UserId, 2.0}, res[0])
			assert.Equal(t, UserMatch{users[5].UserId, users[1].UserId, 2.0}, res[1])
			assert.Equal(t, UserMatch{users[6].UserId, users[2].UserId, 0.0}, res[2])
			assert.Equal(t, UserMatch{users[7].UserId, users[3].UserId, 0.0}, res[3])
		},
		TestName: "Test GreedyUpperLowerYearMatcher recommendations bounded by upper years",
	}
	test.RunTestWithDb(thisTest)
}

func TestGreedyUpperLowerYearRecommendationLowerYearBounded(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			users := seedUpperLowerYearUsers(t, db)
			// Actual test
			opts := UserFetcherOptions{}
			strat := getUpperLowerYearStrat(100, 1, 2021)
			res, err := Recommend(db, opts, strat)
			assert.NoError(t, err)

			sort.Sort(byUserOneId(res))
			assert.Equal(t, 4, len(res))
			assert.Equal(t, UserMatch{users[4].UserId, users[0].UserId, 2.0}, res[0])
			assert.Equal(t, UserMatch{users[5].UserId, users[1].UserId, 2.0}, res[1])
			assert.Equal(t, UserMatch{users[6].UserId, users[0].UserId, 1.0}, res[2])
			assert.Equal(t, UserMatch{users[7].UserId, users[1].UserId, 1.0}, res[3])
		},
		TestName: "Test GreedyUpperLowerYearMatcher recommendations bounded by lower years",
	}
	test.RunTestWithDb(thisTest)
}

func TestGreedyUpperLowerYearRecommendationBoundedByUpperYearsButLess(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			users := seedUpperLowerYearUsers(t, db)
			// Actual test
			opts := UserFetcherOptions{}
			strat := getUpperLowerYearStrat(2, 100, 2021)
			res, err := Recommend(db, opts, strat)
			assert.NoError(t, err)

			sort.Sort(byUserOneId(res))
			assert.Equal(t, 8, len(res))
			assert.Equal(t, UserMatch{users[4].UserId, users[0].UserId, 2.0}, res[0])
			assert.Equal(t, UserMatch{users[4].UserId, users[2].UserId, 1.0}, res[1])
			assert.Equal(t, UserMatch{users[5].UserId, users[1].UserId, 2.0}, res[2])
			assert.Equal(t, UserMatch{users[5].UserId, users[3].UserId, 1.0}, res[3])
			assert.Equal(t, UserMatch{users[6].UserId, users[0].UserId, 1.0}, res[4])
			assert.Equal(t, UserMatch{users[6].UserId, users[2].UserId, 0.0}, res[5])
			assert.Equal(t, UserMatch{users[7].UserId, users[1].UserId, 1.0}, res[6])
			assert.Equal(t, UserMatch{users[7].UserId, users[3].UserId, 0.0}, res[7])
		},
		TestName: "Test GreedyUpperLowerYearMatcher recommendations bounded by upper years but less",
	}
	test.RunTestWithDb(thisTest)
}
