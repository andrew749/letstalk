package request_to_match

import (
	"letstalk/server/core/query"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
)

type GetCredentialOptionsResponse query.CredentialOptions

func GetCredentialOptionsController(c *ctx.Context) errs.Error {
	c.Result = query.GetCredentialOptions()
	return nil
}

// Credential CRUD

type AddUserCredentialResponse struct {
	CredentialId query.CredentialId `json:"credentialId"`
}

func AddUserCredentialController(c *ctx.Context) errs.Error {
	var credential query.CredentialPair

	if err := c.GinContext.BindJSON(&credential); err != nil {
		return errs.NewClientError("Unable to parse request %s", err)
	}

	s := query.UserCredentialStrategy{c.Db, c.SessionData.UserId}
	credentialId, err := query.AddCredentialWithStrategy(s, credential)
	if err != nil {
		return err
	}

	c.Result = AddUserCredentialResponse{*credentialId}
	return nil
}

type RemoveUserCredentialRequest struct {
	CredentialId query.CredentialId `json:"credentialId"`
}

func RemoveUserCredentialController(c *ctx.Context) errs.Error {
	var req RemoveUserCredentialRequest

	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewClientError("Unable to parse request %s", err)
	}

	s := query.UserCredentialStrategy{c.Db, c.SessionData.UserId}
	if err := s.DeleteCredentialForUser(req.CredentialId); err != nil {
		return err
	}
	c.Result = "Removed"

	return nil
}

func GetUserCredentialsController(c *ctx.Context) errs.Error {
	s := query.UserCredentialStrategy{c.Db, c.SessionData.UserId}
	userCredentials, err := query.GetCredentialsWithStrategy(s)
	if err != nil {
		return err
	}
	c.Result = userCredentials
	return nil
}

// Credential Request CRUD

type AddUserCredentialRequestResponse struct {
	CredentialRequestId query.CredentialRequestId `json:"credentialRequestId"`
}

func AddUserCredentialRequestController(c *ctx.Context) errs.Error {
	var credential query.CredentialPair

	if err := c.GinContext.BindJSON(&credential); err != nil {
		return errs.NewClientError("Unable to parse request %s", err)
	}

	s := query.UserCredentialRequestStrategy{c.Db, c.SessionData.UserId}
	credentialId, err := query.AddCredentialWithStrategy(s, credential)
	if err != nil {
		return err
	}

	var isAdded bool
	credentialRequestId := query.CredentialRequestId(*credentialId)
	isAdded, err = query.ResolveRequestToMatch(c.Db, c.SessionData.UserId, credentialRequestId)
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
	CredentialRequestId query.CredentialRequestId `json:"credentialRequestId"`
}

func RemoveUserCredentialRequestController(c *ctx.Context) errs.Error {
	var req RemoveUserCredentialRequestRequest

	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewClientError("Unable to parse request %s", err)
	}

	s := query.UserCredentialRequestStrategy{c.Db, c.SessionData.UserId}
	if err := s.DeleteCredentialForUser(query.CredentialId(req.CredentialRequestId)); err != nil {
		return err
	}
	c.Result = "Removed"

	return nil
}

func GetUserCredentialRequestsController(c *ctx.Context) errs.Error {
	s := query.UserCredentialRequestStrategy{c.Db, c.SessionData.UserId}
	credentials, err := query.GetCredentialsWithStrategy(s)
	if err != nil {
		return err
	}

	credentialRequests := make([]query.CredentialRequestWithId, len(credentials))
	for i, credential := range credentials {
		credentialRequests[i] = query.CredentialRequestWithId{
			Credential:          credential.Credential,
			CredentialRequestId: query.CredentialRequestId(credential.CredentialId),
		}
	}

	c.Result = credentialRequests
	return nil
}
