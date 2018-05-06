package query

import (
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func GetMenteesByMentorId(db *gorm.DB, mentorId int) ([]data.Matching, error) {
	matchings := make([]data.Matching, 0)
	if err := db.Where(&data.Matching{Mentor: mentorId}).Preload("MenteeUser").Find(&matchings).Error; err != nil {
		return nil, err
	}
	return matchings, nil
}

func GetMentorsByMenteeId(db *gorm.DB, menteeId int) ([]data.Matching, error) {
	matchings := make([]data.Matching, 0)
	if err := db.Where(&data.Matching{Mentee: menteeId}).Preload("MentorUser").Find(&matchings).Error; err != nil {
		return nil, err
	}
	return matchings, nil
}
