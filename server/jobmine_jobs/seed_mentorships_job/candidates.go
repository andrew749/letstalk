package seed_mentorships_job

import (
	"time"

	"github.com/jinzhu/gorm"

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
