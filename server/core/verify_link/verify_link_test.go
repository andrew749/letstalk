package verify_link

import (
	"testing"
	"time"

	"letstalk/server/core/errs"
	"letstalk/server/core/test"
	"letstalk/server/data"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestClickLinkOk(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error
			userId := data.TUserID(1)
			linkId, err := CreateLink(db, userId, LINK_TYPE_WHITELIST_WINTER_2019, nil)
			assert.NoError(t, err)

			var link data.UserVerifyLink
			err = db.Where(&data.UserVerifyLink{Id: *linkId}).Find(&link).Error
			assert.NoError(t, err)

			assert.Equal(t, *linkId, link.Id)
			assert.Equal(t, userId, link.UserId)
			assert.False(t, link.Clicked)
			assert.Equal(t, LINK_TYPE_WHITELIST_WINTER_2019, link.Type)
			assert.Nil(t, link.ExpiresAt)

			err = ClickLink(db, *linkId)
			assert.NoError(t, err)

			err = db.Where(&data.UserVerifyLink{Id: *linkId}).Find(&link).Error
			assert.NoError(t, err)

			assert.True(t, link.Clicked)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestClickLinkMissing(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			err := ClickLink(db, data.TVerifyLinkID("not a link"))
			assert.Error(t, err)

			_, isNotFoundError := err.(*errs.NotFoundError)
			assert.True(t, isNotFoundError)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestCreateLinkDistinctLinks(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error
			userId := data.TUserID(1)
			linkId1, err := CreateLink(db, userId, LINK_TYPE_WHITELIST_WINTER_2019, nil)
			assert.NoError(t, err)
			linkId2, err := CreateLink(db, userId, LINK_TYPE_WHITELIST_WINTER_2019, nil)
			assert.NoError(t, err)

			assert.NotEqual(t, *linkId1, *linkId2)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestClickLinkExpiryOk(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error
			userId := data.TUserID(1)
			expiresAt := time.Now().AddDate(0, 0, 1)
			linkId, err := CreateLink(db, userId, LINK_TYPE_WHITELIST_WINTER_2019, &expiresAt)
			assert.NoError(t, err)

			err = ClickLink(db, *linkId)
			assert.NoError(t, err)

			var link data.UserVerifyLink
			err = db.Where(&data.UserVerifyLink{Id: *linkId}).Find(&link).Error
			assert.NoError(t, err)
			assert.True(t, link.Clicked)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestClickLinkExpiryFail(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error
			userId := data.TUserID(1)
			expiresAt := time.Now().AddDate(0, 0, -1)
			linkId, err := CreateLink(db, userId, LINK_TYPE_WHITELIST_WINTER_2019, &expiresAt)
			assert.NoError(t, err)

			err = ClickLink(db, *linkId)
			assert.Error(t, err)

			_, isRequestError := err.(*errs.BadRequest)
			assert.True(t, isRequestError)

			var link data.UserVerifyLink
			err = db.Where(&data.UserVerifyLink{Id: *linkId}).Find(&link).Error
			assert.NoError(t, err)
			assert.False(t, link.Clicked)
		},
	}
	test.RunTestWithDb(thisTest)
}

func TestGetVerifiedUserIds(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			var err error
			user1Id := data.TUserID(1)
			link1Id, err := CreateLink(db, user1Id, LINK_TYPE_WHITELIST_WINTER_2019, nil)
			assert.NoError(t, err)

			user2Id := data.TUserID(2)
			link2Id, err := CreateLink(db, user2Id, LINK_TYPE_WHITELIST_WINTER_2019, nil)
			assert.NoError(t, err)

			user3Id := data.TUserID(3)
			_, err = CreateLink(db, user3Id, LINK_TYPE_WHITELIST_WINTER_2019, nil)
			assert.NoError(t, err)

			err = ClickLink(db, *link1Id)
			assert.NoError(t, err)
			err = ClickLink(db, *link2Id)
			assert.NoError(t, err)

			userIds, err := GetVerifiedUserIds(db, LINK_TYPE_WHITELIST_WINTER_2019)
			assert.NoError(t, err)
			assert.ElementsMatch(t, []data.TUserID{user1Id, user2Id}, userIds)
		},
	}
	test.RunTestWithDb(thisTest)
}
