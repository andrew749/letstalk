package notifications

import (
	"fmt"
	"letstalk/server/aws_utils"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"
	"letstalk/server/jobs"
	"letstalk/server/notifications"

	"github.com/romana/rlog"
)

type NotificationTokenSubmissionRequest struct {
	Token string `json:"token" binding:"required"`
}

func GetNewNotificationToken(c *ctx.Context) errs.Error {
	var request NotificationTokenSubmissionRequest
	err := c.GinContext.BindJSON(&request)
	if err != nil {
		return errs.NewRequestError("Bad Request: %s", err)
	}

	db := c.Db

	tx := db.Begin()
	// TODO(acod): remove hardcoded
	var notificationToken = &data.NotificationToken{
		Token:   request.Token,
		Service: "expo", // hardcoded for now
	}
	// add the token to the
	tx.Create(&notificationToken)
	tx.Model(&data.Session{}).
		Where("session_id = ?", c.SessionData.SessionId).
		Update("notification_token", request.Token)
	if tx.Error != nil {
		tx.Rollback()
		return errs.NewRequestError(tx.Error.Error())
	}

	c.Result = "Ok"
	tx.Commit()

	rlog.Debug("Dispatching notification lambda")
	if err := aws_utils.DispatchLambdaJob(
		jobs.SendNotification,
		notifications.Notification{
			To:    fmt.Sprintf("ExponentPushToken[%s]", notificationToken.Token),
			Body:  "Subscribed for notifications.",
			Title: "Hive",
		},
	); err != nil {
		rlog.Error(err)
	}

	return nil
}
