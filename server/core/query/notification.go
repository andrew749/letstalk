package query

import (
	"encoding/json"

	"letstalk/server/core/api"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

func notificationApiToData(notification api.Notification) (*data.Notification, error) {
	notifData, err := json.Marshal(notification.Data)
	if err != nil {
		return nil, err
	}

	dataNotif := &data.Notification{
		Model: gorm.Model{
			ID: notification.NotificationId,
		},
		UserId: notification.UserId,
		Type:   notification.Type,
		State:  notification.State,
		Data:   notifData,
	}
	if notification.CreatedAt != nil {
		dataNotif.CreatedAt = *notification.CreatedAt
	}

	return dataNotif, nil
}

func notificationDataToApi(notification data.Notification) (*api.Notification, error) {
	dataMap := make(map[string]string)

	err := json.Unmarshal(notification.Data, &dataMap)
	if err != nil {
		return nil, err
	}

	apiNotif := &api.Notification{
		notification.ID,
		notification.UserId,
		notification.Type,
		notification.State,
		dataMap,
		&notification.CreatedAt,
	}

	return apiNotif, nil
}

func notificationsDataToApi(dataNotifs []data.Notification) ([]api.Notification, error) {
	apiNotifs := make([]api.Notification, len(dataNotifs))
	for i, dataNotif := range dataNotifs {
		apiNotif, err := notificationDataToApi(dataNotif)
		if err != nil {
			return nil, err
		}
		apiNotifs[i] = *apiNotif
	}

	return apiNotifs, nil
}

func CreateNotification(
	db *gorm.DB,
	userId data.TUserID,
	tpe data.NotifType,
	dataMap map[string]string,
) (*api.Notification, errs.Error) {
	notifData, err := json.Marshal(dataMap)
	if err != nil {
		return nil, errs.NewInternalError(err.Error())
	}

	dataNotif := &data.Notification{
		UserId: userId,
		Type:   tpe,
		State:  data.NOTIF_STATE_UNREAD,
		Data:   notifData,
	}

	if err := db.Create(dataNotif).Error; err != nil {
		return nil, errs.NewDbError(err)
	}

	// NOTE: This should not error since we just marshalled the data, so didn't add any logic for
	// deleting corrupt notifications.
	apiNotif, err := notificationDataToApi(*dataNotif)
	if err != nil {
		return nil, errs.NewInternalError(err.Error())
	}

	return apiNotif, nil
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

	apiNotifs, err := notificationsDataToApi(dataNotifs)
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

	apiNotifs, err := notificationsDataToApi(dataNotifs)
	if err != nil {
		return nil, errs.NewInternalError(err.Error())
	}

	return apiNotifs, nil
}

func UpdateNotificationState(
	db *gorm.DB,
	userId data.TUserID,
	notificationIds []uint,
	state data.NotifState,
) errs.Error {
	err := db.Model(&data.Notification{}).Where("id in (?) and user_id = ?",
		notificationIds,
		userId,
	).Updates(&data.Notification{State: state}).Error
	if err != nil {
		return errs.NewDbError(err)
	}
	return nil
}