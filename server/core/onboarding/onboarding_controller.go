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

// i.e. fetch a onboarding type and the posible options

/**
*
* Sample post request:
 {
 	"cohortId": UserVectorType
 }
*/

type CohortUpdateRequest struct {
	CohortId int `json:"cohortId" binding:"required"`
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
	var newCohortRequest CohortUpdateRequest
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

	var dbErr error

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
		c.Result = "Successfully added cohort to user."
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
		c.Result = "Successfully updated cohort for user."
	}

	if dbErr != nil {
		return errs.NewDbError(dbErr)
	}

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

type UserVectorPreferenceType int

const (
	MenteePreference UserVectorPreferenceType = iota
	MentorPreference
)

type UpdateUserVectorRequest struct {
	PreferenceType int `json:"isMenteePreference" binding:"required"`
	Sociable       int `json:"sociable" binding:"required"`
	Hard_Working   int `json:"hardWorking" binding:"required"`
	Ambitious      int `json:"ambitious" binding:"required"`
	Energetic      int `json:"energetic" binding:"required"`
	Carefree       int `json:"carefree" binding:"required"`
	Confident      int `json:"confident" binding:"required"`
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

	c.Result = struct{ Status string }{"Ok"}

	return nil
}
