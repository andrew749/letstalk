package onboarding

import (
	"database/sql"
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/mijia/modelq/gmq"
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

func isValidCohort(db *gmq.Db, cohortId int) bool {
	_, err := data.CohortObjs.
		Select().
		Where(data.CohortObjs.FilterCohortId("=", cohortId)).
		One(db)

	return err == nil
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
		rlog.Debug("Invalid cohort.")
		return errs.NewClientError("Unknown cohort.")
	}

	userId := c.SessionData.UserId
	userCohort, err := api.GetUserCohortMappingById(c.Db, userId)

	if err != nil && err != sql.ErrNoRows {
		return errs.NewDbError(err)
	}

	var (
		dbErr          error
		successMessage string
	)

	// if the user doesnt have a cohort
	if userCohort == nil {
		// insert new data from the request
		userCohort = &data.UserCohort{
			UserId:   userId,
			CohortId: newCohortId,
		}

		// try to insert the data
		dbErr = gmq.WithinTx(c.Db, func(tx *gmq.Tx) error {
			_, err = userCohort.Insert(tx)
			if err != nil {
				return err
			}
			return nil
		})
		successMessage = "Successfully added cohort to user."
	} else {
		userCohort.CohortId = newCohortId
		// update the cohort data from the request
		dbErr = gmq.WithinTx(c.Db, func(tx *gmq.Tx) error {
			_, err = userCohort.Update(tx)
			if err != nil {
				return err
			}
			return nil
		})
		successMessage = "Successfully updated cohort for user."
	}

	if dbErr != nil {
		return errs.NewDbError(dbErr)
	}

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
	PreferenceType api.UserVectorPreferenceType `json:"isMenteePreference" binding:"exists"`
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
	stmt, err := c.Db.Prepare(
		`
		REPLACE INTO user_vector (
			user_id,
			preference_type,
			sociable,
			hard_working,
			ambitious,
			energetic,
			carefree,
			confident
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`,
	)

	if err != nil {
		return errs.NewDbError(err)
	}

	_, err = stmt.Query(
		c.SessionData.UserId,
		int(updateUserVectorRequest.PreferenceType),
		updateUserVectorRequest.Sociable,
		updateUserVectorRequest.Hard_Working,
		updateUserVectorRequest.Ambitious,
		updateUserVectorRequest.Energetic,
		updateUserVectorRequest.Carefree,
		updateUserVectorRequest.Confident,
	)
	if err != nil {
		return errs.NewInternalError("Unable to insert new vector", err)
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
