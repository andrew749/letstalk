package seed_mentorships_job

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"letstalk/server/data"
)

const joinStr = `INNER JOIN (
	SELECT user_id, program_id, grad_year FROM user_cohorts
	JOIN cohorts ON user_cohorts.cohort_id = cohorts.cohort_id
) AS cohort ON users.user_id = cohort.user_id`

// Finds users that meet the following conditions:
// - are in any of the given programIds
// - if isLowerYear: graduating any year greater than youngestUpperYear
// - if not isLowerYear: graduating any year less than or equal to youngestUpperYear
// - optionally created within the [createdAfter, createdBefore] range
func GetCandidates(
	db *gorm.DB,
	programIds []string,
	isLowerYear bool,
	youngestUpperYear uint,
	createdAfter *time.Time,
	createdBefore *time.Time,
) ([]data.TUserID, error) {
	query := db.Model(&data.User{}).Joins(joinStr).Where("cohort.program_id IN (?)", programIds)
	if isLowerYear {
		query = query.Where("cohort.grad_year > ?", youngestUpperYear)
	} else {
		query = query.Where("cohort.grad_year <= ?", youngestUpperYear)
	}

	if createdAfter != nil {
		query = query.Where("users.created_at >= ?", *createdAfter)
	}
	if createdBefore != nil {
		query = query.Where("users.created_at <= ?", *createdBefore)
	}

	var users []data.User
	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	userIds := make([]data.TUserID, len(users))
	for i, user := range users {
		userIds[i] = user.UserId
	}

	return userIds, nil
}

// Gets users for both lower and upper years.
// Term start and end times only apply to lower years, since we already downweight upper years
// created out of term during ranking.
func GetLowerUpperYears(
	db *gorm.DB,
	programIds []string,
	youngestUpperYear uint,
	termStartTime *time.Time,
	termEndTime *time.Time,
) ([]data.TUserID, error) {
	lowerYearIds, err := GetCandidates(
		db, programIds, true, youngestUpperYear, termStartTime, termEndTime)
	if err != nil {
		return nil, err
	}
	upperYearIds, err := GetCandidates(db, programIds, false, youngestUpperYear, nil, nil)
	if err != nil {
		return nil, err
	}
	lowerYearIdSet := make(map[data.TUserID]interface{})
	for _, lowerYearId := range lowerYearIds {
		lowerYearIdSet[lowerYearId] = nil
	}

	overlapList := make([]data.TUserID, 0)
	for _, upperYearId := range upperYearIds {
		if _, ok := lowerYearIdSet[upperYearId]; ok {
			overlapList = append(overlapList, upperYearId)
		}
	}

	if len(overlapList) > 0 {
		return nil, errors.New(fmt.Sprintf(
			"Overlapping users between upper and lower year: %v", overlapList))
	}
	userIds := append(lowerYearIds, upperYearIds...)
	return userIds, nil
}
