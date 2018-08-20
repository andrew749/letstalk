package controller

import (
	"fmt"
	"strconv"

	"letstalk/server/aws_utils"
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	notification_helper "letstalk/server/core/notifications"
	"letstalk/server/core/query"
	"letstalk/server/data"
	"letstalk/server/jobs"
	"letstalk/server/notifications"

	"github.com/jinzhu/gorm"
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

	// TODO(acod): remove hardcoded
	var notificationToken = &data.NotificationToken{
		Token:   request.Token,
		Service: "expo", // hardcoded for now
	}

	err = c.WithinTx(func(db *gorm.DB) error {
		tx := db.Begin()
		// add the token to the
		tx.Create(&notificationToken)
		tx.Model(&data.Session{}).
			Where("session_id = ?", c.SessionData.SessionId).
			Update("notification_token", request.Token)
		if tx.Error != nil {
			return tx.Error
		}

		return nil
	})

	if err != nil {
		return errs.NewRequestError(err.Error())
	}
	c.Result = "Ok"

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

func GetNotifications(c *ctx.Context) errs.Error {
	db := c.Db
	userId := c.SessionData.UserId
	q := c.GinContext.Request.URL.Query()
	var (
		limitStrs []string
		pastStrs  []string
		ok        bool
		err       errs.Error
		apiNotifs []api.Notification
	)

	if limitStrs, ok = q["limit"]; !ok || len(limitStrs) == 0 {
		return errs.NewRequestError("Must provide query param `limit`")
	}

	limit, convErr := strconv.Atoi(limitStrs[0])
	if convErr != nil {
		return errs.NewRequestError(convErr.Error())
	}

	if pastStrs, ok = q["past"]; ok && len(pastStrs) > 0 {
		past, convErr := strconv.Atoi(pastStrs[0])
		if convErr != nil {
			return errs.NewRequestError(convErr.Error())
		}
		apiNotifs, err = query.GetNotificationsForUser(db, userId, past, limit)
		if err != nil {
			return err
		}
	} else {
		apiNotifs, err = query.GetNewestNotificationsForUser(db, userId, limit)
		if err != nil {
			return err
		}
	}

	// dataMap := make(map[string]string)
	// dataMap["credentialName"] = "Software Engineer at Quora"
	// dataMap["userName"] = "Wojtek Swiderski"
	// dataMap["side"] = "ASKER"
	//
	// // TODO: Remove
	// _, err = notification_helper.CreateNotification(db, userId, data.NOTIF_TYPE_NEW_CREDENTIAL_MATCH, "New match", nil, time.Now(), dataMap)
	// if err != nil {
	// 	return err
	// }

	c.Result = apiNotifs
	return nil
}

func UpdateNotificationState(c *ctx.Context) errs.Error {
	var req api.UpdateNotificationStateRequest

	if err := c.GinContext.BindJSON(&req); err != nil {
		return errs.NewRequestError(err.Error())
	}

	if err := notification_helper.UpdateNotificationState(
		c.Db,
		c.SessionData.UserId,
		req.NotificationIds,
		req.State,
	); err != nil {
		return err
	}

	return nil
}
