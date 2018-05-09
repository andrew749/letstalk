package query

import (
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

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
