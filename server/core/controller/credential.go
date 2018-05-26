package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
)

const RESOLVE_WAIT_TIME = 5000 // ms

func GetAllCredentialsController(c *ctx.Context) errs.Error {
	credentials, err := query.GetAllCredentials(c.Db)
	if err != nil {
		return err
	}
	c.Result = credentials
	return nil
}

func AddUserCredentialRequestController(c *ctx.Context) errs.Error {
	var req api.AddUserCredentialRequestRequest

	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewClientError(err.Error())
	}

	if err := query.AddUserCredentialRequest(
		c.Db,
		c.SessionData.UserId,
		req.CredentialId,
	); err != nil {
		return err
	}

	go query.ResolveRequestToMatchWithDelay(
		c,
		query.RESOLVE_TYPE_ASKER,
		req.CredentialId,
		RESOLVE_WAIT_TIME,
	)

	return nil
}

func RemoveUserCredentialRequestController(c *ctx.Context) errs.Error {
	var req api.RemoveUserCredentialRequestRequest

	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewClientError(err.Error())
	}

	if err := query.RemoveUserCredentialRequest(
		c.Db,
		c.SessionData.UserId,
		req.CredentialId,
	); err != nil {
		return err
	}

	return nil
}

func GetUserCredentialRequestsController(c *ctx.Context) errs.Error {
	credentials, err := query.GetUserCredentialRequests(c.Db, c.SessionData.UserId)
	if err != nil {
		return err
	}

	c.Result = credentials
	return nil
}

func AddUserCredentialController(c *ctx.Context) errs.Error {
	var req api.AddUserCredentialRequest

	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewClientError(err.Error())
	}

	credentialId, err := query.AddUserCredential(c.Db, c.SessionData.UserId, req.Name)
	if err != nil {
		return err
	}

	go query.ResolveRequestToMatchWithDelay(
		c,
		query.RESOLVE_TYPE_ANSWERER,
		*credentialId,
		RESOLVE_WAIT_TIME,
	)

	c.Result = api.AddUserCredentialResponse{CredentialId: *credentialId}
	return nil
}

func RemoveUserCredentialController(c *ctx.Context) errs.Error {
	var req api.RemoveUserCredentialRequest

	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewClientError(err.Error())
	}

	if err := query.RemoveUserCredential(c.Db, c.SessionData.UserId, req.CredentialId); err != nil {
		return err
	}

	return nil
}

func GetUserCredentialsController(c *ctx.Context) errs.Error {
	credentials, err := query.GetUserCredentials(c.Db, c.SessionData.UserId)
	if err != nil {
		return err
	}

	c.Result = credentials
	return nil
}
