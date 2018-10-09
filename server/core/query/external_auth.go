package query

import (
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// GetExternalAuth Returns a matching external auth.
func GetExternalAuthRecord(db *gorm.DB, userID data.TUserID) (*data.ExternalAuthData, error) {
	var auth data.ExternalAuthData
	if err := db.Where(&data.ExternalAuthData{UserId: userID}).First(&auth).Error; err != nil {
		return nil, errors.Errorf("Unable to get User with id")
	}
	return &auth, nil
}

func GetExternalAuthRecordByFBIDNoAssert(db *gorm.DB, fbUserID *string) (*data.ExternalAuthData, error) {
	var externalAuthRecord data.ExternalAuthData
	if err := db.
		Where(&data.ExternalAuthData{FbUserId: fbUserID}).
		First(&externalAuthRecord).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &externalAuthRecord, nil
}
