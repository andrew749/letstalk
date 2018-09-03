package onboarding

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	"letstalk/server/data"
)

// i.e. fetch a onboarding type and the possible options

/**
*
* Sample post request:
 {
 	"cohortId": UserVectorType
 }
*/

// Update a user with new information for their school
// try to match this data to an existing sequence.
func UpdateUserCohort(c *ctx.Context) errs.Error {
	var req api.UpdateCohortRequest

	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError("%s", err.Error())
	}

	if err := query.UpdateUserCohortAndAdditionalInfo(
		c.Db,
		c.Es,
		c.SessionData.UserId,
		req.CohortId,
		req.MentorshipPreference,
		req.Bio,
		req.Hometown,
	); err != nil {
		return err
	}

	onboardingInfo, err := GetOnboardingInfo(c.Db, c.SessionData.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}

	onboardingStatus := &api.OnboardingStatus{
		onboardingInfo.State,
		onboardingInfo.UserType,
	}
	successMessage := "Successfully added cohort to user."
	c.Result = api.OnboardingUpdateResponse{successMessage, onboardingStatus}
	return nil
}

type UserVectorType string

const (
	Sociable     UserVectorType = "Sociable"
	Hard_Working UserVectorType = "Hard Working"
	Ambitious    UserVectorType = "Ambitious"
	Energetic    UserVectorType = "Energetic"
	Carefree     UserVectorType = "Carefree"
	Confident    UserVectorType = "Confident"
)

type UpdateUserVectorRequest struct {
	PreferenceType api.UserVectorPreferenceType `json:"preferenceType" binding:"exists"`
	Sociable       int                          `json:"sociable" binding:"exists"`
	Hard_Working   int                          `json:"hardWorking" binding:"exists"`
	Ambitious      int                          `json:"ambitious" binding:"exists"`
	Energetic      int                          `json:"energetic" binding:"exists"`
	Carefree       int                          `json:"carefree" binding:"exists"`
	Confident      int                          `json:"confident" binding:"exists"`
}

/**
 * Update the user vector.
 */
func UserVectorUpdateController(c *ctx.Context) errs.Error {
	var updateUserVectorRequest UpdateUserVectorRequest
	err := c.GinContext.BindJSON(&updateUserVectorRequest)
	if err != nil {
		return errs.NewRequestError("Unable to parse request %s", err)
	}

	// check if the user already has a vector for this
	err = c.Db.Where(&data.UserVector{
		UserId:         c.SessionData.UserId,
		PreferenceType: int(updateUserVectorRequest.PreferenceType),
	}).Assign(data.UserVector{
		Sociable:    updateUserVectorRequest.Sociable,
		HardWorking: updateUserVectorRequest.Hard_Working,
		Ambitious:   updateUserVectorRequest.Ambitious,
		Energetic:   updateUserVectorRequest.Energetic,
		Carefree:    updateUserVectorRequest.Carefree,
		Confident:   updateUserVectorRequest.Confident,
	}).FirstOrCreate(&data.UserVector{}).Error

	if err != nil {
		return errs.NewRequestError("Unable to insert new user vector")
	}

	onboardingInfo, err := GetOnboardingInfo(c.Db, c.SessionData.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}

	onboardingStatus := &api.OnboardingStatus{
		onboardingInfo.State,
		onboardingInfo.UserType,
	}
	c.Result = api.OnboardingUpdateResponse{"Ok", onboardingStatus}

	return nil
}
