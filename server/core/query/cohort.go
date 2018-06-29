package query

import (
	"errors"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

/**
 * Api controller to get the cohort data for a specific user
 */
func GetCohortController(c *ctx.Context) errs.Error {
	userId := c.SessionData.UserId
	cohort, err := GetUserCohort(c.Db, userId)
	if err != nil {
		return errs.NewRequestError(err.Error())
	}

	c.Result = cohort
	return nil
}

/**
 * Try to see if there is school data associated with this account.
 * If there is no data, return nil
 */
func GetUserCohort(db *gorm.DB, userId int) (*data.Cohort, error) {
	cohortIdMapping, err := GetUserCohortMappingById(db, userId)

	if err != nil {
		return nil, err
	}

	var cohort data.Cohort
	if db.Where("cohort_id = ?", cohortIdMapping.CohortId).First(&cohort).RecordNotFound() {
		return nil, errors.New("Unable to find cohort")
	}
	return &cohort, nil
}

/**
 * Get the particular cohort for a user.
 */
func GetUserCohortMappingById(db *gorm.DB, userId int) (*data.UserCohort, error) {
	var cohort data.UserCohort
	if err := db.Where("user_id = ?", userId).First(&cohort).Error; err != nil {
		return nil, err
	}
	return &cohort, nil
}
