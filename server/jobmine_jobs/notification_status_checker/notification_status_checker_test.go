package notification_status_checker

import (
	"letstalk/server/core/test"
	"letstalk/server/data"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func cleanPendingNotificationTable(db *gorm.DB) {
	err := db.Delete(&data.ExpoPendingNotification{}).Error
	if err != nil {
		panic(err)
	}
}

func TestNotificationStatusChecker(t *testing.T) {
	tests := []test.Test{
		test.Test{
			Test: func(db *gorm.DB) {
				cleanPendingNotificationTable(db)
				pendingNotification, err := data.CreateNewPendingNotification(db, 2, "testDevice")
				assert.NoError(t, err)
				pendingNotification.Checked = true
				assert.NoError(t, db.Save(pendingNotification).Error)
				err = pendingNotification.MarkNotificationChecked(db)
				assert.NoError(t, err)
				notifications, err := PendingNotificationsToCheck(db, nil)
				assert.NoError(t, err)
				assert.Len(t, notifications, 0)
			},
			TestName: "Test don't find notifications that are checked",
		},
		test.Test{
			Test: func(db *gorm.DB) {
				cleanPendingNotificationTable(db)
				pendingNotification, err := data.CreateNewPendingNotification(db, 3, "testDevice")
				assert.NoError(t, err)
				currentTime := time.Now()
				hourLater := currentTime.Add(time.Hour * time.Duration(1))
				pendingNotification.CreatedAt = hourLater
				assert.NoError(t, db.Save(pendingNotification).Error)
				notifications, err := PendingNotificationsToCheck(db, &currentTime)
				assert.NoError(t, err)
				assert.Len(t, notifications, 0)
			},
			TestName: "Test don't find notifications that newer than threshold",
		},
		test.Test{
			Test: func(db *gorm.DB) {
				cleanPendingNotificationTable(db)
				_, err := data.CreateNewPendingNotification(db, 1, "testDevice")
				assert.NoError(t, err)
				notifications, err := PendingNotificationsToCheck(db, nil)
				assert.NoError(t, err)
				assert.Len(t, notifications, 1)
			},
			TestName: "Test finding notifications that aren't checked",
		},

		test.Test{
			Test: func(db *gorm.DB) {
				cleanPendingNotificationTable(db)
				pendingNotification, err := data.CreateNewPendingNotification(db, 4, "testDevice")
				assert.NoError(t, err)
				currentTime := time.Now()
				hourBefore := currentTime.Add(time.Hour * time.Duration(-1))
				pendingNotification.CreatedAt = hourBefore
				assert.NoError(t, db.Save(pendingNotification).Error)
				notifications, err := PendingNotificationsToCheck(db, &currentTime)
				assert.NoError(t, err)
				assert.Len(t, notifications, 1)
			},
			TestName: "Test find notifications that are older than threshold",
		},
	}

	test.RunTestsWithDb(tests)
}
