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

type AddUserCredentialResponse struct {
	CredentialId uint `json:"credentialId"`
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
	CredentialId uint `json:"credentialId"`
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
