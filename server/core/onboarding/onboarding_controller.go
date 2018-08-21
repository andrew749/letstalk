package onboarding

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"letstalk/server/core/api"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

// i.e. fetch a onboarding type and the possible options

/**
*
* Sample post request:
 {
 	"cohortId": UserVectorType
 }
*/

// TODO: Move into `query`
func isValidCohort(db *gorm.DB, cohortId data.TCohortID) bool {
	var numCohorts int = 0
	db.Model(&data.Cohort{}).Where("cohort_id = ?", cohortId).Count(&numCohorts)
	return numCohorts > 0
}

// Update a user with new information for their school
// try to match this data to an existing sequence.
func UpdateUserCohort(c *ctx.Context) errs.Error {
	var newCohortRequest api.UpdateCohortRequest

	if err := c.GinContext.BindJSON(&newCohortRequest); err != nil {
		return errs.NewRequestError("%s", err.Error())
	}

	newCohortId := newCohortRequest.CohortId

	// check that the new cohort is valid
	if !isValidCohort(c.Db, newCohortId) {
		return errs.NewRequestError("Unknown cohort: %s", newCohortId)
	}

	userId := c.SessionData.UserId

	var (
		dbErr          error
		successMessage string
	)

	tx := c.Db.Begin()

	rlog.Debug("No cohort found for user. Adding cohort.")
	// insert new data from the request
	var (
		userCohort         data.UserCohort
		userAdditionalData data.UserAdditionalData
	)

	dbErr = tx.Where(&data.UserCohort{UserId: userId}).Assign(
		&data.UserCohort{CohortId: newCohortId},
	).FirstOrCreate(&userCohort).Error

	if dbErr != nil {
		tx.Rollback()
		return errs.NewInternalError(dbErr.Error())
	}

	dbErr = tx.Where(
		&data.UserAdditionalData{UserId: userId},
	).Assign(
		&data.UserAdditionalData{
			MentorshipPreference: newCohortRequest.MentorshipPreference,
			Bio:                  newCohortRequest.Bio,
			Hometown:             newCohortRequest.Hometown,
		},
	).FirstOrCreate(&userAdditionalData).Error

	if dbErr != nil {
		tx.Rollback()
		return errs.NewInternalError(dbErr.Error())
	}

	successMessage = "Successfully added cohort to user."
	tx.Commit()

	onboardingInfo, err := GetOnboardingInfo(c.Db, c.SessionData.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}

	onboardingStatus := &api.OnboardingStatus{
		onboardingInfo.State,
		onboardingInfo.UserType,
	}
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
