package controller

import (
	"strconv"

	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/core/query"
	// "letstalk/server/data"
)

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

	// // TODO: Remove
	// _, err = query.CreateNotification(db, userId, data.NOTIF_TYPE_NEW_CREDENTIAL_MATCH, dataMap)
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

	if err := query.UpdateNotificationState(
		c.Db,
		c.SessionData.UserId,
		req.NotificationIds,
		req.State,
	); err != nil {
		return err
	}

	return nil
}
