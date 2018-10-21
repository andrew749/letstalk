package api

import (
	"time"

	"letstalk/server/data"
)

type Notification struct {
	NotificationId uint                   `json:"notificationId"`
	UserId         data.TUserID           `json:"userId"`
	Type           data.NotifType         `json:"type"`
	State          data.NotifState        `json:"state"`
	Title          string                 `json:"title"`
	Message        string                 `json:"message"`
	Timestamp      time.Time              `json:"timestamp"`
	ThumbnailLink  *string                `json:"thumbnail"`
	Data           map[string]interface{} `json:"data"`
	Link           *string                `json:"link"`
	RunId          *string                `json:"runId"`
}

type UpdateNotificationStateRequest struct {
	NotificationIds []uint          `json:"notificationIds"`
	State           data.NotifState `json:"state"`
}

type SendAdhocNotificationRequest struct {
	Recipient      int     `json:"recipient" binding:"required"`
	Message        string  `json:"message" binding:"required"`
	Title          string  `json:"title" binding:"required"`
	Thumbnail      *string `json:"thumbnail"`
	TemplatePath   string  `json:"templatePath" binding:"required"`
	TemplateParams string  `json:"templateParams"`
	Link           string  `json:"link"`
}
