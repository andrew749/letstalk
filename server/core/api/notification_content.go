package api

import "letstalk/server/data"

type NotificationContentPageRequest struct {
	NotificationContentId data.EntID `json:"notificationContentId" binding:"required"`
}
