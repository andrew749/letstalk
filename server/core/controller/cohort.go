package controller

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
)

/**
 * Api controller to get the cohort data for a specific user
 */
func GetCohortController(c *ctx.Context) errs.Error {
	userId := c.SessionData.UserId
	cohort, err := query.GetUserCohort(c.Db, userId)
	if err != nil {
		return errs.NewRequestError(err.Error())
	}

	c.Result = cohort
	return nil
}

func GetAllCohortsController(c *ctx.Context) errs.Error {
	cohorts, err := query.GetAllCohorts(c.Db)
	if err != nil {
		return err
	}

	c.Result = cohorts
	return nil
}
