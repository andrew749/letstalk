package login

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/utility"
	"letstalk/server/data"
	"time"

	"github.com/romana/rlog"
)

type LoginRequestData struct {
	UserId   int    `json:"userId" binding:"required"`
	Password string `json:"password" binding:"required"`
	// optional token to associate with this session
	NotificationToken *string `json:"notificationToken"`
}

type LoginResponse struct {
	SessionId  string    `json:"sessionId"`
	ExpiryDate time.Time `json:"expiry"`
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

	var authData data.AuthenticationData
	if c.Db.Where("user_id = ?", req.UserId).First(&authData).RecordNotFound() {
		return errs.NewClientError("Couldn't find an account")
	}

	// check if the password is correct
	if !utility.CheckPasswordHash(req.Password, authData.PasswordHash) {
		return errs.NewClientError("Bad password")
	}

	rlog.Debug("Successfully Checked Password")

	session, err := (*sm).CreateNewSessionForUserId(req.UserId, req.NotificationToken)
	if err != nil {
		return errs.NewClientError("%s", err)
	}
	c.Result = LoginResponse{*session.SessionId, session.ExpiryDate}

	return nil
}
