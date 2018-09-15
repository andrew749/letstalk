package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"

	"github.com/romana/rlog"
)

func SimpleTraitAutocompleteController(c *ctx.Context) errs.Error {
	var req api.AutocompleteRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	searchClient := c.SearchClientWithContext()
	rlog.Debugf("[ANDREW] %#v", req)
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

	apiRoles := make([]api.Role, len(roles))
	for i, role := range roles {
		apiRoles[i] = api.Role{
			role.Id,
			role.Name,
		}
	}

	c.Result = apiRoles
	return nil
}

func OrganizationAutocompleteController(c *ctx.Context) errs.Error {
	var req api.AutocompleteRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	searchClient := c.SearchClientWithContext()

	organizations, err := searchClient.CompletionSuggestionOrganizations(req.Prefix, req.Size)
	if err != nil {
		return errs.NewEsError(err)
	}

	apiOrganizations := make([]api.Organization, len(organizations))
	for i, organization := range organizations {
		apiOrganizations[i] = api.Organization{
			organization.Id,
			organization.Name,
			organization.Type,
		}
	}

	c.Result = apiOrganizations
	return nil
}

func MultiTraitAutocompleteController(c *ctx.Context) errs.Error {
	var req api.AutocompleteRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	searchClient := c.SearchClientWithContext()
	traits, err := searchClient.QueryMultiTraitsByName(req.Prefix, req.Size)
	if err != nil {
		return errs.NewEsError(err)
	}

	c.Result = traits
	return nil
}
