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

func CreateTestConnection(
	t *testing.T,
	db *gorm.DB,
	mentorUserId data.TUserID,
	menteeUserId data.TUserID,
) data.TConnectionID {
	createdAt := time.Now()
	conn := data.Connection{
		UserOneId:  mentorUserId,
		UserTwoId:  menteeUserId,
		CreatedAt:  createdAt,
		AcceptedAt: &createdAt, // Automatically accept.
	}
	assert.NoError(t, db.Create(&conn).Error)
	return conn.ConnectionId
}

func CreateTestMentorship(
	t *testing.T,
	db *gorm.DB,
	mentorUserId data.TUserID,
	connectionId data.TConnectionID,
) {
	mentorship := data.Mentorship{
		ConnectionId: connectionId,
		MentorUserId: mentorUserId,
	}
	assert.NoError(t, db.Create(&mentorship).Error)
}

func CreateTestConnectionMatchRound(
	t *testing.T,
	db *gorm.DB,
	connectionId data.TConnectionID,
	matchRoundId data.TMatchRoundID,
) {
	round := data.ConnectionMatchRound{
		ConnectionId: connectionId,
		MatchRoundId: matchRoundId,
	}
	assert.NoError(t, db.Create(&round).Error)
}

func CreateTestMatchRound(db *gorm.DB) (*data.MatchRound, error) {
	matchRound := data.MatchRound{
		Name:    "Some match round",
		GroupId: data.TGroupID("Hey friend"),
		MatchParameters: data.MatchParameters(map[string]interface{}{
			"param_a": interface{}(123),
			"param_b": interface{}(234),
		}),
	}
	if err := db.Create(&matchRound).Error; err != nil {
		return nil, err
	}
	return &matchRound, nil
}

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

func CreateUsersForMatching(db *gorm.DB) ([]data.User, error) {
	sequenceId := "8_STREAM"
	sequenceName := "8 Stream"
	now := time.Now()

	user1, err := CreateTestUser(db, 1)
	if err != nil {
		return nil, err
	}
	user1.CreatedAt = now.AddDate(0, 0, -2)
	user1.Gender = data.GENDER_FEMALE
	err = db.Save(user1).Error
	if err != nil {
		return nil, err
	}
	err = CreateCohortForUser(
		db, user1, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
		&sequenceId, &sequenceName)
	if err != nil {
		return nil, err
	}

	user2, err := CreateTestUser(db, 2)
	if err != nil {
		return nil, err
	}
	err = CreateCohortForUser(
		db, user2, "SOFTWARE_ENGINEERING", "Software Engineering", 2022, true,
		&sequenceId, &sequenceName)
	if err != nil {
		return nil, err
	}

	user3, err := CreateTestUser(db, 3)
	if err != nil {
		return nil, err
	}
	user3.CreatedAt = now.AddDate(0, 0, 2)
	err = db.Save(user3).Error
	if err != nil {
		return nil, err
	}
	err = CreateCohortForUser(
		db, user3, "COMPUTER_ENGINEERING", "Computer Engineering", 2022, true,
		&sequenceId, &sequenceName)
	if err != nil {
		return nil, err
	}

	user4, err := CreateTestUser(db, 4)
	if err != nil {
		return nil, err
	}
	user4.CreatedAt = now.AddDate(0, 0, -2)
	user4.Gender = data.GENDER_FEMALE
	err = db.Save(user4).Error
	if err != nil {
		return nil, err
	}
	err = CreateCohortForUser(
		db, user4, "SOFTWARE_ENGINEERING", "Software Engineering", 2021, true,
		&sequenceId, &sequenceName)
	if err != nil {
		return nil, err
	}

	user5, err := CreateTestUser(db, 5)
	if err != nil {
		return nil, err
	}
	err = CreateCohortForUser(
		db, user5, "SOFTWARE_ENGINEERING", "Software Engineering", 2021, true,
		&sequenceId, &sequenceName)
	if err != nil {
		return nil, err
	}

	user6, err := CreateTestUser(db, 6)
	if err != nil {
		return nil, err
	}
	user6.CreatedAt = now.AddDate(0, 0, 2)
	err = db.Save(user6).Error
	if err != nil {
		return nil, err
	}
	err = CreateCohortForUser(
		db, user6, "COMPUTER_ENGINEERING", "Computer Engineering", 2021, true,
		&sequenceId, &sequenceName)
	if err != nil {
		return nil, err
	}

	// These users should be ignored
	user7, err := CreateTestUser(db, 7)
	if err != nil {
		return nil, err
	}
	err = CreateCohortForUser(
		db, user7, "SOFTWARE_ENGINEERING", "Software Engineering", 2020, true,
		&sequenceId, &sequenceName)
	if err != nil {
		return nil, err
	}

	user8, err := CreateTestUser(db, 8)
	if err != nil {
		return nil, err
	}
	err = CreateCohortForUser(
		db, user8, "COMPUTER_ENGINEERING", "Computer Engineering", 2023, true,
		&sequenceId, &sequenceName)
	if err != nil {
		return nil, err
	}
	return []data.User{*user1, *user2, *user3, *user4, *user5, *user6, *user7, *user8}, nil
}
