package query

import (
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

func GetMenteesByMentorId(db *gorm.DB, mentorId int) ([]data.Matching, error) {
	matchings := make([]data.Matching, 0)
	err := db.Where(
		&data.Matching{Mentor: mentorId},
	).Preload("MenteeUser").Preload("MenteeUser.ExternalAuthData").Find(&matchings).Error
	if err != nil {
		return nil, err
	}
	return matchings, nil
}

func GetMentorsByMenteeId(db *gorm.DB, menteeId int) ([]data.Matching, error) {
	matchings := make([]data.Matching, 0)
	err := db.Where(
		&data.Matching{Mentee: menteeId},
	).Preload("MentorUser").Preload("MentorUser.ExternalAuthData").Find(&matchings).Error
	if err != nil {
		return nil, err
	}
	return matchings, nil
}
