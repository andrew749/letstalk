package converters

import (
	"encoding/json"
	"letstalk/server/core/api"
	"letstalk/server/data"

	raven "github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
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
	dataMap := make(map[string]interface{})

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

func NotificationsDataToApi(dataNotifs []data.Notification) ([]api.Notification, error) {
	apiNotifs := make([]api.Notification, 0, len(dataNotifs))
	for _, dataNotif := range dataNotifs {
		apiNotif, err := NotificationDataToApi(dataNotif)
		if err != nil {
			rlog.Errorf("Unable to deserialize notification: %+v", err)
			raven.CaptureError(err, nil)
			continue
		}
		apiNotifs = append(apiNotifs, *apiNotif)
	}

	return apiNotifs, nil
}
