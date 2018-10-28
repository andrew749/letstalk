package controller

import (
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/notifications"
	"letstalk/server/data"

	"github.com/romana/rlog"
)

// NotificationCampaignController Send out a notification campaign
func NotificationCampaignController(ctx *ctx.Context) errs.Error {
	var notificationCampaignReq api.NotificationCampaignSendRequest
	if err := ctx.GinContext.BindJSON(&notificationCampaignReq); err != nil {
		return errs.NewRequestError("Bad request: %+v", err)
	}

	db := ctx.Db

	if exists, err := data.ExistsCampaign(db, notificationCampaignReq.RunId); exists || err != nil {
		if err != nil {
			return errs.NewDbError(err)
		}
		return errs.NewRequestError("Campaign with id %s already exists", notificationCampaignReq.RunId)
	}

	rlog.Infof("Creating notification campaign with RunId %s", notificationCampaignReq.RunId)

	if err := notifications.SendNotificationCampaign(db, notificationCampaignReq); err != nil {
		return errs.NewInternalError(err.Error())
	}
	ctx.Result = "Ok"

	return nil
}
