package query

import (
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

// GetMatchingByUserIds returns details on a matching between two users.
func GetMatchingByUserIds(db *gorm.DB, firstUser int, secondUser int) (*data.Matching, error) {
	var matching *data.Matching = nil
	err := db.
		Where(&data.Matching{Mentor: firstUser, Mentee: secondUser}).
		Or(&data.Matching{Mentor: secondUser, Mentee: firstUser}).
		Preload("MenteeUser").
		Preload("MentorUser").
		First(matching).Error
	if err != nil {
		return nil, err
	}
	return matching, nil
}

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
