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
	UserId   int    `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
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
	rlog.Debug("Handling route")
	// create new session
	sm := c.SessionManager

	var req LoginRequestData
	err := c.GinContext.BindJSON(&req)
	if err != nil {
		return errs.NewClientError("Bad login request data %s", err)
	}

	authenticationEntry, err := data.AuthenticationDataObjs.Select().Where(data.AuthenticationDataObjs.FilterUserId("=", req.UserId)).List(c.Db)
	if len(authenticationEntry) == 0 {
		return errs.NewClientError("Couldn't find an account")
	}

	// check if the password is correct
	if !utility.CheckPasswordHash(req.Password, authenticationEntry[0].PasswordHash) {
		return errs.NewClientError("Bad password")
	}
	rlog.Debug("Successfully Checked Password")

	// if all preconditions pass, then create a new session
	session, errSession := (*sm).CreateNewSessionForUserId(req.UserId)

	if errSession != nil {
		return errSession
	}
	rlog.Debug("Successfully Created Session")

	c.Result = LoginResponse{*session.SessionId, session.ExpiryDate}

	return nil
}
