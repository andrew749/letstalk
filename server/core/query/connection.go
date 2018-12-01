package query

import (
	"time"

	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

// GetConnectionDetails returns details on an existing directed connection between two users.
func GetConnectionDetails(
	db *gorm.DB,
	requestingUser data.TUserID,
	connectedUser data.TUserID,
) (*data.Connection, error) {
	var connection data.Connection
	q := db.
		Where(&data.Connection{UserOneId: requestingUser, UserTwoId: connectedUser}).
		Preload("Intent").
		Preload("Mentorship").
		First(&connection)
	if q.RecordNotFound() {
		return nil, nil
	}
	if q.Error != nil {
		return nil, q.Error
	}
	return &connection, nil
}

// GetConnectionDetailsUndirected returns details on an existing connection between two users in either direction.
func GetConnectionDetailsUndirected(
	db *gorm.DB,
	firstUser data.TUserID,
	secondUser data.TUserID,
) (*data.Connection, error) {
	if connection, err := GetConnectionDetails(db, firstUser, secondUser); err != nil {
		return nil, err
	} else if connection != nil {
		return connection, nil
	}
	if connection, err := GetConnectionDetails(db, secondUser, firstUser); err != nil {
		return nil, err
	} else {
		return connection, nil
	}
}

// GetAllConnections returns all of a user's connections.
func GetAllConnections(db *gorm.DB, userId data.TUserID) ([]data.Connection, error) {
	connections := make([]data.Connection, 0)
	q := db.Where(&data.Connection{UserOneId: userId}).
		Or(&data.Connection{UserTwoId: userId}).
		Where("deleted_at IS NULL").
		Preload("Intent").
		Preload("Mentorship").
		Preload("UserOne.Cohort.Cohort").
		Preload("UserTwo.Cohort.Cohort").
		Preload("UserOne.ExternalAuthData").
		Preload("UserTwo.ExternalAuthData").
		Find(&connections)
	if q.RecordNotFound() {
		return []data.Connection{}, nil
	}
	if q.Error != nil {
		return nil, q.Error
	}
	return connections, nil
}

const (
	leftJoinStr = "LEFT JOIN mentorships ON mentorships.connection_id = connections.connection_id"
	whereStr    = "mentorships.connection_id IS NOT NULL"
)

func GetMentorshipConnectionsByDate(
	db *gorm.DB,
	startDate *time.Time,
	endDate *time.Time,
) ([]data.Connection, error) {
	var connections []data.Connection
	query := db.Model(&data.Connection{}).Joins(leftJoinStr).Where(
		whereStr,
	).Preload("Mentorship").Preload("UserOne").Preload("UserTwo")
	if startDate != nil {
		query = query.Where("mentorships.created_at >= ?", *startDate)
	}
	if endDate != nil {
		query = query.Where("mentorships.created_at <= ?", *endDate)
	}

	if err := query.Find(&connections).Error; err != nil {
		return nil, err
	}
	return connections, nil
}
