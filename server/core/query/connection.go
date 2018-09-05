package query

import (
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

// GetConnectionDetails returns details on a connection between two users.
func GetConnectionDetails(db *gorm.DB, firstUser data.TUserID, secondUser data.TUserID) (*data.Connection, error) {
	var connection data.Connection
	q := db.
		Where(&data.Connection{UserOneId: firstUser, UserTwoId: secondUser}).
		Or(&data.Connection{UserOneId: secondUser, UserTwoId: firstUser}).
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
