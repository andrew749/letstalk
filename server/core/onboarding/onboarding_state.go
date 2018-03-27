package onboarding

import (
	"letstalk/server/core/api"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/mijia/modelq/gmq"
)

/**
 * The onboarding state specifies which onboarding information we want to get next for the user.
 */
type OnboardingState string

const (
	ONBOARDING_COHORT     OnboardingState = "onboarding_cohort"     // get cohort info
	ONBOARDING_VECTOR_ME  OnboardingState = "onboarding_vector_me"  // get my personality vector
	ONBOARDING_VECTOR_YOU OnboardingState = "onboarding_vector_you" // get personality vector for others
	ONBOARDING_DONE       OnboardingState = "onboarding_done"       // finished
)

type OnboardingInfo struct {
	State      OnboardingState `json:"state" binding:"required"`
	UserType   api.UserType    `json:"userType" binding:"required"`
	UserCohort *data.Cohort    `json:"userCohort"`
}

func GetOnboardingInfo(db *gmq.Db, userId int) (*OnboardingInfo, errs.Error) {
	userCohort, err := api.GetUserCohort(db, userId)
	onboardingInfo := &OnboardingInfo{ONBOARDING_DONE, api.USER_TYPE_UNKNOWN, userCohort}
	// TODO: Should probably check what the errors here are.
	if err != nil {
		onboardingInfo.State = ONBOARDING_COHORT
		return onboardingInfo, nil
	}

	userVectors, err := api.GetUserVectorsById(db, userId)
	if err != nil {
		return nil, errs.NewDbError(err)
	}

	// TODO: determine user type based on cohort response
	onboardingInfo.UserType = api.USER_TYPE_MENTOR
	if userVectors.Me == nil {
		onboardingInfo.State = ONBOARDING_VECTOR_ME
	} else if userVectors.You == nil {
		onboardingInfo.State = ONBOARDING_VECTOR_YOU
	}

	return onboardingInfo, nil
}
