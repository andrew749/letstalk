package recommendations

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"letstalk/server/core/survey"
	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/test_helpers"
)

func TestSurveyPairwiseScoreRequiredObjects(t *testing.T) {
	assert.Equal(t, []string{"UserSurveys"}, SurveyPairwiseScore{}.RequiredObjects())
}

func TestSurveyPairwiseScoreCalculateSomeOverlap(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)
			responses11 := map[data.SurveyQuestionKey]data.SurveyOptionKey{
				"free_time":  "reading",
				"group_size": "both",
			}
			err = test_helpers.CreateSurveyForUser(db, user1, responses11, survey.Generic_v1.Group, 1)
			assert.NoError(t, err)
			responses12 := map[data.SurveyQuestionKey]data.SurveyOptionKey{
				"interests": "distributed",
			}
			err = test_helpers.CreateSurveyForUser(db, user1, responses12, survey.Se_soc_v1.Group, 1)
			assert.NoError(t, err)

			user2, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)
			responses21 := map[data.SurveyQuestionKey]data.SurveyOptionKey{
				"free_time":  "reading",
				"group_size": "larger",
			}
			err = test_helpers.CreateSurveyForUser(db, user2, responses21, survey.Generic_v1.Group, 1)
			assert.NoError(t, err)
			responses22 := map[data.SurveyQuestionKey]data.SurveyOptionKey{
				"interests": "distributed",
			}
			err = test_helpers.CreateSurveyForUser(db, user2, responses22, survey.Wics_v1.Group, 1)
			assert.NoError(t, err)

			score, err := SurveyPairwiseScore{}.Calculate(user1, user2)
			assert.NoError(t, err)
			assert.Equal(t, Score(0.5), score)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestSurveyPairwiseScoreCalculateNoOverlap(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			user1, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)
			responses11 := map[data.SurveyQuestionKey]data.SurveyOptionKey{
				"free_time": "reading",
			}
			err = test_helpers.CreateSurveyForUser(db, user1, responses11, survey.Generic_v1.Group, 1)
			assert.NoError(t, err)
			responses12 := map[data.SurveyQuestionKey]data.SurveyOptionKey{
				"interests": "distributed",
			}
			err = test_helpers.CreateSurveyForUser(db, user1, responses12, survey.Se_soc_v1.Group, 1)
			assert.NoError(t, err)

			user2, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)
			responses21 := map[data.SurveyQuestionKey]data.SurveyOptionKey{
				"group_size": "larger",
			}
			err = test_helpers.CreateSurveyForUser(db, user2, responses21, survey.Generic_v1.Group, 1)
			assert.NoError(t, err)
			responses22 := map[data.SurveyQuestionKey]data.SurveyOptionKey{
				"interests": "distributed",
			}
			err = test_helpers.CreateSurveyForUser(db, user2, responses22, survey.Wics_v1.Group, 1)
			assert.NoError(t, err)

			score, err := SurveyPairwiseScore{}.Calculate(user1, user2)
			assert.NoError(t, err)
			assert.Equal(t, Score(0.0), score)
		},
	}
	test.RunTestWithDb(thisTest)
}
