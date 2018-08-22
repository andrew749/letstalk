package controller

import (
	"math/rand"

	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/search"
	"letstalk/server/data"

	"github.com/romana/rlog"
)

type AddSimpleTraitToESRequest struct {
	Name string
}

type SimpleTraitAutocompleteRequest struct {
	Prefix string `json:"prefix" binding:"required"`
	Size   int    `json:"size" binding:"required"`
}

func AddSimpleTraitToES(c *ctx.Context) errs.Error {
	var req AddSimpleTraitToESRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	trait := search.SimpleTrait{
		Id:              data.TSimpleTraitID(rand.Int()),
		Name:            req.Name,
		Type:            data.SIMPLE_TRAIT_TYPE_UNDETERMINED,
		IsSensitive:     false,
		IsUserGenerated: true,
	}

	err := c.SearchClient.IndexSimpleTrait(trait)
	if err != nil {
		rlog.Error(err)
		return errs.NewDbError(err)
	}
	return nil
}

func SimpleTraitAutocompleteController(c *ctx.Context) errs.Error {
	var req SimpleTraitAutocompleteRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	traits, err := c.SearchClient.CompletionSuggestionSimpleTraits(req.Prefix, req.Size)
	if err != nil {
		// TODO: New error type
		return errs.NewDbError(err)
	}

	c.Result = traits
	return nil
}
