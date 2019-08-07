package match_round

import (
	"fmt"
	"letstalk/server/core/api"
	"letstalk/server/core/query"
	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/jobmine"
	"letstalk/server/jobmine_jobs/match_round_commit_job"
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
				{mentorIdIndex: 3, menteeIdIndex: 0},
				{mentorIdIndex: 4, menteeIdIndex: 1},
				{mentorIdIndex: 5, menteeIdIndex: 2},
				{mentorIdIndex: 4, menteeIdIndex: 0},
				{mentorIdIndex: 3, menteeIdIndex: 1},
				{mentorIdIndex: 5, menteeIdIndex: 1},
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
				{mentorIdIndex: 3, menteeIdIndex: 0},
				{mentorIdIndex: 4, menteeIdIndex: 1},
				{mentorIdIndex: 5, menteeIdIndex: 2},
				{mentorIdIndex: 4, menteeIdIndex: 0},
				{mentorIdIndex: 3, menteeIdIndex: 1},
				{mentorIdIndex: 4, menteeIdIndex: 2},
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

var specStore = jobmine.JobSpecStore{
	JobSpecs: map[jobmine.JobType]jobmine.JobSpec{
		match_round_commit_job.MATCH_ROUND_COMMIT_JOB: match_round_commit_job.CommitJobSpec,
	},
}

func checkMatchUser(t *testing.T, user *api.MatchUser, userMap map[data.TUserID]data.User) {
	eUser := userMap[user.User.UserId]
	assert.Equal(t, eUser.FirstName, user.User.FirstName)
	assert.Equal(t, eUser.LastName, user.User.LastName)
	assert.Equal(t, eUser.Email, user.Email)
	assert.Equal(t, eUser.Cohort.Cohort.ProgramName, user.Cohort.ProgramName)
	assert.Equal(t, *eUser.Cohort.Cohort.SequenceName, *user.Cohort.SequenceName)
	assert.Equal(t, eUser.Cohort.Cohort.GradYear, user.Cohort.GradYear)
}

// Pretty big integration test of pretty much all the functionality of the match rounds module.
func TestGetMatchRoundsControllerHappy(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error
			groupName := "WICS"
			adminIndex := 7
			managedGroup, users := createMatchRoundTestSetup(t, db, groupName, adminIndex)
			admin := users[adminIndex]

			request := api.CreateMatchRoundRequest{
				Parameters: api.MatchRoundParameters{
					MaxLowerYearsPerUpperYear: 1,
					MaxUpperYearsPerLowerYear: 1,
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

			numRounds := 5
			matchRounds := make([]api.MatchRound, numRounds)
			// 0 - Created
			// 1 - Deleted
			// 2 - Committing
			// 3 - Committed
			// 4 - Failed

			for i := 0; i < numRounds; i++ {
				matchRound, err := handleCreateMatchRound(db, admin.UserId, request)
				assert.NoError(t, err)
				matchRounds[i] = *matchRound
			}

			err = handleDeleteMatchRound(db, admin.UserId, matchRounds[1].MatchRoundId)
			assert.NoError(t, err)

			for i := 2; i < numRounds; i++ {
				err = handleCommitMatchRound(db, admin.UserId, matchRounds[i].MatchRoundId)
				assert.NoError(t, err)
			}

			// Forcing the statuses for two of the jobs
			// TODO(wojtek): I noticed when trying to actually run the job that there is a problem
			// when running tests that rely on SendGrid offline. We should remove this depenedency during
			// tests.
			// Need to get run id for 3rd and 4th match rounds which should exist
			var matchRound3 data.MatchRound
			err = db.Where(
				&data.MatchRound{Id: matchRounds[3].MatchRoundId},
			).Preload("CommitJob").Find(&matchRound3).Error
			assert.NoError(t, err)
			assert.NotNil(t, matchRound3.CommitJob)
			matchRound3.CommitJob.Status = jobmine.STATUS_SUCCESS
			err = db.Save(matchRound3.CommitJob).Error
			assert.NoError(t, err)

			var matchRound4 data.MatchRound
			err = db.Where(
				&data.MatchRound{Id: matchRounds[4].MatchRoundId},
			).Preload("CommitJob").Find(&matchRound4).Error
			assert.NoError(t, err)
			assert.NotNil(t, matchRound4.CommitJob)
			matchRound4.CommitJob.Status = jobmine.STATUS_FAILED
			err = db.Save(matchRound4.CommitJob).Error
			assert.NoError(t, err)

			response, err := handleGetMatchRounds(db, admin.UserId, managedGroup.GroupId)
			assert.NoError(t, err)
			assert.Len(t, response, 4)

			resMap := make(map[data.TMatchRoundID]api.MatchRound)
			for _, matchRound := range response {
				resMap[matchRound.MatchRoundId] = matchRound
			}

			// Check states
			assert.Equal(t, api.MATCH_ROUND_STATE_CREATED, resMap[matchRounds[0].MatchRoundId].State)
			assert.Equal(t, api.MATCH_ROUND_STATE_COMMITTING, resMap[matchRounds[2].MatchRoundId].State)
			assert.Equal(t, api.MATCH_ROUND_STATE_COMMITTED, resMap[matchRounds[3].MatchRoundId].State)
			assert.Equal(t, api.MATCH_ROUND_STATE_FAILED, resMap[matchRounds[4].MatchRoundId].State)

			// Check matches in rounds
			expectedMatches := []match{
				{mentorId: users[3].UserId, menteeId: users[0].UserId},
				{mentorId: users[4].UserId, menteeId: users[1].UserId},
				{mentorId: users[5].UserId, menteeId: users[2].UserId},
			}
			userMap := make(map[data.TUserID]data.User)
			for _, user := range users {
				userMap[user.UserId] = user
			}

			for _, matchRound := range response {
				matches := matchesFromApiMatchRound(&matchRound)
				assert.ElementsMatch(t, expectedMatches, matches)

				for _, match := range matchRound.Matches {
					checkMatchUser(t, &match.Mentee, userMap)
					checkMatchUser(t, &match.Mentor, userMap)
				}
			}

			// Test that you can't delete committing, committed, failed match rounds
			err = handleDeleteMatchRound(db, admin.UserId, matchRounds[2].MatchRoundId)
			assert.EqualError(t, err, fmt.Sprintf("Cannot delete match round in %s state",
				api.MATCH_ROUND_STATE_COMMITTING))
			err = handleDeleteMatchRound(db, admin.UserId, matchRounds[3].MatchRoundId)
			assert.EqualError(t, err, fmt.Sprintf("Cannot delete match round in %s state",
				api.MATCH_ROUND_STATE_COMMITTED))
			err = handleDeleteMatchRound(db, admin.UserId, matchRounds[4].MatchRoundId)
			assert.EqualError(t, err, fmt.Sprintf("Cannot delete match round in %s state",
				api.MATCH_ROUND_STATE_FAILED))
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetMatchRoundsControllerNotAdmin(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			groupName := "WICS"
			managedGroup, users := createMatchRoundTestSetup(t, db, groupName, 6)

			_, err := handleGetMatchRounds(db, users[7].UserId, managedGroup.GroupId)
			assert.EqualError(t, err, "You do not have rights to do this operation")
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestDeleteMatchRoundControllerNotAdmin(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			groupName := "WICS"
			_, users, matchRoundId := commitMatchRoundTestSetup(t, db, groupName, 6)

			err := handleDeleteMatchRound(db, users[7].UserId, matchRoundId)
			assert.EqualError(t, err, "You do not have rights to do this operation")
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetGroupMembersControllerNotAdmin(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			groupName := "WICS"
			managedGroup, users := createMatchRoundTestSetup(t, db, groupName, 6)

			_, err := handleGetGroupMembers(db, users[7].UserId, managedGroup.GroupId)
			assert.EqualError(t, err, "You do not have rights to do this operation")
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetGroupMembersControllerHappy(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			groupName := "WICS"

			userAdmin, err := test_helpers.CreateTestUser(db, 1)
			assert.NoError(t, err)

			managedGroup, err := query.CreateManagedGroup(db, userAdmin.UserId, groupName)
			assert.NoError(t, err)

			userSignedUp, err := test_helpers.CreateTestUser(db, 2)
			assert.NoError(t, err)

			userEmailVerified, err := test_helpers.CreateTestUser(db, 3)
			assert.NoError(t, err)
			userEmailVerified.IsEmailVerified = true
			err = db.Save(userEmailVerified).Error
			assert.NoError(t, err)

			userBasicInfo, err := test_helpers.CreateTestUser(db, 4)
			assert.NoError(t, err)
			userBasicInfo.IsEmailVerified = true
			err = db.Save(userBasicInfo).Error
			assert.NoError(t, err)
			cohort := &data.Cohort{
				ProgramId:   "ARTS",
				ProgramName: "Arts",
				GradYear:    2018,
				IsCoop:      false,
			}
			err = db.Save(cohort).Error
			assert.NoError(t, err)
			userCohort := &data.UserCohort{
				UserId:   userBasicInfo.UserId,
				CohortId: cohort.CohortId,
			}
			err = db.Save(userCohort).Error
			assert.NoError(t, err)

			userSetup, err := test_helpers.CreateTestSetupUser(db, 5)
			assert.NoError(t, err)

			// False connection not attached to group
			connection := &data.Connection{
				UserOneId: userSetup.UserId,
				UserTwoId: userAdmin.UserId,
			}
			err = db.Save(connection).Error
			assert.NoError(t, err)

			userConn1, err := test_helpers.CreateTestSetupUser(db, 6)
			assert.NoError(t, err)
			userConn2, err := test_helpers.CreateTestSetupUser(db, 7)
			assert.NoError(t, err)

			memberUsers := []data.User{
				*userSignedUp,
				*userEmailVerified,
				*userBasicInfo,
				*userSetup,
				*userConn1,
				*userConn2,
			}

			for _, user := range memberUsers {
				err = query.EnrollUserInManagedGroup(db, user.UserId, managedGroup.GroupId)
				assert.NoError(t, err)
			}

			matchRound := &data.MatchRound{GroupId: managedGroup.GroupId}
			err = db.Save(matchRound).Error
			assert.NoError(t, err)
			connection = &data.Connection{
				UserOneId: userConn1.UserId,
				UserTwoId: userConn2.UserId,
				ConnectionMatchRound: &data.ConnectionMatchRound{
					MatchRoundId: matchRound.Id,
				},
			}
			err = db.Save(connection).Error
			assert.NoError(t, err)

			groupMembers, err := handleGetGroupMembers(db, userAdmin.UserId, managedGroup.GroupId)
			assert.NoError(t, err)
			assert.Len(t, groupMembers, len(memberUsers))

			memberMap := make(map[data.TUserID]api.GroupMember)
			for _, groupMember := range groupMembers {
				memberMap[groupMember.User.UserId] = groupMember
			}

			assert.Equal(t, api.GROUP_MEMBER_STATUS_SIGNED_UP,
				memberMap[userSignedUp.UserId].Status)
			assert.Equal(t, api.GROUP_MEMBER_STATUS_SIGNED_UP,
				memberMap[userEmailVerified.UserId].Status)
			assert.Equal(t, api.GROUP_MEMBER_STATUS_SIGNED_UP,
				memberMap[userBasicInfo.UserId].Status)
			assert.Equal(t, api.GROUP_MEMBER_STATUS_ONBOARDED,
				memberMap[userSetup.UserId].Status)
			assert.Equal(t, api.GROUP_MEMBER_STATUS_MATCHED,
				memberMap[userConn1.UserId].Status)
			assert.Equal(t, api.GROUP_MEMBER_STATUS_MATCHED,
				memberMap[userConn2.UserId].Status)

			for _, user := range memberUsers {
				groupMember := memberMap[user.UserId]
				assert.Equal(t, user.FirstName, groupMember.User.FirstName)
				assert.Equal(t, user.LastName, groupMember.User.LastName)
				assert.Equal(t, user.Email, groupMember.Email)
				if user.Cohort != nil && user.Cohort.Cohort != nil {
					assert.Equal(t, user.Cohort.Cohort.ProgramName, groupMember.Cohort.ProgramName)
					assert.Equal(t, user.Cohort.Cohort.GradYear, groupMember.Cohort.GradYear)
					if user.Cohort.Cohort.SequenceName == nil {
						assert.Nil(t, user.Cohort.Cohort.SequenceName)
					} else {
						assert.Equal(t, *user.Cohort.Cohort.SequenceName, *groupMember.Cohort.SequenceName)
					}
				}
			}
		},
	}
	test.RunTestWithDb(thisTest)
}
