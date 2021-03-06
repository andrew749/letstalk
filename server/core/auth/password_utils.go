package auth

import (
	"letstalk/server/core/errs"
	"letstalk/server/core/utility"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

type HashingError struct {
	errs.IError
}

// ChangeUserPassword update the specified user password
func ChangeUserPassword(db *gorm.DB, userId data.TUserID, newPassword string) error {
	var err error
	var hashedPassword string
	if hashedPassword, err = utility.HashPassword(newPassword); err != nil {
		return &HashingError{errs.NewInternalError(err.Error())}
	}

	authData := data.AuthenticationData{
		UserId:       userId,
		PasswordHash: hashedPassword,
	}

	if err := db.Save(&authData).Error; err != nil {
		return errs.NewDbError(err)
	}

	rlog.Infof("Changed user password for user %d to %s", userId, hashedPassword)

	return nil
}
