package query

import (
	"letstalk/server/data"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// First parameter should be a db transaction.
func GenerateNewVerifyEmailId(tx *gorm.DB, userId data.TUserID, emailAddr string) (*data.VerifyEmailId, error) {
	var id = uuid.New()
	verifyEmailData := data.VerifyEmailId{
		Id:             id.String(),
		UserId:         userId,
		Email:          emailAddr,
		IsActive:       true,
		IsUsed:         false,
		ExpirationDate: time.Now().AddDate(0, 0, 1), // Verification email valid for 24 hours.
	}
	// Set all existing VerifyEmailId entries for this user to inactive.
	err := tx.Model(&data.VerifyEmailId{}).
		Where(&data.VerifyEmailId{UserId: userId}).
		Update("is_active", false).
		Error
	if err != nil {
		return nil, err
	}
	// Insert the new VerifyEmailId entry.
	if err := tx.Save(&verifyEmailData).Error; err != nil {
		return nil, err
	}
	return &verifyEmailData, nil
}
