package query

import (
	"fmt"

	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func CreateTestUser(db *gorm.DB, num int) (*data.User, error) {
	cohort := &data.Cohort{
		ProgramId:   "ARTS",
		ProgramName: "Arts",
		GradYear:    2018 + uint(num),
		IsCoop:      false,
	}
	err := db.Save(cohort).Error
	if err != nil {
		return nil, err
	}

	birthdate := "1996-11-07"
	user, err := data.CreateUser(
		db,
		fmt.Sprintf("john.doe%d@gmail.com", num),
		fmt.Sprintf("John%d", num),
		fmt.Sprintf("Doe%d", num),
		data.GENDER_MALE,
		&birthdate,
		data.USER_ROLE_DEFAULT,
	)
	if err != nil {
		return nil, err
	}

	userCohort := &data.UserCohort{
		UserId:   user.UserId,
		CohortId: cohort.CohortId,
	}
	err = db.Save(userCohort).Error
	if err != nil {
		return nil, err
	}
	userCohort.Cohort = cohort
	user.Cohort = userCohort
	return user, nil
}
