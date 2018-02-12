package onboarding

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/mijia/modelq/gmq"
)

// TODO(acod): make this flow dynamic in nature
// i.e. fetch a onboarding type and the posible options

/**
*
* Sample post request:
 {
 	"cohortId": string
 }
*/

type CohortUpdateRequest struct {
	CohortId int `json:"cohortId" binding:"required"`
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
	userId := c.SessionData.UserId
	userCohort, err := api.GetUserCohortMappingById(c.Db, userId)

	if err != nil {
		return errs.NewDbError(err)
	}

	var dbErr error
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
