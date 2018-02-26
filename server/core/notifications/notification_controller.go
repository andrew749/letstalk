package notifications

import (
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/mijia/modelq/gmq"
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
		if _, err = notification_token.Insert(tx); err != nil {
			return err
		}
		// send test notification
		go SendNotification(
			notification_token.Token,
			"Successfully registered for notification",
			"Hive",
		)
		return nil
	})

	if err != nil {
		return errs.NewInternalError("Internal error: %s", err)
	}
	c.Result = "Ok"

	return nil
}
