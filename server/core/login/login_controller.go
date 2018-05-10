package login

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/utility"
	"letstalk/server/data"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

type LoginRequestData struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	// optional token to associate with this session
	NotificationToken *string `json:"notificationToken"`
}

type LoginResponse struct {
	SessionId  string    `json:"sessionId"`
	ExpiryDate time.Time `json:"expiry"`
}

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

	var req LoginRequestData
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
	c.Result = LoginResponse{*session.SessionId, session.ExpiryDate}

	return nil
}
