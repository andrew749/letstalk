package notifications

import (
	"fmt"
	"letstalk/server/aws_utils"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/jobs"

	"github.com/mijia/modelq/gmq"
	"github.com/romana/rlog"
)

type Notification struct {
	To    string `json:"to"`
	Title string `json:"title"`
	Body  string `json:"body"`

	// extra stuff
	Data *interface{} `json:"data,omitempty"`

	// default to play, nothing to play no sound
	Sound *string `json:"sound,omitempty"`

	// how long to keep message for redelivery
	TTL *int `json:"ttl,omitempty"`

	// unix timestamp for when message should go away
	Expiration *int `json:"expiration,omitempty"`

	// default, normal or high
	Priority *string `json:"priority,omitempty"`

	// unread notification count
	Badge *int `json:"badge,omitempty"`
}

/**
 * Handle new token submission.
 */
func NewNotificationTokenHandler(c *ctx.Context, notificationToken string) errs.Error {

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
