package user

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func LogoutHandler(c *ctx.Context) errs.Error {
	if c.SessionData.SessionId == nil {
		return errs.NewInternalError("Bad session token.")
	}

	// TODO: Refactor this to not live in the controller
	// remove the session from list of active session
	err := c.Db.Where("session_id = ?", c.SessionData.SessionId).Delete(data.Session{}).Error
	if err != nil {
		return errs.NewDbError(err)
	}

	err = c.Db.Where("session_id = ?", c.SessionData.SessionId).Delete(data.NotificationToken{}).Error
	if err != nil {
		return errs.NewDbError(err)
	}

	var session data.Session
	if err = c.Db.
		Where("session_id = ?", c.SessionData.SessionId).
		Preload("NotificationToken").
		First(&session).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return errs.NewDbError(err)
		}

		// try to delete all devices for this session
		if err = c.Db.Delete(&data.UserDevice{
			NotificationToken: session.NotificationToken.Token,
		}).Error; err != nil {
			return errs.NewDbError(errors.Wrap(err, "Cannot delete notification token"))
		}
	}

	c.Result = "Ok"

	return nil
}
