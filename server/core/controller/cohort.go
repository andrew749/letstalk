package controller

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
)

/**
 * Returns all cohorts in the db
 */
func GetAllCohortsController(c *ctx.Context) errs.Error {
	cohorts, err := query.GetAllCohorts(c.Db)
	if err != nil {
		return err
	}

	c.Result = cohorts
	return nil
}
