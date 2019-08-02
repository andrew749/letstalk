package match_round

import (
	"fmt"
	"letstalk/server/core/api"
	"letstalk/server/core/query"
	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/jobmine"
	"letstalk/server/test_helpers"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

type match struct {
	mentorId data.TUserID
	menteeId data.TUserID
}

type matchIndex struct {
	mentorIdIndex int
	menteeIdIndex int
}

type parameters struct {
	maxLowerYearsPerUpperYear uint
	maxUpperYearsPerLowerYear uint
	youngestUpperGradYear     uint
}

func matchesFromDataMatchRound(matchRound *data.MatchRound) []match {
	matches := make([]match, 0, len(matchRound.Matches))
	for _, mrMatch := range matchRound.Matches {
		matches = append(matches, match{
			mentorId: mrMatch.MentorUserId,
			menteeId: mrMatch.MenteeUserId,
		})
	}
	return matches
}

func matchesFromApiMatchRound(matchRound *api.MatchRound) []match {
	matches := make([]match, 0, len(matchRound.Matches))
	for _, mrMatch := range matchRound.Matches {
		matches = append(matches, match{
			mentorId: mrMatch.Mentor.User.UserId,
			menteeId: mrMatch.Mentee.User.UserId,
		})
	}
	return matches
}

func createMatchRoundTestSetup(
	t *testing.T,
	db *gorm.DB,
	groupName string,
	adminIndex int,
) (*data.ManagedGroup, []data.User) {
	users, err := test_helpers.CreateUsersForMatching(db)
	assert.NoError(t, err)
	admin := users[adminIndex]

	managedGroup, err := query.CreateManagedGroup(db, admin.UserId, groupName)
	assert.NoError(t, err)

	for _, user := range users {
		err = query.EnrollUserInManagedGroup(db, user.UserId, managedGroup.GroupId)
		assert.NoError(t, err)
	}

	return managedGroup, users
}

func checkCreateMatchRound(
	t *testing.T,
	db *gorm.DB,
	params parameters,
	expectedMatchIndexes []matchIndex,
) {
	var err error
	groupName := "WICS"
	adminIndex := 7
	managedGroup, users := createMatchRoundTestSetup(t, db, groupName, adminIndex)
	admin := users[adminIndex]

	expectedMatches := make([]match, 0, len(expectedMatchIndexes))
	for _, matchIndex := range expectedMatchIndexes {
		expectedMatches = append(expectedMatches, match{
			mentorId: users[matchIndex.mentorIdIndex].UserId,
			menteeId: users[matchIndex.menteeIdIndex].UserId,
		})
	}
	expectedParamMap := data.MatchParameters(map[string]interface{}{
		"maxLowerYearsPerUpperYear": float64(params.maxLowerYearsPerUpperYear),
		"maxUpperYearsPerLowerYear": float64(params.maxUpperYearsPerLowerYear),
		"youngestUpperGradYear":     float64(params.youngestUpperGradYear),
	})

	request := api.CreateMatchRoundRequest{
		Parameters: api.MatchRoundParameters{
			MaxLowerYearsPerUpperYear: params.maxLowerYearsPerUpperYear,
			MaxUpperYearsPerLowerYear: params.maxUpperYearsPerLowerYear,
			YoungestUpperGradYear:     params.youngestUpperGradYear,
		},
		GroupId: managedGroup.GroupId,
		UserIds: []data.TUserID{
			users[0].UserId,
			users[1].UserId,
			users[2].UserId,
			users[3].UserId,
			users[4].UserId,
			users[5].UserId,
		},
	}

	apiMatchRound, err := handleCreateMatchRound(db, admin.UserId, request)
	assert.NoError(t, err)

	matches := matchesFromApiMatchRound(apiMatchRound)
	assert.ElementsMatch(t, expectedMatches, matches)
	assert.True(t, strings.HasPrefix(apiMatchRound.Name, groupName))
	assert.Equal(t, api.MATCH_ROUND_STATE_CREATED, apiMatchRound.State)

	var dataMatchRound data.MatchRound
	err = db.Where(
		&data.MatchRound{Id: apiMatchRound.MatchRoundId},
	).Preload("Matches").Find(&dataMatchRound).Error
	assert.NoError(t, err)

	matches = matchesFromDataMatchRound(&dataMatchRound)
	assert.ElementsMatch(t, expectedMatches, matches)
	assert.True(t, strings.HasPrefix(dataMatchRound.Name, groupName))
	assert.Equal(t, dataMatchRound.MatchParameters, expectedParamMap)
	assert.Equal(t, managedGroup.GroupId, dataMatchRound.GroupId)
	assert.Nil(t, dataMatchRound.RunId)
}

func TestCreateMatchRoundControllerHappyBoundMaxLower(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			params := parameters{
				maxLowerYearsPerUpperYear: 2,
				maxUpperYearsPerLowerYear: 100,
				youngestUpperGradYear:     2021,
			}
			expectedMatches := []matchIndex{
				{3, 0},
				{4, 1},
				{5, 2},
				{4, 0},
				{3, 1},
				{5, 1},
			}
			checkCreateMatchRound(t, db, params, expectedMatches)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestCreateMatchRoundControllerHappyBoundMaxUpper(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			params := parameters{
				maxLowerYearsPerUpperYear: 100,
				maxUpperYearsPerLowerYear: 2,
				youngestUpperGradYear:     2021,
			}
			expectedMatches := []matchIndex{
				{3, 0},
				{4, 1},
				{5, 2},
				{4, 0},
				{3, 1},
				{4, 2},
			}
			checkCreateMatchRound(t, db, params, expectedMatches)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestCreateMatchRoundControllerNotAdmin(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			groupName := "WICS"
			managedGroup, users := createMatchRoundTestSetup(t, db, groupName, 6)

			request := api.CreateMatchRoundRequest{
				Parameters: api.MatchRoundParameters{
					MaxLowerYearsPerUpperYear: 100,
					MaxUpperYearsPerLowerYear: 100,
					YoungestUpperGradYear:     2021,
				},
				GroupId: managedGroup.GroupId,
				UserIds: []data.TUserID{
					users[0].UserId,
					users[1].UserId,
					users[2].UserId,
					users[3].UserId,
					users[4].UserId,
					users[5].UserId,
				},
			}
			_, err := handleCreateMatchRound(db, users[7].UserId, request)
			assert.EqualError(t, err, "You do not have rights to do this operation")
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestCreateMatchRoundControllerUserNotInGroup(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			groupName := "WICS"
			adminIndex := 7
			managedGroup, users := createMatchRoundTestSetup(t, db, groupName, adminIndex)
			admin := users[adminIndex]

			otherUser, err := test_helpers.CreateTestUser(db, 69)
			assert.NoError(t, err)

			request := api.CreateMatchRoundRequest{
				Parameters: api.MatchRoundParameters{
					MaxLowerYearsPerUpperYear: 100,
					MaxUpperYearsPerLowerYear: 100,
					YoungestUpperGradYear:     2021,
				},
				GroupId: managedGroup.GroupId,
				UserIds: []data.TUserID{
					users[0].UserId,
					users[1].UserId,
					users[2].UserId,
					users[3].UserId,
					users[4].UserId,
					users[5].UserId,
					otherUser.UserId,
				},
			}
			_, err = handleCreateMatchRound(db, admin.UserId, request)
			assert.EqualError(t, err, fmt.Sprintf("Users not in group: [%d]", otherUser.UserId))
			assert.True(t, db.Find(&data.MatchRound{}).RecordNotFound())
		},
	}
	test.RunTestWithDb(thisTest)
}

func commitMatchRoundTestSetup(
	t *testing.T,
	db *gorm.DB,
	groupName string,
	adminIndex int,
) (*data.ManagedGroup, []data.User, data.TMatchRoundID) {
	managedGroup, users := createMatchRoundTestSetup(t, db, groupName, adminIndex)
	request := api.CreateMatchRoundRequest{
		Parameters: api.MatchRoundParameters{
			MaxLowerYearsPerUpperYear: 100,
			MaxUpperYearsPerLowerYear: 100,
			YoungestUpperGradYear:     2021,
		},
		GroupId: managedGroup.GroupId,
		UserIds: []data.TUserID{
			users[0].UserId,
			users[1].UserId,
			users[2].UserId,
			users[3].UserId,
			users[4].UserId,
			users[5].UserId,
		},
	}

	matchRound, err := handleCreateMatchRound(db, users[adminIndex].UserId, request)
	assert.NoError(t, err)

	return managedGroup, users, matchRound.MatchRoundId
}

func TestCommitMatchRoundControllerHappy(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error
			groupName := "WICS"
			adminIndex := 7
			_, users, matchRoundId := commitMatchRoundTestSetup(t, db, groupName, adminIndex)
			admin := users[adminIndex]

			err = handleCommitMatchRound(db, admin.UserId, matchRoundId)
			assert.NoError(t, err)

			var matchRound data.MatchRound
			err = db.Where(
				&data.MatchRound{Id: matchRoundId},
			).Preload("CommitJob").Find(&matchRound).Error
			assert.NoError(t, err)
			assert.NotNil(t, matchRound.CommitJob)
			assert.Equal(t, float64(matchRoundId), matchRound.CommitJob.Metadata["matchRoundId"])

			// Check error saying that it already exists
			err = handleCommitMatchRound(db, admin.UserId, matchRoundId)
			assert.EqualError(t, err, fmt.Sprintf(
				"Encountered database error: Job record for match round %d already exists",
				matchRoundId))

			// Only one job should exists
			var records []jobmine.JobRecord
			err = db.Find(&records).Error
			assert.NoError(t, err)
			assert.Equal(t, 1, len(records))
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestCommitMatchRoundControllerNotAdmin(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error
			groupName := "WICS"
			_, users, matchRoundId := commitMatchRoundTestSetup(t, db, groupName, 6)

			err = handleCommitMatchRound(db, users[7].UserId, matchRoundId)
			assert.EqualError(t, err, "You do not have rights to do this operation")
		},
	}
	test.RunTestWithDb(thisTest)
}
