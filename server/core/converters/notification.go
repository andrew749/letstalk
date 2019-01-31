package converters

import (
	"encoding/json"
	"letstalk/server/core/api"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func NotificationApiToData(notification api.Notification) (*data.Notification, error) {
	notifData, err := json.Marshal(notification.Data)
	if err != nil {
		return nil, err
	}

	dataNotif := &data.Notification{
		Model: gorm.Model{
			ID: notification.NotificationId,
		},
		UserId:        notification.UserId,
		Type:          notification.Type,
		State:         notification.State,
		Message:       notification.Message,
		ThumbnailLink: notification.ThumbnailLink,
		Timestamp:     notification.Timestamp,
		Data:          notifData,
		Link:          notification.Link,
		RunId:         notification.RunId,
	}

	return dataNotif, nil
}

func NotificationDataToApi(notification data.Notification) (*api.Notification, error) {
	var dataMap map[string]interface{}

	err := json.Unmarshal(notification.Data, &dataMap)
	if err != nil {
		return nil, err
	}

	apiNotif := &api.Notification{
		NotificationId: notification.ID,
		UserId:         notification.UserId,
		Type:           notification.Type,
		State:          notification.State,
		Timestamp:      notification.Timestamp,
		Title:          notification.Title,
		Message:        notification.Message,
		ThumbnailLink:  notification.ThumbnailLink,
		Data:           dataMap,
		Link:           notification.Link,
	}

	return apiNotif, nil
}
