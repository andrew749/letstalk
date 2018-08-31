package query

import (
	"letstalk/server/core/api"
	"letstalk/server/core/converters"
	"letstalk/server/core/errs"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
)

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

	apiNotifs, err := converters.NotificationsDataToApi(dataNotifs)
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

	apiNotifs, err := converters.NotificationsDataToApi(dataNotifs)
	if err != nil {
		return nil, errs.NewInternalError(err.Error())
	}

	return apiNotifs, nil
}
