package api

import (
	"time"

	"letstalk/server/data"
)

type Notification struct {
	NotificationId uint              `json:"notificationId"`
	UserId         data.TUserID      `json:"userId"`
	Type           data.NotifType    `json:"type"`
	State          data.NotifState   `json:"state"`
	Message        string            `json:"message"`
	Timestamp      time.Time         `json:"timestamp"`
	ThumbnailLink  *string           `json:"thumbnail"`
	Data           map[string]string `json:"data"`
}

type UpdateNotificationStateRequest struct {
	NotificationIds []uint          `json:"notificationIds"`
	State           data.NotifState `json:"state"`
}
