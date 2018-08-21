package controller

import (
	"math/rand"

	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/search"
	"letstalk/server/data"
)

type AddSimpleTraitToESRequest struct {
	Name string
}

func AddSimpleTraitToES(c *ctx.Context) errs.Error {
	var req AddSimpleTraitToESRequest
	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	trait := search.SimpleTrait{
		Id:              data.TSimpleTraitID(rand.Int()),
		Name:            req.Name,
		Type:            data.SIMPLE_TRAIT_TYPE_NONE,
		IsSensitive:     false,
		IsUserGenerated: true,
	}

	// TODO: Error handling
	search.InsertSimpleTrait(c.Es, trait)
	search.PrintAllSimpleTraits(c.Es)
	return nil
}
