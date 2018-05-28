package query

import (
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// GetExternalAuth Returns a matching external auth.
func GetExternalAuthRecord(db *gorm.DB, userID int) (*data.ExternalAuthData, error) {
	var auth data.ExternalAuthData
	if err := db.Where(&data.ExternalAuthData{UserId: userID}).First(&auth).Error; err != nil {
		return nil, errors.Errorf("Unable to get User with id")
	}
	return &auth, nil
}
