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

		// send each notification individually
		// TODO: migrate to one project to fix this bug so we can send batches
		for _, sendNotification := range *sendNotifications {
			rlog.Debugf("Sending notification to %s", sendNotification.To)

			// this is what we need to do in the future
			// res, err := notification_api.SendNotifications(*sendNotifications)

			// ensure idempotent behaviour for each device
			// if there is already a sent notification for this, don't send one
			exists, err := data.ExistsPendingNotification(db, notification.ID, sendNotification.To)
			if err != nil {
				rlog.Error(err)
				continue
			}

			if exists == true {
				rlog.Debug("Notification already sent")
				continue
			}

			// instead send a single notification
			res, err := notification_api.SendNotifications([]notification_api.ExpoNotification{sendNotification})

			// if there was an error sending this notification (not 200 status code)
			if err != nil {
				rlog.Error(err)
				raven.CaptureError(err, nil)
				return err
			}

			rlog.Debugf("Sent Notifications and got response: %v", res)

			// for the single notification response, add the receipt
			// none of these should fail since we have a successful response from expo
			for i, response := range res.Data {
				rlog.Debugf("Processing response: %+v", response)
				temp, err := data.CreateNewPendingNotification(db, notification.ID, (*sendNotifications)[i].To)
				if err != nil {
					raven.CaptureError(err, nil)
					continue
				}

				// was this message errored
				if response.Status == notification_api.ERROR_STATUS {
					failureType := notification_api.ExpoNotificationFailureType(response.Details.Error)
					if err := temp.MarkNotificationError(
						db,
						response.Message,
						response.Details,
						&failureType,
					); err != nil {
						rlog.Error(err)
						raven.CaptureError(err, nil)
						continue
					}
				}

				if err = temp.MarkNotificationSent(db, response.Id); err != nil {
					rlog.Error(err)
					raven.CaptureError(err, nil)
					continue
				}

				rlog.Debug("Done processing response")
			}

		}
	}
	// will never reach here
	return nil
}
