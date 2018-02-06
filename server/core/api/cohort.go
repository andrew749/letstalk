package api

import (
	"database/sql"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/mijia/modelq/gmq"
)

/**
 * Api controller to get the cohort data for a specific user
 */
func GetCohortController(c *ctx.Context) errs.Error {
	userId := c.SessionData.UserId
	cohort, err := GetUserCohort(c.Db, userId)
	if err != nil {
		return errs.NewClientError("Bad request for user cohort")
	}

	c.Result = cohort
	return nil
}

/**
 * Try to see if there is school data assiociated with this account.
 * If there is no data, return nil
 */
func GetUserCohort(db *gmq.Db, userId int) (*data.Cohort, error) {
	cohortIdMapping, err := GetUserCohortMappingById(db, userId)
	cohort, err := data.CohortObjs.
		Select().
		Where(data.CohortObjs.FilterCohortId("=", cohortIdMapping.CohortId)).One(db)
	if err != nil {
		return nil, err
	}
	return &cohort, nil
}

/**
 * Get the particular cohort for a user.
 */
func GetUserCohortMappingById(db *gmq.Db, userId int) (*data.UserCohort, error) {
	cohortIdMapping, err := data.UserCohortObjs.
		Select().
		Where(data.UserCohortObjs.FilterUserId("=", userId)).
		One(db)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &cohortIdMapping, nil
}

/**
 * Get all the current cohorts known.
 */
func GetAllCohorts(db *gmq.Db) ([]data.Cohort, error) {
	cohorts, err := data.CohortObjs.Select().List(db)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return cohorts, nil
}
