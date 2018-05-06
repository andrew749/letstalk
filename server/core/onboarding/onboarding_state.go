package onboarding

import (
	"letstalk/server/core/query"
	"letstalk/server/core/errs"

	"github.com/jinzhu/gorm"
	"letstalk/server/core/api"
)

func GetOnboardingInfo(db *gorm.DB, userId int) (*api.OnboardingInfo, errs.Error) {
	userCohort, err := query.GetUserCohort(db, userId)
	onboardingInfo := &api.OnboardingInfo{api.ONBOARDING_DONE, api.USER_TYPE_UNKNOWN, userCohort}
	// TODO: Should probably check what the errors here are.
	if err != nil {
		onboardingInfo.State = api.ONBOARDING_COHORT
		return onboardingInfo, nil
	}

	userVectors, err := query.GetUserVectorsById(db, userId)
	if err != nil {
		return nil, errs.NewDbError(err)
	}

	// TODO: determine user type based on cohort response
	onboardingInfo.UserType = api.USER_TYPE_MENTOR
	if userVectors.Me == nil {
		onboardingInfo.State = api.ONBOARDING_VECTOR_ME
	} else if userVectors.You == nil {
		onboardingInfo.State = api.ONBOARDING_VECTOR_YOU
	}

	return onboardingInfo, nil
}
