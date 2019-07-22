package match_round

import (
	"fmt"
	"letstalk/server/core/api"
	"letstalk/server/core/converters"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"
	"letstalk/server/jobmine_jobs/match_round_commit_job"
	"letstalk/server/recommendations"
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
		request.Parameters.MaxLowerYearsPerUpperYear,
		request.Parameters.MaxUpperYearsPerLowerYear,
		request.Parameters.YoungestUpperGradYear,
		request.GroupId,
		request.UserIds,
	)
	if err != nil {
		return err
	}

	c.Result = matchRound
	return nil
}

func handleCreateMatchRound(
	db *gorm.DB,
	maxLowerYearsPerUpperYear uint,
	maxUpperYearsPerLowerYear uint,
	youngestUpperGradYear uint,
	groupId data.TGroupID,
	userIds []data.TUserID,
) (*api.MatchRound, errs.Error) {
	if userIds == nil {
		return nil, errs.NewRequestError("Expected non-nil user ids")
	}

	// TODO(match-api): Check if users are in given group

	strat := recommendations.MentorMenteeStrat(
		maxLowerYearsPerUpperYear,
		maxUpperYearsPerLowerYear,
		youngestUpperGradYear,
	)

	fetcherOptions := recommendations.UserFetcherOptions{UserIds: userIds}
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
		maxLowerYearsPerUpperYear,
		maxUpperYearsPerLowerYear,
		youngestUpperGradYear,
	)

	matchRound, err := createMatchRound(
		db,
		groupId,
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
	userId data.TUserID,
	matchRoundId data.TMatchRoundID,
) errs.Error {
	// TODO(match-api): Check that user is authorized to commit this round
	err := ctx.WithinTx(db, func(db *gorm.DB) error {
		runId, err := match_round_commit_job.CreateCommitJob(db, matchRoundId)
		if err != nil {
			return err
		}

		var matchRound data.MatchRound
		if err := db.Where(&data.MatchRound{Id: matchRoundId}).Find(&matchRound).Error; err != nil {
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

// TODO(match-api): Use the correct group id/name
func generateMatchRoundName(groupId data.TGroupID) string {
	groupName := string(groupId)
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
	var roundMatches []data.MatchRoundMatch

	err := ctx.WithinTx(db, func(db *gorm.DB) error {
		matchRound = &data.MatchRound{
			Name:            generateMatchRoundName(groupId),
			GroupId:         groupId,
			MatchParameters: parameters,
			RunId:           nil,
		}

		if err := db.Create(matchRound).Error; err != nil {
			return err
		}

		roundMatches = make([]data.MatchRoundMatch, 0, len(matches))
		for _, match := range matches {
			roundMatches = append(roundMatches, data.MatchRoundMatch{
				MatchRoundId: matchRound.Id,
				MenteeUserId: match.UserOneId,
				MentorUserId: match.UserTwoId,
				Score:        float32(match.Score),
			})
		}

		if err := db.Create(roundMatches).Error; err != nil {
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

	apiMatchRound := converters.ApiMatchRoundFromDataEntities(matchRound, matchRoundMatches)
	return &apiMatchRound, nil
}
