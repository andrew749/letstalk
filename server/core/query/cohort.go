package query

import (
	"errors"
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic"
)

/**
 * Try to see if there is school data associated with this account.
 * If there is no data, return nil
 */
func GetUserCohort(db *gorm.DB, userId data.TUserID) (*data.Cohort, error) {
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
func GetUserCohortMappingById(db *gorm.DB, userId data.TUserID) (*data.UserCohort, error) {
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
		// NOTE: New API will allow for null sequence ids.
		sequenceId := ""
		if row.SequenceId != nil {
			sequenceId = *row.SequenceId
		}
		cohorts[i] = api.Cohort{
			row.CohortId,
			row.ProgramId,
			sequenceId,
			row.GradYear,
		}
	}

	return cohorts, nil
}

func getCohort(db *gorm.DB, cohortId data.TCohortID) (*data.Cohort, errs.Error) {
	var cohort data.Cohort
	if err := db.Where(&data.Cohort{CohortId: cohortId}).First(&cohort).Error; err != nil {
		return nil, errs.NewDbError(err)
	}
	return &cohort, nil
}

// TODO: Maybe make this return the cohort that is newly added/updated.
// TODO: Would probably be preferable to break these up in the future.
func UpdateUserCohortAndAdditionalInfo(
	db *gorm.DB,
	es *elastic.Client,
	userId data.TUserID,
	cohortId data.TCohortID,
	mentorshipPreference *int,
	bio *string,
	hometown *string,
) errs.Error {
	cohort, err := getCohort(db, cohortId)
	if err != nil {
		return err
	}

	var (
		userCohort         data.UserCohort
		userAdditionalData data.UserAdditionalData
	)

	dbErr := ctx.WithinTx(db, func(db *gorm.DB) error {
		if err := db.Where(&data.UserCohort{UserId: userId}).Assign(
			&data.UserCohort{CohortId: cohortId},
		).FirstOrCreate(&userCohort).Error; err != nil {
			return err
		}

		if err := db.Where(
			&data.UserAdditionalData{UserId: userId},
		).Assign(
			&data.UserAdditionalData{
				MentorshipPreference: mentorshipPreference,
				Bio:                  bio,
				Hometown:             hometown,
			},
		).FirstOrCreate(&userAdditionalData).Error; err != nil {
			return err
		}
		return nil
	})
	if dbErr != nil {
		return errs.NewDbError(dbErr)
	}

	userCohort.Cohort = cohort
	go indexCohortMultiTrait(es, &userCohort)

	return nil
}
