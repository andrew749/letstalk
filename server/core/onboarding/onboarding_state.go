package onboarding

import (
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"letstalk/server/core/api"

	"github.com/jinzhu/gorm"
)

func GetOnboardingInfo(db *gorm.DB, userId int) (*api.OnboardingInfo, errs.IError) {
	var user data.User

	err := db.Where(&data.User{UserId: userId}).Preload(
		"Cohort.Cohort",
	).Preload("AdditionalData").First(
		&user,
	).Error

	if err != nil || user.Cohort == nil || user.AdditionalData == nil {
		return &api.OnboardingInfo{State: api.ONBOARDING_COHORT}, nil
	}
	onboardingInfo := &api.OnboardingInfo{
		api.ONBOARDING_DONE,
		api.USER_TYPE_UNKNOWN,
		user.Cohort.Cohort,
		user.AdditionalData,
	}
	return onboardingInfo, nil
}
