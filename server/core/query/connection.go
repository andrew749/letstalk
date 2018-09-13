package query

import (
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

// GetConnectionDetails returns details on an existing directed connection between two users.
func GetConnectionDetails(db *gorm.DB, requestingUser data.TUserID, connectedUser data.TUserID) (*data.Connection, error) {
	var connection data.Connection
	q := db.
	Where(&data.Connection{UserOneId: requestingUser, UserTwoId: connectedUser}).
		Where("deleted_at IS NULL").
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
func GetConnectionDetailsUndirected(db *gorm.DB, firstUser data.TUserID, secondUser data.TUserID) (*data.Connection, error) {
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
	connections := make ([]data.Connection, 0)
	q := db.Where(&data.Connection{UserOneId: userId}).
		Or(&data.Connection{UserTwoId: userId}).
		Where("deleted_at IS NULL").
		Preload("Intent").
		Preload("Mentorship").
		Preload("UserOne").
		Preload("UserTwo").
		Find(&connections)
	if q.RecordNotFound() {
		return []data.Connection{}, nil
	}
	if q.Error != nil {
		return nil, q.Error
	}
	return connections, nil
}
