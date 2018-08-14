package api

import (
	"time"

	"letstalk/server/data"
)

type Notification struct {
	NotificationId uint              `json:"notificationId"`
	UserId         int               `json:"userId"`
	Type           data.NotifType    `json:"type"`
	State          data.NotifState   `json:"state"`
	Data           map[string]string `json:"data"`
	CreatedAt      *time.Time        `json:"createdAt"`
}

type UpdateNotificationStateRequest struct {
	NotificationIds []uint          `json:"notificationIds"`
	State           data.NotifState `json:"state"`
}
