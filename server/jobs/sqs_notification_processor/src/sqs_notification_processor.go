package sqs_notification_processor

import (
	"encoding/json"
	"letstalk/server/core/notifications"
	"letstalk/server/data"
	"letstalk/server/queue/queues/notification_queue"
	"letstalk/server/utility"

	notification_api "letstalk/server/notifications"

	"github.com/aws/aws-lambda-go/events"
	raven "github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

var db *gorm.DB

func SendNotificationLambda(sqsEvent events.SQSEvent) error {
	// if this is a lambda job in a new execution environment then create a new
	// db connection
	// See https://docs.aws.amazon.com/lambda/latest/dg/running-lambda-code.html
	// which explains execution environments and variable reuse between lambda methods
	if db == nil {
		conn, err := utility.GetDB()
		if err != nil {
			raven.CaptureError(err, nil)
			return err
		}
		db = conn
	}

	rlog.Printf("Received message %#v\n", sqsEvent)
	var queueMessage notification_queue.NotificationQueueData
	records := sqsEvent.Records

	// Only handles one notification record since each sqs message only contains
	// at most one notification.
	for _, record := range records {
		// get the serialized data in sqs
		rlog.Debug("Unmarshalling message")
		err := json.Unmarshal([]byte(record.Body), &queueMessage)

		if err != nil {
			rlog.Error(err)
			raven.CaptureError(err, nil)
			return err
		}
		rlog.Debugf("Unmarshalled message into %v", queueMessage)

		// get the notification model from db given the queue model
		notification, err := notification_queue.QueueModelToDataNotificationModel(db, queueMessage)
		if err != nil {
			rlog.Error(err)
			raven.CaptureError(err, nil)
			return err
		}
		rlog.Debugf("Retrieved Data model from db: %v", notification)

		// create a set of notifications to send off to expo
		sendNotifications, err := notifications.NotificationsFromNotificationDataModel(db, notification)
		if err != nil {
			rlog.Error(err)
			raven.CaptureError(err, nil)
			return err
		}
		rlog.Debugf("Generated Notifications: %v", sendNotifications)

		// send the batch
		res, err := notification_api.SendNotifications(*sendNotifications)
		if err != nil {
			rlog.Error(err)
			raven.CaptureError(err, nil)
			return err
		}
		rlog.Debugf("Sent Notifications got response: %v", res)

		// create pending notification for each message we tried to send
		tx := db.Begin()

		// for each notification response, add the receipt
		// none of these should fail since we have a successful response from expo
		for i, response := range res.Data {
			rlog.Debug("Processing response: %v", response)
			temp, err := data.CreateNewPendingNotification(tx, notification.ID, (*sendNotifications)[i].To)

			// was this message errored
			if response.Status == notification_api.ERROR_STATUS {
				err := temp.MarkNotificationError(tx, response.Message, response.Details)
				rlog.Error(err)
				raven.CaptureError(err, nil)
				continue
			}

			err = temp.MarkNotificationSent(tx, response.Id)
			if err != nil {
				rlog.Error(err)
				raven.CaptureError(err, nil)
				continue
			}
			rlog.Debug("Done processing response")
		}

		return tx.Commit().Error
	}
	// will never reach here
	return nil
}
