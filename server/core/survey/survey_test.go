package survey

import (
	"testing"

	"letstalk/server/core/ctx"
	"letstalk/server/core/test"
	"letstalk/server/core/user"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestSaveSurvey(t *testing.T) {
	tests := []test.Test{
		{
			Test: func(db *gorm.DB) {
				c := ctx.NewContext(nil, db, nil, nil, nil)
				userOne := user.CreateUserForTest(t, c.Db)
				getResponses, err := getSurveyResponses(db, userOne.UserId, Generic_v1.Group)
				// Assert no responses initially.
				assert.Nil(t, getResponses)
				assert.NoError(t, err)
				// Save some new responses.
				responses := data.SurveyResponses{}
				for _, question := range Generic_v1.Questions {
					responses[question.Key] = question.Options[0].Key
				}
				err = saveSurveyResponses(db, userOne.UserId, Generic_v1.Group, Generic_v1.Version, responses)
				assert.NoError(t, err)
				getResponses, err = getSurveyResponses(db, userOne.UserId, Generic_v1.Group)
				// Assert no responses initially.
				assert.NoError(t, err)
				assert.NotNil(t, getResponses)
				for _, question := range Generic_v1.Questions {
					assert.Equal(t, question.Options[0].Key, (*getResponses)[question.Key])
				}
			},
			TestName: "Test set survey responses",
		},
	}
	test.RunTestsWithDb(tests)
}
