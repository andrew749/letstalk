package test_helpers

import (
	"fmt"
	"testing"
	"time"

	"letstalk/server/core/ctx"
	"letstalk/server/core/sessions"
	"letstalk/server/core/survey"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func CreateTestContext(t *testing.T, db *gorm.DB, userId data.TUserID) *ctx.Context {
	expiry := time.Now()
	expiry = expiry.AddDate(1, 0, 0)
	sessionData, err := sessions.CreateSessionData(userId, expiry)
	assert.NoError(t, err)
	return ctx.NewContext(nil, db, nil, sessionData, nil)
}

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

func CreateCohortForUser(
	db *gorm.DB,
	user *data.User,
	programId string,
	programName string,
	gradYear uint,
	isCoop bool,
	sequenceId *string,
	sequenceName *string,
) error {
	cohort := &data.Cohort{
		ProgramId:    programId,
		ProgramName:  programName,
		GradYear:     gradYear,
		IsCoop:       isCoop,
		SequenceId:   sequenceId,
		SequenceName: sequenceName,
	}
	err := db.Save(cohort).Error
	if err != nil {
		return err
	}

	userCohort := &data.UserCohort{
		UserId:   user.UserId,
		CohortId: cohort.CohortId,
	}
	err = db.Save(userCohort).Error
	if err != nil {
		return err
	}
	userCohort.Cohort = cohort
	user.Cohort = userCohort
	return nil
}

func CreateSurveyForUser(
	db *gorm.DB,
	user *data.User,
	responses map[data.SurveyQuestionKey]data.SurveyOptionKey,
	group data.SurveyGroup,
	version int,
) error {
	userSurvey := data.UserSurvey{
		UserId:    user.UserId,
		Group:     group,
		Version:   version,
		Responses: responses,
	}
	err := db.Save(&userSurvey).Error
	if err != nil {
		return err
	}

	if user.UserSurveys == nil {
		user.UserSurveys = []data.UserSurvey{userSurvey}
	} else {
		user.UserSurveys = append(user.UserSurveys, userSurvey)
	}

	return nil
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
	err = CreateCohortForUser(db, user, "ARTS", "Arts", 2018+uint(num), false, nil, nil)

	responses := map[data.SurveyQuestionKey]data.SurveyOptionKey{
		"free_time":   "reading",
		"group_size":  "both",
		"exercise":    "daily",
		"school_work": "minimally",
		"working_on":  "school",
	}
	err = CreateSurveyForUser(db, user, responses, survey.Generic_v1.Group, 1)
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
