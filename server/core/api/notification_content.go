package api

import (
	"letstalk/server/data"
)

type NotificationContentPageRequest struct {
	NotificationContentId data.EntID `json:"notificationContentId" binding:"required"`
}

// Dummy data sent to render using server
type NotificationEchoRequest struct {
	TemplateLink string                 `json:"templateLink"`
	Data         map[string]interface{} `json:"data"`
}
