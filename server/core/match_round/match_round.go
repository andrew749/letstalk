package match_round

import (
	"fmt"
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"
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

	err := handleCreateMatchRound(
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

	return nil
}

func handleCreateMatchRound(
	db *gorm.DB,
	maxLowerYearsPerUpperYear uint,
	maxUpperYearsPerLowerYear uint,
	youngestUpperGradYear uint,
	groupId data.TGroupID,
	userIds []data.TUserID,
) errs.Error {
	if userIds == nil {
		return errs.NewRequestError("Expected non-nil user ids")
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
		return errs.NewRequestError(errStr)
	}

	if len(matches) == 0 {
		return errs.NewRequestError("Parameters result in no matches")
	}

	parameters := createMatchParameters(
		maxLowerYearsPerUpperYear,
		maxUpperYearsPerLowerYear,
		youngestUpperGradYear,
	)

	_, err = createMatchRound(
		db,
		groupId,
		matches,
		parameters,
	)
	// TODO(match-api): Make this return the users matched to each other (first/last names and programs)

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

func createMatchRound(
	db *gorm.DB,
	groupId data.TGroupID,
	matches []recommendations.UserMatch,
	parameters data.MatchParameters,
) (*data.MatchRound, error) {
	var matchRound *data.MatchRound

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

		roundMatches := make([]data.MatchRoundMatch, 0, len(matches))
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

	return matchRound, nil
}
