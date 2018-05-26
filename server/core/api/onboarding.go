package api

import "letstalk/server/data"

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
	State          OnboardingState      `json:"state" binding:"required"`
	UserType       UserType             `json:"userType" binding:"required"`
	UserCohort     *data.Cohort         `json:"userCohort"`
	UserPreference *data.UserPreference `json:"userPreference"`
}

type OnboardingStatus struct {
	State    OnboardingState `json:"state" binding:"required"`
	UserType UserType        `json:"userType" binding:"required"`
}

type UpdateCohortRequest struct {
	CohortId             int `json:"cohortId" binding:"required"`
	MentorshipPreference int `json:"mentorshipPreference" binding:"required"`
}

type OnboardingUpdateResponse struct {
	Message          string            `json:"message" binding:"required"`
	OnboardingStatus *OnboardingStatus `json:"onboardingStatus" binding:"required"`
}
