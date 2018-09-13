package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"

	"github.com/jinzhu/gorm"
)

func SimpleTraitUserSearchController(c *ctx.Context) errs.Error {
	var req api.SimpleTraitUserSearchRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	res, err := query.SearchUsersBySimpleTrait(c.Db, req, c.SessionData.UserId)
	if err != nil {
		return errs.NewEsError(err)
	}

	c.Result = res
	return nil
}

func PositionUserSearchController(c *ctx.Context) errs.Error {
	var req api.PositionUserSearchRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	res, err := query.SearchUsersByPosition(c.Db, req, c.SessionData.UserId)
	if err != nil {
		return errs.NewEsError(err)
	}

	c.Result = res
	return nil
}

func CohortUserSearchController(c *ctx.Context) errs.Error {
	var req api.CohortUserSearchRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	res, err := query.SearchUsersByCohort(c.Db, req, c.SessionData.UserId)
	if err != nil {
		return errs.NewEsError(err)
	}

	c.Result = res
	return nil
}

func MyCohortUserSearchController(c *ctx.Context) errs.Error {
	var req api.CommonUserSearchRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	cohort, err := query.GetUserCohortMappingById(c.Db, c.SessionData.UserId)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errs.NewRequestError("Set a cohort to see other users in your class")
		}
		return errs.NewDbError(err)
	}

	cohortReq := api.CohortUserSearchRequest{
		CommonUserSearchRequest: req,
		CohortId:                cohort.CohortId,
	}
	res, err := query.SearchUsersByCohort(c.Db, cohortReq)
	if err != nil {
		return errs.NewEsError(err)
	}

	c.Result = res
	return nil
}
