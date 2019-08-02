package match_round

import (
	"fmt"
	"letstalk/server/core/api"
	"letstalk/server/core/converters"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"
	"letstalk/server/jobmine"
	"letstalk/server/jobmine_jobs/match_round_commit_job"
	"letstalk/server/recommendations"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

// Controller for the create_match_round admin endpoint
func CreateMatchRoundController(c *ctx.Context) errs.Error {
	var request api.CreateMatchRoundRequest
	if err := c.GinContext.BindJSON(&request); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}

	matchRound, err := handleCreateMatchRound(
		c.Db,
		c.SessionData.UserId,
		request,
	)
	if err != nil {
		return err
	}

	c.Result = matchRound
	return nil
}

func handleCreateMatchRound(
	db *gorm.DB,
	adminId data.TUserID,
	req api.CreateMatchRoundRequest,
) (*api.MatchRound, errs.Error) {
	if err := checkIsAdmin(db, adminId, req.GroupId); err != nil {
		return nil, err
	}

	if req.UserIds == nil {
		return nil, errs.NewRequestError("Expected non-nil user ids")
	}

	if err := checkUsersInGroup(db, req.UserIds, req.GroupId); err != nil {
		return nil, err
	}

	strat := recommendations.MentorMenteeStrat(
		req.Parameters.MaxLowerYearsPerUpperYear,
		req.Parameters.MaxUpperYearsPerLowerYear,
		req.Parameters.YoungestUpperGradYear,
	)

	fetcherOptions := recommendations.UserFetcherOptions{UserIds: req.UserIds}
	matches, err := recommendations.Recommend(db, fetcherOptions, strat)

	if err != nil {
		errStr := fmt.Sprintf("Error when generating matches, %+v", err)
		rlog.Error(errStr)
		return nil, errs.NewRequestError(errStr)
	}

	if len(matches) == 0 {
		return nil, errs.NewRequestError("Parameters result in no matches")
	}

	parameters := createMatchParameters(
		req.Parameters.MaxLowerYearsPerUpperYear,
		req.Parameters.MaxUpperYearsPerLowerYear,
		req.Parameters.YoungestUpperGradYear,
	)

	matchRound, err := createMatchRound(
		db,
		req.GroupId,
		matches,
		parameters,
	)

	return matchRound, nil
}

// Controller for commit_match_round admin endpoint
// Creates the job to commit the match round, updates the match round model and returns.
func CommitMatchRoundController(c *ctx.Context) errs.Error {
	var request api.CommitMatchRoundRequest
	if err := c.GinContext.BindJSON(&request); err != nil {
		return errs.NewRequestError("Failed to parse input")
	}

	err := handleCommitMatchRound(
		c.Db,
		c.SessionData.UserId,
		request.MatchRoundId,
	)
	if err != nil {
		return err
	}

	c.Result = "Success"
	return nil
}

func handleCommitMatchRound(
	db *gorm.DB,
	adminId data.TUserID,
	matchRoundId data.TMatchRoundID,
) errs.Error {
	if err := checkIsAdminMatchRound(db, adminId, matchRoundId); err != nil {
		return err
	}

	err := ctx.WithinTx(db, func(db *gorm.DB) error {
		var matchRound data.MatchRound
		if err := db.Where(&data.MatchRound{Id: matchRoundId}).Find(&matchRound).Error; err != nil {
			return err
		}

		runId, err := match_round_commit_job.CreateCommitJob(db, matchRoundId)
		if err != nil {
			return err
		}

		matchRound.RunId = runId
		if err := db.Save(&matchRound).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return errs.NewDbError(err)
	}

	return nil
}

// Controller for GET match_rounds admin endpoint
// Returns match rounds for a given group, including matches in that match round and its status
func GetMatchRoundsController(c *ctx.Context) errs.Error {
	groupId := data.TGroupID(c.GinContext.Param("groupId"))

	matchRounds, err := handleGetMatchRounds(
		c.Db,
		c.SessionData.UserId,
		groupId,
	)
	if err != nil {
		return err
	}

	c.Result = matchRounds
	return nil
}

func handleGetMatchRounds(
	db *gorm.DB,
	adminId data.TUserID,
	groupId data.TGroupID,
) ([]api.MatchRound, errs.Error) {
	if err := checkIsAdmin(db, adminId, groupId); err != nil {
		return nil, err
	}

	var matchRounds []data.MatchRound
	if err := db.Where(
		&data.MatchRound{GroupId: groupId},
	).Preload(
		"CommitJob",
	).Preload(
		"Matches.MenteeUser.Cohort.Cohort",
	).Preload(
		"Matches.MentorUser.Cohort.Cohort",
	).Find(matchRounds).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	apiMatchRounds := make([]api.MatchRound, 0, len(matchRounds))
	for _, matchRound := range matchRounds {
		state := getMatchRoundState(&matchRound)
		apiMatchRounds = append(
			apiMatchRounds, converters.ApiMatchRoundFromDataEntities(&matchRound, state))
	}
	return apiMatchRounds, nil
}

// Controller for DELETE match_round endpoint
func DeleteMatchRoundController(c *ctx.Context) errs.Error {
	matchRoundIdStr := c.GinContext.Param("matchRoundId")
	tempMatchRoundId, convErr := strconv.Atoi(matchRoundIdStr)
	matchRoundId := data.TMatchRoundID(tempMatchRoundId)

	if convErr != nil {
		return errs.NewRequestError(convErr.Error())
	}

	if err := handleDeleteMatchRound(c.Db, c.SessionData.UserId, matchRoundId); err != nil {
		return err
	}

	c.Result = "Success"
	return nil
}

func handleDeleteMatchRound(
	db *gorm.DB,
	adminId data.TUserID,
	matchRoundId data.TMatchRoundID,
) errs.Error {
	if err := checkIsAdminMatchRound(db, adminId, matchRoundId); err != nil {
		return err
	}

	return ctx.WithinTxRequestErr(db, func(db *gorm.DB) errs.Error {
		var matchRound data.MatchRound
		if err := db.Where(
			&data.MatchRound{Id: matchRoundId},
		).Preload("CommitJob").Find(&matchRound).Error; err != nil {
			return errs.NewDbError(err)
		}

		state := getMatchRoundState(&matchRound)
		if state != api.MATCH_ROUND_STATE_CREATED {
			return errs.NewRequestError(fmt.Sprintf("Cannot delete match round in %s state", state))
		}

		if err := db.Where(&data.MatchRound{Id: matchRoundId}).Delete(
			&data.MatchRound{},
		).Error; err != nil {
			return errs.NewDbError(err)
		}
		return nil
	})
}

func createMatchParameters(
	maxLowerYearsPerUpperYear uint,
	maxUpperYearsPerLowerYear uint,
	youngestUpperGradYear uint,
) data.MatchParameters {
	return data.MatchParameters(map[string]interface{}{
		"maxLowerYearsPerUpperYear": interface{}(maxLowerYearsPerUpperYear),
		"maxUpperYearsPerLowerYear": interface{}(maxUpperYearsPerLowerYear),
		"youngestUpperGradYear":     interface{}(youngestUpperGradYear),
	})
}

func getMatchRoundState(matchRound *data.MatchRound) api.MatchRoundState {
	if matchRound.CommitJob == nil {
		return api.MATCH_ROUND_STATE_CREATED
	} else if matchRound.CommitJob.Status == jobmine.STATUS_SUCCESS {
		return api.MATCH_ROUND_STATE_COMMITTED
	} else if matchRound.CommitJob.Status == jobmine.STATUS_FAILED {
		return api.MATCH_ROUND_STATE_FAILED
	} else {
		// Created or running counts as committings since it happens after the admin has "committed"
		// the matches. The only difference is whether the jobmine cron has picked it up yet or not.
		return api.MATCH_ROUND_STATE_COMMITTING
	}
}

func generateMatchRoundName(groupName string) string {
	now := time.Now()
	return fmt.Sprintf("%s: %s", groupName, now.Format("2006-01-02 at 15:04:05"))
}

// Assumes matches is non-empty
func createMatchRound(
	db *gorm.DB,
	groupId data.TGroupID,
	matches []recommendations.UserMatch,
	parameters data.MatchParameters,
) (*api.MatchRound, error) {
	var matchRound *data.MatchRound

	var group data.Group
	err := db.Where(&data.Group{GroupId: groupId}).Find(&group).Error
	if err != nil {
		return nil, err
	}

	err = ctx.WithinTx(db, func(db *gorm.DB) error {
		matchRound = &data.MatchRound{
			Name:            generateMatchRoundName(group.GroupName),
			GroupId:         groupId,
			MatchParameters: parameters,
			RunId:           nil,
		}

		matchRound.Matches = make([]data.MatchRoundMatch, 0, len(matches))
		for _, match := range matches {
			matchRound.Matches = append(matchRound.Matches, data.MatchRoundMatch{
				MenteeUserId: match.UserOneId,
				MentorUserId: match.UserTwoId,
				Score:        float32(match.Score),
			})
		}

		if err := db.Create(&matchRound).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Not sure how to avoid reloading the match round matches here, since Preload assumes you're
	// loading in new data.
	var matchRoundMatches []data.MatchRoundMatch

	err = db.Preload(
		"MenteeUser.Cohort.Cohort",
	).Preload(
		"MentorUser.Cohort.Cohort",
	).Where(
		&data.MatchRoundMatch{MatchRoundId: matchRound.Id},
	).Find(&matchRoundMatches).Error

	// NOTE: Will not get not found error since we are guaranteed to have at least one match
	// Look at error checking in `handleCreateMatchRound`.
	if err != nil {
		return nil, err
	}
	matchRound.Matches = matchRoundMatches

	apiMatchRound := converters.ApiMatchRoundFromDataEntities(
		matchRound, api.MATCH_ROUND_STATE_CREATED)
	return &apiMatchRound, nil
}

func checkIsAdminMatchRound(
	db *gorm.DB,
	adminId data.TUserID,
	matchRoundId data.TMatchRoundID,
) errs.Error {
	var matchRound data.MatchRound
	if err := db.Where(
		&data.MatchRound{Id: matchRoundId},
	).Find(&matchRound).Error; err != nil {
		return errs.NewDbError(err)
	}
	return checkIsAdmin(db, adminId, matchRound.GroupId)
}

func checkIsAdmin(db *gorm.DB, adminId data.TUserID, groupId data.TGroupID) errs.Error {
	err := db.Where(
		&data.ManagedGroup{AdministratorId: adminId, GroupId: groupId},
	).Find(&data.ManagedGroup{}).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errs.NewUnauthorizedError("You do not have rights to do this operation")
		} else {
			return errs.NewDbError(err)
		}
	}
	return nil
}

func checkUsersInGroup(db *gorm.DB, userIds []data.TUserID, groupId data.TGroupID) errs.Error {
	notInGroup := make([]data.TUserID, 0)
	var users []data.User
	if err := db.Where(
		"user_id IN (?)", userIds,
	).Preload("UserGroups").Find(&users).Error; err != nil {
		return errs.NewDbError(err)
	}

	for _, user := range users {
		hasGroup := false
		for _, group := range user.UserGroups {
			if groupId == group.GroupId {
				hasGroup = true
			}
		}

		if !hasGroup {
			notInGroup = append(notInGroup, user.UserId)
		}
	}

	if len(notInGroup) > 0 {
		return errs.NewRequestError(fmt.Sprintf("Users not in group: %v", notInGroup))
	}
	return nil
}
