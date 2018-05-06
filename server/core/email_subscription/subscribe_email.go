package email_subscription

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/romana/rlog"
)

type SubscriptionRequest struct {
	ClassYear    int    `json:"classYear" binding:"required"`
	ProgramName  string `json:"programName" binding:"required"`
	EmailAddress string `json:"emailAddress" binding:"required"`
	FirstName    string `json:"firstName" binding:"required"`
	LastName     string `json:"lastName" binding:"required"`
}

type SubscriptionResponse struct {
	Status string `json:"status"`
}

func AddSubscription(ctx *ctx.Context) errs.Error {
	var request SubscriptionRequest

	err := ctx.GinContext.BindJSON(&request)
	rlog.Debug(ctx.GinContext.GetRawData())

	if err != nil {
		return errs.NewClientError(err.Error())
	}

	var subscriber data.Subscriber

	subscriber.ClassYear = request.ClassYear
	subscriber.Email = request.EmailAddress
	subscriber.ProgramName = request.ProgramName
	subscriber.FirstName = request.FirstName
	subscriber.LastName = request.LastName

	// create new subscription
	if err := ctx.Db.Create(subscriber).Error; err != nil {
		return errs.NewClientError("Unable to create new subscription")
	}

	ctx.Result = SubscriptionResponse{"Ok"}
	return nil
}
