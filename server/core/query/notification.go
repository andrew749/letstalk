package query

import (
	"encoding/json"

	"letstalk/server/core/api"
	"letstalk/server/core/errs"
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
	}

	return dataNotif, nil
}

func NotificationDataToApi(notification data.Notification) (*api.Notification, error) {
	dataMap := make(map[string]string)

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
		Message:        notification.Message,
		ThumbnailLink:  notification.ThumbnailLink,
		Data:           dataMap,
	}

	return apiNotif, nil
}

func NotificationsDataToApi(dataNotifs []data.Notification) ([]api.Notification, error) {
	apiNotifs := make([]api.Notification, len(dataNotifs))
	for i, dataNotif := range dataNotifs {
		apiNotif, err := NotificationDataToApi(dataNotif)
		if err != nil {
			return nil, err
		}
		apiNotifs[i] = *apiNotif
	}

	return apiNotifs, nil
}

func GetNewestNotificationsForUser(
	db *gorm.DB,
	userId data.TUserID,
	limit int,
) ([]api.Notification, errs.Error) {
	var dataNotifs []data.Notification

	err := db.Order("id desc").Where(
		&data.Notification{UserId: userId},
	).Limit(limit).Find(&dataNotifs).Error
	if err != nil {
		return nil, errs.NewDbError(err)
	}

	apiNotifs, err := NotificationsDataToApi(dataNotifs)
	if err != nil {
		return nil, errs.NewInternalError(err.Error())
	}

	return apiNotifs, nil
}

func GetNotificationsForUser(
	db *gorm.DB,
	userId data.TUserID,
	past int,
	limit int,
) ([]api.Notification, errs.Error) {
	var dataNotifs []data.Notification

	err := db.Order("id desc").Where(
		"user_id = ? and id < ?",
		userId,
		past,
	).Limit(limit).Find(&dataNotifs).Error
	if err != nil {
		return nil, errs.NewDbError(err)
	}

	apiNotifs, err := NotificationsDataToApi(dataNotifs)
	if err != nil {
		return nil, errs.NewInternalError(err.Error())
	}

	return apiNotifs, nil
}
