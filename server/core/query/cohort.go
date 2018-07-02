package query

import (
	"errors"
	"letstalk/server/core/api"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

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

func GetAllCohorts(db *gorm.DB) ([]api.Cohort, errs.Error) {
	var rows []data.Cohort
	if err := db.Find(&rows).Error; err != nil {
		return nil, errs.NewDbError(err)
	}
	cohorts := make([]api.Cohort, len(rows))
	for i, row := range rows {
		cohorts[i] = api.Cohort{row.CohortId, row.ProgramId, row.SequenceId, row.GradYear}
	}

	return cohorts, nil
}
