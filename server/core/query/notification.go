package query

import (
	"database/sql"
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
	// HACK: for some reason data is getting corrupted when being passed on the stack
	rows, err := db.DB().Query("select * from notifications where user_id=? order by id desc limit ?", userId, limit)
	if err != nil {
		return nil, errs.NewDbError(err)
	}
	defer rows.Close()

	apiNotifs, err := getApiNotificationFromRows(rows)

	if err != nil {
		return nil, errs.NewInternalError(err.Error())
	}

	return apiNotifs, nil
}

func getApiNotificationFromRows(rows *sql.Rows) ([]api.Notification, error) {
	apiNotifs := make([]api.Notification, 0)

	for rows.Next() {
		var notification data.Notification
		err := rows.Scan(
			&notification.ID,
			&notification.CreatedAt,
			&notification.UpdatedAt,
			&notification.DeletedAt,
			&notification.UserId,
			&notification.Type,
			&notification.Timestamp,
			&notification.State,
			&notification.Title,
			&notification.Message,
			&notification.ThumbnailLink,
			&notification.Data,
			&notification.Link,
			&notification.RunId,
		)

		converted, err := converters.NotificationDataToApi(notification)
		if err != nil {
			return nil, errs.NewInternalError(err.Error())
		}

		apiNotifs = append(apiNotifs, *converted)
	}

	return apiNotifs, nil
}

func GetNotificationsForUser(
	db *gorm.DB,
	userId data.TUserID,
	past int,
	limit int,
) ([]api.Notification, errs.Error) {
	// HACK: for some reason data is getting corrupted when being passed on the stack
	rows, err := db.DB().Query("select * from notifications where user_id = ? and id < ? order by id desc limit ?", userId, past, limit)
	if err != nil {
		return nil, errs.NewDbError(err)
	}
	defer rows.Close()

	apiNotifs, err := getApiNotificationFromRows(rows)

	if err != nil {
		return nil, errs.NewInternalError(err.Error())
	}

	return apiNotifs, nil
}
