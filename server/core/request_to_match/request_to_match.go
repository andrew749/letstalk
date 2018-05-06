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

	s := api.UserCredentialStrategy{c.Db, c.SessionData.UserId}
	credentialId, err := api.AddCredentialWithStrategy(s, credential)
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

	s := api.UserCredentialStrategy{c.Db, c.SessionData.UserId}
	if err := s.DeleteCredentialForUser(req.CredentialId); err != nil {
		return err
	}
	c.Result = "Removed"

	return nil
}

func GetUserCredentialsController(c *ctx.Context) errs.Error {
	s := api.UserCredentialStrategy{c.Db, c.SessionData.UserId}
	userCredentials, err := api.GetCredentialsWithStrategy(s)
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

	s := api.UserCredentialRequestStrategy{c.Db, c.SessionData.UserId}
	credentialId, err := api.AddCredentialWithStrategy(s, credential)
	if err != nil {
		return err
	}

	var isAdded bool
	credentialRequestId := api.CredentialRequestId(*credentialId)
	isAdded, err = api.ResolveRequestToMatch(c.Db, c.SessionData.UserId, credentialRequestId)
	if err != nil {
		return err
	}

	if isAdded {
		return errs.NewClientError("Found a match right away", err)
	}

	c.Result = AddUserCredentialRequestResponse{credentialRequestId}
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

	s := api.UserCredentialRequestStrategy{c.Db, c.SessionData.UserId}
	if err := s.DeleteCredentialForUser(api.CredentialId(req.CredentialRequestId)); err != nil {
		return err
	}
	c.Result = "Removed"

	return nil
}

func GetUserCredentialRequestsController(c *ctx.Context) errs.Error {
	s := api.UserCredentialRequestStrategy{c.Db, c.SessionData.UserId}
	credentials, err := api.GetCredentialsWithStrategy(s)
	if err != nil {
		return err
	}

	credentialRequests := make([]api.CredentialRequestWithId, len(credentials))
	for i, credential := range credentials {
		credentialRequests[i] = api.CredentialRequestWithId{
			Credential:          credential.Credential,
			CredentialRequestId: api.CredentialRequestId(credential.CredentialId),
		}
	}

	c.Result = credentialRequests
	return nil
}
