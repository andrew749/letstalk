package survey

import (
	"testing"

	"letstalk/server/core/ctx"
	"letstalk/server/core/test"
	"letstalk/server/core/user"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"letstalk/server/core/api"
	"net/http"
)

func TestSaveSurvey(t *testing.T) {
	tests := []test.Test{
		{
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				userOne := user.CreateUserForTest(t, c.Db)
				survey, err := getSurvey(db, userOne.UserId, Generic_v1.Group)
				// Assert no responses initially.
				assert.NotNil(t, survey)
				assert.Nil(t, survey.Responses)
				assert.NoError(t, err)
				// Save some new responses.
				responses := data.SurveyResponses{}
				for _, question := range Generic_v1.Questions {
					responses[question.Key] = question.Options[0].Key
				}
				err = saveSurveyResponses(db, userOne.UserId, Generic_v1.Group, Generic_v1.Version, responses)
				assert.NoError(t, err)
				survey, err = getSurvey(db, userOne.UserId, Generic_v1.Group)
				// Assert responses are fetched.
				assert.NoError(t, err)
				assert.NotNil(t, survey.Responses)
				for _, question := range Generic_v1.Questions {
					assert.Equal(t, question.Options[0].Key, (*survey.Responses)[question.Key])
				}
			},
			TestName: "Test set survey responses",
		},
	}
	test.RunTestsWithDb(tests)
}

func TestSurveyGroups(t *testing.T) {
	testSurvey := api.Survey{
		Group: "testSurveyGroup",
		Version: 1,
		Questions: []api.SurveyQuestion{{
			Key: "q1",
			Prompt: "test survey prompt",
			Options: []api.SurveyOption{{
				Key: "q1a",
				Text: "test survey option",
			}},
		}},
	}
	tests := []test.Test{
		{
			TestName: "Test survey groups",
			Test: func(db *gorm.DB) {
				userId := data.TUserID(1)
				genericResponse, err := getSurveyResponses(db, userId, Generic_v1.Group)
				testResponse, err := getSurveyResponses(db, userId, testSurvey.Group)
				// Assert no responses initially.
				assert.Nil(t, genericResponse)
				assert.Nil(t, testResponse)
				// Save some new responses.
				responses := data.SurveyResponses{Generic_v1.Questions[0].Key: Generic_v1.Questions[0].Options[0].Key}
				err = saveSurveyResponses(db, userId, Generic_v1.Group, Generic_v1.Version, responses)
				assert.NoError(t, err)
				responses = data.SurveyResponses{testSurvey.Questions[0].Key: testSurvey.Questions[0].Options[0].Key}
				err = saveSurveyResponses(db, userId, testSurvey.Group, testSurvey.Version, responses)
				assert.NoError(t, err)
				genericResponse, err = getSurveyResponses(db, userId, Generic_v1.Group)
				assert.NoError(t, err)
				testResponse, err = getSurveyResponses(db, userId, testSurvey.Group)
				assert.NoError(t, err)
				assert.Equal(t, Generic_v1.Questions[0].Options[0].Key, (*genericResponse)[Generic_v1.Questions[0].Key])
				assert.Equal(t, testSurvey.Questions[0].Options[0].Key, (*testResponse)[testSurvey.Questions[0].Key])
			},
		},
		{
			TestName: "Test survey get bad group",
			Test: func(db *gorm.DB) {
				userId := data.TUserID(1)
				survey, err := getSurvey(db, userId, data.SurveyGroup("bad_survey_group"))
				assert.Nil(t, survey)
				assert.Equal(t, http.StatusNotFound, err.GetHTTPCode())
			},
		},
	}
	test.RunTestsWithDb(tests)
}
