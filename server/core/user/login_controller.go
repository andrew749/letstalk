package user

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/utility"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

/**
 * Method to get called in login route
 * ```
 *  {"user_id": string, "password": string}
 * ```
 */
func LoginUser(c *ctx.Context) errs.Error {
	// create new session
	sm := c.SessionManager

	var req api.LoginRequestData
	err := c.GinContext.BindJSON(&req)
	if err != nil {
		return errs.NewRequestError("Bad login request data %s", err)
	}

	var userModel data.User
	if err := c.Db.Where("email = ?", req.Email).Preload("AuthData").First(&userModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return errs.InvalidPassError()
		} else {
			return errs.NewDbError(err)
		}
	}

	// check if the password is correct
	if !utility.CheckPasswordHash(req.Password, userModel.AuthData.PasswordHash) {
		return errs.InvalidPassError()
	}

	rlog.Debug("Successfully Checked Password")

	session, err := (*sm).CreateNewSessionForUserId(userModel.UserId)
	if err != nil {
		return errs.NewRequestError("%s", err)
	}

	notificationToken := req.NotificationToken
	// add device token to db
	if notificationToken != nil {
		if err := data.AddExpoDeviceTokenforUser(c.Db, userModel.UserId, *notificationToken); err != nil {
			return errs.NewInternalError("Unable to register device in db.")
		}
	}

	c.Result = api.LoginResponse{
		SessionId:  *session.SessionId,
		ExpiryDate: session.ExpiryDate,
	}

	return nil
}
