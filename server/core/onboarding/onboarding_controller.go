package onboarding

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

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

type UpdateCohortRequest struct {
	CohortId int `json:"cohortId" binding:"required"`
}

type OnboardingUpdateResponse struct {
	Message          string            `json:"message" binding:"required"`
	OnboardingStatus *OnboardingStatus `json:"onboardingStatus" binding:"required"`
}

func isValidCohort(db *gorm.DB, cohortId int) bool {
	var numCohorts int = 0
	db.Model(&data.Cohort{}).Where("cohort_id = ?", cohortId).Count(&numCohorts)
	return numCohorts > 0
}

// Update a user with new information for their school
// try to match this data to an existing sequence.
func UpdateUserCohort(c *ctx.Context) errs.Error {
	var newCohortRequest UpdateCohortRequest
	err := c.GinContext.BindJSON(&newCohortRequest)

	if err != nil {
		return errs.NewClientError("%s", err)
	}

	newCohortId := newCohortRequest.CohortId

	// check that the new cohort is valid
	if !isValidCohort(c.Db, newCohortId) {
		return errs.NewClientError("Unknown cohort.")
	}

	userId := c.SessionData.UserId
	userCohort, err := api.GetUserCohortMappingById(c.Db, userId)

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return errs.NewDbError(err)
	}

	var (
		dbErr          error
		successMessage string
	)

	tx := c.Db.Begin()

	// if the user doesnt have a cohort
	if userCohort == nil {
		rlog.Debug("No cohort found for user. Adding cohort.")
		// insert new data from the request
		userCohort = &data.UserCohort{
			UserId:   userId,
			CohortId: newCohortId,
		}

		dbErr = tx.Create(&userCohort).Error
	} else {
		rlog.Debug("Updating cohort information for user.")
		userCohort.CohortId = newCohortId
		// update the cohort data from the request
		dbErr = c.Db.Save(&userCohort).Error
	}

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

	onboardingStatus := &OnboardingStatus{
		onboardingInfo.State,
		onboardingInfo.UserType,
	}
	c.Result = OnboardingUpdateResponse{successMessage, onboardingStatus}
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
		return errs.NewClientError("Unable to parse request %s", err)
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
		return errs.NewClientError("Unable to insert new user vector")
	}

	onboardingInfo, err := GetOnboardingInfo(c.Db, c.SessionData.UserId)
	if err != nil {
		return errs.NewDbError(err)
	}

	onboardingStatus := &OnboardingStatus{
		onboardingInfo.State,
		onboardingInfo.UserType,
	}
	c.Result = OnboardingUpdateResponse{"Ok", onboardingStatus}

	return nil
}
