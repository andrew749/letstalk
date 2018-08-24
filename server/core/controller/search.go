package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
)

func SimpleTraitAutocompleteController(c *ctx.Context) errs.Error {
	var req api.AutocompleteRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	searchClient := c.SearchClientWithContext()
	traits, err := searchClient.CompletionSuggestionSimpleTraits(req.Prefix, req.Size)
	if err != nil {
		return errs.NewEsError(err)
	}

	apiTraits := make([]api.SimpleTrait, len(traits))
	for i, trait := range traits {
		apiTraits[i] = api.SimpleTrait{
			trait.Id,
			trait.Name,
			trait.Type,
			trait.IsSensitive,
		}
	}

	c.Result = apiTraits
	return nil
}

func RoleAutocompleteController(c *ctx.Context) errs.Error {
	var req api.AutocompleteRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	searchClient := c.SearchClientWithContext()
	roles, err := searchClient.CompletionSuggestionRoles(req.Prefix, req.Size)
	if err != nil {
		return errs.NewEsError(err)
	}

	apiroles := make([]api.Role, len(roles))
	for i, role := range roles {
		apiroles[i] = api.Role{
			role.Id,
			role.Name,
		}
	}

	c.Result = apiroles
	return nil
}
