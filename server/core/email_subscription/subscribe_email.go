package email_subscription

import (
	"fmt"
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"
	"letstalk/server/email"

	"github.com/romana/rlog"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func AddSubscription(ctx *ctx.Context) errs.Error {
	var request api.SubscriptionRequest

	var err error

	if err = ctx.GinContext.BindJSON(&request); err != nil {
		return errs.NewClientError(err.Error())
	}

	var subscribers []data.Subscriber

	// if there is already a subscription
	if err = ctx.Db.Where(
		"email = ?",
		request.EmailAddress,
	).Find(&subscribers).Error; err != nil {
		return errs.NewInternalError(err.Error())
	}

	if len(subscribers) > 0 {
		return errs.NewClientError("Subscription already created")
	}

	var subscriber data.Subscriber
	// create new subscription
	if err = ctx.Db.FirstOrCreate(&subscriber, data.Subscriber{
		ClassYear:   request.ClassYear,
		Email:       request.EmailAddress,
		ProgramName: request.ProgramName,
		FirstName:   request.FirstName,
		LastName:    request.LastName,
	}).Error; err != nil {
		return errs.NewClientError("Unable to create new subscription")
	}

	// send verification email
	to := mail.NewEmail(
		fmt.Sprintf("%s %s", subscriber.FirstName, subscriber.LastName),
		subscriber.Email,
	)

	err = email.SendSubscribeEmail(to, subscriber.FirstName)
	if err != nil {
		rlog.Error("Unable to send email to ", subscriber.Email)
	}

	ctx.Result = api.SubscriptionResponse{
		Status: "Ok",
	}
	return nil
}
