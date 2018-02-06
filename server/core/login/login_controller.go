package login

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/utility"
	"letstalk/server/data"
	"time"
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
	// create new session
	sm := c.SessionManager

	var req LoginRequestData
	err := c.GinContext.BindJSON(req)
	if err != nil {
		return errs.NewClientError("Bad login request data")
	}

	authenticationEntry, err := data.AuthenticationDataObjs.Select().Where(data.AuthenticationDataObjs.FilterUserId("=", req.UserId)).List(c.Db)

	// check if the password is correct
	if !utility.CheckPasswordHash(req.Password, authenticationEntry[0].PasswordHash) {
		return errs.NewClientError("Bad password")
	}

	// if all preconditions pass, then create a new session
	session, errSession := (*sm).CreateNewSessionForUserId(req.UserId)

	if errSession != nil {
		return errSession
	}

	c.Result = LoginResponse{*session.SessionId, session.ExpiryDate}

	return nil
}
