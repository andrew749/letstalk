package test_helpers

import (
	"fmt"

	"letstalk/server/core/survey"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

// Creates a plain user to be used in tests.
// Doesn't have anything setup yet.
func CreateTestUser(db *gorm.DB, num int) (*data.User, error) {
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
	return user, nil
}

// Creates a user that has already gone through onboarding.
func CreateTestSetupUser(db *gorm.DB, num int) (*data.User, error) {
	user, err := CreateTestUser(db, num)
	if err != nil {
		return nil, err
	}
	user.IsEmailVerified = true
	err = db.Save(user).Error
	if err != nil {
		return nil, err
	}
	cohort := &data.Cohort{
		ProgramId:   "ARTS",
		ProgramName: "Arts",
		GradYear:    2018 + uint(num),
		IsCoop:      false,
	}
	err = db.Save(cohort).Error
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

	responses := map[data.SurveyQuestionKey]data.SurveyOptionKey{
		"free_time":   "reading",
		"group_size":  "both",
		"exercise":    "daily",
		"school_work": "minimally",
		"working_on":  "school",
	}
	err = db.Save(&data.UserSurvey{
		UserId:    user.UserId,
		Group:     survey.Generic_v1.Group,
		Version:   1,
		Responses: responses,
	}).Error
	if err != nil {
		return nil, err
	}

	mentorshipPreference := 1
	bio := "yolo"
	hometown := "Richmond Hill, ON"
	additionalData := &data.UserAdditionalData{
		UserId:               user.UserId,
		MentorshipPreference: &mentorshipPreference,
		Bio:                  &bio,
		Hometown:             &hometown,
	}
	err = db.Save(additionalData).Error
	if err != nil {
		return nil, err
	}
	user.AdditionalData = additionalData
	return user, nil
}
