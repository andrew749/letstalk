package login

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/utility"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

func invalidPassError() errs.Error {
	return errs.NewClientError("Invalid Password. Try again.")
}

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
		return errs.NewClientError("Bad login request data %s", err)
	}

	var userModel data.User
	if err := c.Db.Where("email = ?", req.Email).Preload("AuthData").First(&userModel).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return invalidPassError()
		} else {
			return errs.NewDbError(err)
		}
	}

	// check if the password is correct
	if !utility.CheckPasswordHash(req.Password, userModel.AuthData.PasswordHash) {
		return invalidPassError()
	}

	rlog.Debug("Successfully Checked Password")

	session, err := (*sm).CreateNewSessionForUserId(userModel.UserId, req.NotificationToken)
	if err != nil {
		return errs.NewClientError("%s", err)
	}
	c.Result = api.LoginResponse{
		SessionId:  *session.SessionId,
		ExpiryDate: session.ExpiryDate,
	}

	return nil
}
