package controller

import (
	"testing"

	"letstalk/server/core/api"
	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/test_helpers"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestAddExpoDeviceToken(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			userId := data.TUserID(69)
			token := "sdkljahkljfga843hfal843afewahjgla9493hlg854lh3daf4389lrahldf453"
			c := test_helpers.CreateTestContext(t, db, userId)
			req := api.AddExpoDeviceTokenRequest{Token: token}

			err := handleAddExpoDeviceToken(c, req)
			assert.NoError(t, err)

			var devices []data.UserDevice
			dbErr := db.Where(&data.UserDevice{UserId: userId}).Find(&devices).Error
			assert.NoError(t, dbErr)
			assert.Equal(t, 1, len(devices))
			assert.Equal(t, data.UserDevice{
				UserId:                userId,
				NotificationToken:     token,
				NotificationTokenType: data.EXPO_PUSH,
			}, devices[0])

			// Adding again shouldn't add another token (should be idempotent)
			err = handleAddExpoDeviceToken(c, req)
			assert.NoError(t, err)

			dbErr = db.Where(&data.UserDevice{UserId: userId}).Find(&devices).Error
			assert.NoError(t, dbErr)
			assert.Equal(t, 1, len(devices))
			assert.Equal(t, data.UserDevice{
				UserId:                userId,
				NotificationToken:     token,
				NotificationTokenType: data.EXPO_PUSH,
			}, devices[0])
		},
	}
	test.RunTestWithDb(thisTest)
}
