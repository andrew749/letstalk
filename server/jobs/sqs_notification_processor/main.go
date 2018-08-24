package main

import (
	"context"
	"encoding/json"
	"letstalk/server/core/notifications"
	"letstalk/server/data"
	notification_api "letstalk/server/notifications"
	"letstalk/server/queue/queues/notification_queue"
	"letstalk/server/utility"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	raven "github.com/getsentry/raven-go"
	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

var db *gorm.DB

// HandleRequest Handle the message data passed to the lambda from sqs
func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
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
		err := json.Unmarshal([]byte(record.Body), &queueMessage)

		if err != nil {
			rlog.Error(err)
			raven.CaptureError(err, nil)
			return err
		}

		// get the notification model from db given the queue model
		notification, err := notification_queue.QueueModelToDataNotificationModel(db, queueMessage)
		if err != nil {
			rlog.Error(err)
			raven.CaptureError(err, nil)
			return err
		}

		// create a set of notifications to send off to expo
		sendNotifications, err := notifications.NotificationsFromNotificationDataModel(db, notification)
		if err != nil {
			rlog.Error(err)
			raven.CaptureError(err, nil)
			return err
		}

		// send the batch
		res, err := notification_api.SendNotifications(*sendNotifications)
		if err != nil {
			rlog.Error(err)
			raven.CaptureError(err, nil)
			return err
		}

		// create pending notification for each message we tried to send
		tx := db.Begin()

		// for each notification response, add the receipt
		// none of these should fail since we have a successful response from expo
		for i, response := range res.Data {
			temp, err := data.CreateNewPendingNotification(tx, notification.ID, (*sendNotifications)[i].To)

			// was this message errored
			if response.Status == notification_api.ERROR_STATUS {
				err := temp.MarkNotificationError(tx, response.Message, response.Details)
				rlog.Error(err)
				raven.CaptureError(err, nil)
				continue
			}

			err = temp.MarkNotificationSent(db, response.Id)
			if err != nil {
				rlog.Error(err)
				raven.CaptureError(err, nil)
				continue
			}
		}

		return tx.Commit().Error
	}
	// will never reach here
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
