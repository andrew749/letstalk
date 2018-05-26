package onboarding

import (
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"letstalk/server/core/api"
)

func GetOnboardingInfo(db *gorm.DB, userId int) (*api.OnboardingInfo, errs.Error) {
	var user data.User

	err := db.Where(&data.User{UserId: userId}).Preload("Cohort.Cohort").Preload("Preference").First(
		&user,
	).Error

	if err != nil || user.Cohort == nil || user.Preference == nil {
		return &api.OnboardingInfo{State: api.ONBOARDING_COHORT}, nil
	}
	onboardingInfo := &api.OnboardingInfo{
		api.ONBOARDING_DONE,
		api.USER_TYPE_UNKNOWN,
		user.Cohort.Cohort,
		user.Preference,
	}
	return onboardingInfo, nil
}
