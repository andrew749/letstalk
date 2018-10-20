package notifications

import (
	"encoding/json"
	"errors"
	"letstalk/server/core/api"
	"letstalk/server/core/ctx"
	"letstalk/server/core/errs"
	"letstalk/server/data"
	"net/http"
	"strconv"

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

	ctx.GinContext.HTML(
		http.StatusOK,
		notificationPage.TemplateLink,
		map[string]interface{}{"WebpageTitle": d["title"], "Data": d},
	)

	return nil
}

func createTestUser() data.User {
	birthdate := "Oct 7 1996"
	return data.User{
		UserId:    1,
		FirstName: "Andrew",
		LastName:  "Codispoti",
		Email:     "thegripper@gmail.com",
		Secret:    "somesecret",
		Gender:    data.GENDER_MALE,
		Birthdate: &birthdate,
	}
}

func EchoNotificationPage(ctx *ctx.Context) errs.Error {
	var notificationReq api.NotificationEchoRequest
	if err := ctx.GinContext.BindJSON(&notificationReq); err != nil {
		return errs.NewRequestError("Unable to bind json %s", err)
	}

	notificationReq.Data["user"] = createTestUser()
	ctx.GinContext.HTML(
		http.StatusOK,
		notificationReq.TemplateLink,
		map[string]interface{}{"Data": notificationReq.Data},
	)
	return nil
}
