package email_subscription

import (
	"fmt"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"
	"letstalk/server/push"

	"github.com/romana/rlog"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

	var err error

	if err = ctx.GinContext.BindJSON(&request); err != nil {
		return errs.NewClientError(err.Error())
	}

	var subscriber data.Subscriber

	subscriber.ClassYear = request.ClassYear
	subscriber.Email = request.EmailAddress
	subscriber.ProgramName = request.ProgramName
	subscriber.FirstName = request.FirstName
	subscriber.LastName = request.LastName

	// create new subscription
	if err = ctx.Db.Create(subscriber).Error; err != nil {
		return errs.NewClientError("Unable to create new subscription")
	}

	// send verification email
	to := mail.NewEmail(
		fmt.Sprintf("%s %s", subscriber.FirstName, subscriber.LastName),
		subscriber.Email,
	)

	err = push.SendSubscribeEmail(to, subscriber.FirstName)
	if err != nil {
		rlog.Error("Unable to send email to ", subscriber.Email)
	}

	ctx.Result = SubscriptionResponse{"Ok"}
	return nil
}
