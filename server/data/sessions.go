package data

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Session struct {
	SessionId  string    `gorm:"not null;primary_key;size:190"`
	User       User      `gorm:"foreignkey:UserId"`
	UserId     TUserID   `gorm:"not null"`
	ExpiryDate time.Time `gorm:"not null"`
}

func DeleteSession(db *gorm.DB, sessionId string) error {
	if err := db.
		Where("session_id = ?", sessionId).
		Delete(Session{}).
		Error; err != nil {
		return errors.Wrap(err, "Unable to delete session")
	}
	return nil
}
