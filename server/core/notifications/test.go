package notifications

import (
	"letstalk/server/core/test"
	"letstalk/server/data"
	"letstalk/server/queue/queues/notification_queue"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

var mockQueue chan notification_queue.NotificationQueueData

func TestSerializeDeserializeNotification(t *testing.T) {
	thisTest := test.Test{
		Test: func(db *gorm.DB) {
			thumbnail := "thumbnail"
			var err error
			data, err := CreateNotification(
				db,
				1,
				data.NOTIF_TYPE_ADHOC,
				"title",
				"message",
				&thumbnail,
				time.Now(),
				map[string]interface{}{"test": "test"},
				"",
				nil,
			)
			assert.NoError(t, err)
			mockQueue <- notification_queue.DataNotificationModelToQueueModel(*data)
			rcv := <-mockQueue
			rcvData, err := notification_queue.QueueModelToDataNotificationModel(db, rcv)
			assert.NoError(t, err)
			assert.Equal(t, rcvData.Message, data.Message)
		},
		TestName: "Test Serializing a data notification and then deserializing it",
	}
	test.RunTestWithDb(thisTest)
}
