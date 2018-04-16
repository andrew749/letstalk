package request_to_match

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
)

type GetCredentialOptionsResponse api.CredentialOptions

func GetCredentialOptionsController(c *ctx.Context) errs.Error {
	c.Result = api.GetCredentialOptions()
	return nil
}

// Credential CRUD

type AddUserCredentialResponse struct {
	CredentialId api.CredentialId `json:"credentialId"`
}

func AddUserCredentialController(c *ctx.Context) errs.Error {
	var credential api.CredentialPair

	if err := c.GinContext.BindJSON(&credential); err != nil {
		return errs.NewClientError("Unable to parse request %s", err)
	}

	credentialId, err := api.AddUserCredential(c.Db, c.SessionData.UserId, credential)
	if err != nil {
		return err
	}

	c.Result = AddUserCredentialResponse{*credentialId}
	return nil
}

type RemoveUserCredentialRequest struct {
	CredentialId api.CredentialId `json:"credentialId"`
}

func RemoveUserCredentialController(c *ctx.Context) errs.Error {
	var req RemoveUserCredentialRequest

	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewClientError("Unable to parse request %s", err)
	}

	if err := api.RemoveUserCredential(c.Db, c.SessionData.UserId, req.CredentialId); err != nil {
		return err
	}
	c.Result = "Removed"

	return nil
}

func GetUserCredentialsController(c *ctx.Context) errs.Error {
	userCredentials, err := api.GetUserCredentials(c.Db, c.SessionData.UserId)
	if err != nil {
		return err
	}
	c.Result = userCredentials
	return nil
}

// Credential Request CRUD

type AddUserCredentialRequestResponse struct {
	CredentialRequestId api.CredentialRequestId `json:"credentialRequestId"`
}

func AddUserCredentialRequestController(c *ctx.Context) errs.Error {
	var credential api.CredentialPair

	if err := c.GinContext.BindJSON(&credential); err != nil {
		return errs.NewClientError("Unable to parse request %s", err)
	}

	credentialId, err := api.AddUserCredentialRequest(c.Db, c.SessionData.UserId, credential)
	if err != nil {
		return err
	}

	c.Result = AddUserCredentialRequestResponse{*credentialId}
	return nil
}

type RemoveUserCredentialRequestRequest struct {
	CredentialRequestId api.CredentialRequestId `json:"credentialRequestId"`
}

func RemoveUserCredentialRequestController(c *ctx.Context) errs.Error {
	var req RemoveUserCredentialRequestRequest

	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewClientError("Unable to parse request %s", err)
	}

	if err := api.RemoveUserCredentialRequest(
		c.Db,
		c.SessionData.UserId,
		req.CredentialRequestId,
	); err != nil {
		return err
	}
	c.Result = "Removed"

	return nil
}

func GetUserCredentialRequestsController(c *ctx.Context) errs.Error {
	userCredentialRequests, err := api.GetUserCredentialRequests(c.Db, c.SessionData.UserId)
	if err != nil {
		return err
	}
	c.Result = userCredentialRequests
	return nil
}
