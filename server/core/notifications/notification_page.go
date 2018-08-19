package notifications

import (
	"encoding/json"
	"errors"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"
	"net/http"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

func getNotificationTemplateForId(db *gorm.DB, id uint) (*data.NotificationPage, error) {
	var content data.NotificationPage

	// get the notification content from db
	if err := db.Where(&data.NotificationPage{NotificationId: id}).First(&content).Error; err != nil {
		return nil, errs.NewRequestError(err.Error())
	}

	return &content, nil
}

// gets the query param for notificationId
func getNotificationIdFromContext(ctx *ctx.Context) (*uint, error) {
	c := ctx.GinContext
	rawNotificationID := c.DefaultQuery("notificationId", "")
	if rawNotificationID == "" {
		return nil, errors.New("Missing notificationId query param")
	}

	notificationID, err := strconv.Atoi(rawNotificationID)
	if err != nil {
		return nil, err
	}
	uNotificationID := uint(notificationID)

	return &uNotificationID, nil
}

func createTestNotificationPage(db *gorm.DB, userId uint) error {
	d := []byte("{\"title\":\"Title\", \"body\":\"This is a body\"}")
	n, _ := CreateNotification(db, data.TUserID(userId), data.NOTIF_TYPE_ADHOC, "Test Message", nil, time.Now(), nil)
	req := data.NotificationPage{
		NotificationId: n.ID,
		UserId:         data.TUserID(userId),
		TemplateLink:   "sample_template.html",
		Attributes:     data.JSONBlob(d),
	}
	return db.Save(&req).Error
}

// GetNotificationContentPage Render page with specific attributes for a user
func GetNotificationContentPage(ctx *ctx.Context) errs.Error {
	notificationId, err := getNotificationIdFromContext(ctx)
	if err != nil {
		return errs.NewRequestError(err.Error())
	}

	db := ctx.Db
	// createTestNotificationPage(db, 1)
	notificationPage, err := getNotificationTemplateForId(db, *notificationId)
	if err != nil {
		return errs.NewRequestError(err.Error())
	}

	// access control check
	if notificationPage.UserId != ctx.SessionData.UserId {
		return errs.NewForbiddenError("Not allowed to access")
	}

	var d map[string]interface{}
	err = json.Unmarshal(notificationPage.Attributes, &d)
	if err != nil {
		return errs.NewInternalError(err.Error())
	}
	ctx.GinContext.HTML(http.StatusOK, notificationPage.TemplateLink, &d)

	return nil
}