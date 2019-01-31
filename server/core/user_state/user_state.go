// This package maintains a state machine for a user's account and is used during onboariding.
// It tells us what information, if any, we need to collect from the user before we let them into
// the app.
package user_state

import (
	"letstalk/server/core/api"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/core/survey"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func GetUserState(db *gorm.DB, userId data.TUserID) (*api.UserState, errs.Error) {
	state := api.ACCOUNT_CREATED
	user, err := query.GetUserById(db, userId)
	if err != nil {
		return nil, err
	}
	if !user.IsEmailVerified {
		return &state, nil
	}
	state = api.ACCOUNT_EMAIL_VERIFIED

	cohort, err := query.GetUserCohort(db, userId)
	if err != nil {
		return nil, err
	}
	additionalData, err := query.GetUserAdditionalData(db, userId)
	if err != nil {
		return nil, err
	}
	if cohort == nil || additionalData == nil {
		return &state, nil
	}
	state = api.ACCOUNT_HAS_BASIC_INFO

	survey, err := query.GetSurvey(db, userId, survey.Generic_v1.Group)
	if err != nil {
		return nil, err
	}
	if survey.Responses == nil {
		return &state, nil
	}
	state = api.ACCOUNT_SETUP
	return &state, nil
}
