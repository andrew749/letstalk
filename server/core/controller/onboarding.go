package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
)

/**
*
* Sample post request:
 {
 	"cohortId": UserVectorType
 }
*/

// Update a user with new information for their school
// try to match this data to an existing sequence.
func UpdateUserCohortAndAdditionalInfo(c *ctx.Context) errs.Error {
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

	c.Result = "Success"
	return nil
}
