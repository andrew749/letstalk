package notifications

import (
	"fmt"
	"letstalk/server/aws_utils"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"
	"letstalk/server/jobs"

	"github.com/mijia/modelq/gmq"
	"github.com/romana/rlog"
)

type NotificationTokenSubmissionRequest struct {
	Token string `json:"token" binding:"required"`
}

func GetNewNotificationToken(c *ctx.Context) errs.Error {
	var request NotificationTokenSubmissionRequest
	err := c.GinContext.BindJSON(&request)
	if err != nil {
		return errs.NewClientError("Bad Request: %s", err)
	}

	notification_token := data.NotificationTokens{
		UserId: c.SessionData.UserId,
		Token:  request.Token,
	}

	err = gmq.WithinTx(c.Db, func(tx *gmq.Tx) error {
		// check if this token already exists
		if _, err = notification_token.Insert(tx); err != nil {
			return err
		}
		// send test notification
		rlog.Debug("Dispatching notification lambda")
		aws_utils.DispatchLambdaJob(
			jobs.SendNotification,
			Notification{
				To:    fmt.Sprintf("ExponentPushToken[%s]", notification_token.Token),
				Body:  "Subscribed for notifications.",
				Title: "Hive",
			},
		)

		return nil
	})

	if err != nil {
		return errs.NewInternalError("Internal error: %s", err)
	}
	c.Result = "Ok"

	return nil
}
