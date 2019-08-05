// This package maintains a state machine for a user's account and is used during onboariding.
// It tells us what information, if any, we need to collect from the user before we let them into
// the app.
package user_state

import (
	"fmt"
	"letstalk/server/core/api"
	"letstalk/server/core/errs"
	"letstalk/server/core/survey"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func GetUserState(db *gorm.DB, userId data.TUserID) (*api.UserState, errs.Error) {
	userWithStates, err := BulkGetUsersWithStates(db, []data.TUserID{userId})
	if err != nil {
		return nil, err
	}
	if len(userWithStates) == 0 {
		return nil, errs.NewRequestError(fmt.Sprintf("User with id %d not found", userId))
	}
	return &userWithStates[0].State, err
}

type UserWithState struct {
	User  data.User
	State api.UserState
}

var REQUIRED_ONBOARDING_SURVEY_GROUP = survey.Generic_v1.Group

func containsOnboardingSurvey(surveys []data.UserSurvey) bool {
	if surveys != nil {
		for _, survey := range surveys {
			if survey.Group == REQUIRED_ONBOARDING_SURVEY_GROUP {
				return true
			}
		}
	}
	return false
}

func userStateFromAssociations(user *data.User) api.UserState {
	if !user.IsEmailVerified {
		return api.ACCOUNT_CREATED
	} else if user.Cohort == nil || user.AdditionalData == nil {
		return api.ACCOUNT_EMAIL_VERIFIED
	} else if !containsOnboardingSurvey(user.UserSurveys) {
		return api.ACCOUNT_HAS_BASIC_INFO
	} else {
		return api.ACCOUNT_SETUP
	}
}

func BulkGetUsersWithStates(db *gorm.DB, userIds []data.TUserID) ([]UserWithState, errs.Error) {
	var users []data.User
	err := db.Where("user_id IN (?)", userIds).Preload("Cohort").
		Preload("AdditionalData").Preload("UserSurveys").Find(&users).Error
	if err != nil {
		return nil, errs.NewDbError(err)
	}

	userWithStates := make([]UserWithState, 0, len(users))
	for _, user := range users {
		userWithStates = append(userWithStates, UserWithState{
			User:  user,
			State: userStateFromAssociations(&user),
		})
	}
	return userWithStates, nil
}
