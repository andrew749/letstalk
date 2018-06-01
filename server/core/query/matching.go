package query

import (
	"letstalk/server/core/api"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

// GetMatchingByUserIds returns details on a matching between two users.
func GetMatchingByUserIds(db *gorm.DB, firstUser int, secondUser int) (*data.Matching, error) {
	matchings := make([]data.Matching, 0)
	err := db.
		Where(&data.Matching{Mentor: firstUser, Mentee: secondUser}).
		Or(&data.Matching{Mentor: secondUser, Mentee: firstUser}).
		Preload("MenteeUser").
		Preload("MentorUser").
		First(&matchings).Error
	if err != nil {
		return nil, err
	}
	if len(matchings) == 0 {
		return nil, nil
	}
	return &matchings[0], nil
}

func GetMenteesByMentorId(
	db *gorm.DB,
	mentorId int,
	flag api.MatchingInfoFlag,
) ([]data.Matching, error) {
	matchings := make([]data.Matching, 0)
	req := db.Where(&data.Matching{Mentor: mentorId}).Preload("MenteeUser")

	if flag&api.MATCHING_INFO_FLAG_AUTH_DATA != 0 {
		req = req.Preload("MenteeUser.ExternalAuthData")
	}
	if flag&api.MATCHING_INFO_FLAG_COHORT != 0 {
		req = req.Preload("MenteeUser.Cohort.Cohort")
	}

	if err := req.Find(&matchings).Error; err != nil {
		return nil, err
	}
	return matchings, nil
}

func GetMentorsByMenteeId(
	db *gorm.DB,
	menteeId int,
	flag api.MatchingInfoFlag,
) ([]data.Matching, error) {
	matchings := make([]data.Matching, 0)
	req := db.Where(&data.Matching{Mentee: menteeId}).Preload("MentorUser")

	if flag&api.MATCHING_INFO_FLAG_AUTH_DATA != 0 {
		req = req.Preload("MentorUser.ExternalAuthData")
	}
	if flag&api.MATCHING_INFO_FLAG_COHORT != 0 {
		req = req.Preload("MentorUser.Cohort.Cohort")
	}

	if err := req.Find(&matchings).Error; err != nil {
		return nil, err
	}
	return matchings, nil
}
