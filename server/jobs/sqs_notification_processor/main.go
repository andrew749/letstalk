package main

import (
	"context"
	"encoding/json"
	"letstalk/server/core/notifications"
	"letstalk/server/data"
	notification_api "letstalk/server/notifications"
	"letstalk/server/utility"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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
			return err
		}
		db = conn
	}
	// TODO: handle error
	rlog.Printf("Received message %#v\n", sqsEvent)
	var notification data.Notification
	records := sqsEvent.Records
	// Only handles one record since each sqs message only contains at most notification.
	for _, record := range records {
		err := json.Unmarshal([]byte(record.Body), &notification)
		if err != nil {
			return err
		}

		sendNotifications, err := notifications.NotificationsFromNotificationDataModel(db, notification)

		if err != nil {
			return err
		}

		for _, not := range *sendNotifications {
			res, err := notification_api.SendNotifications(not)
			if err != nil {
				rlog.Error(err)
				return err
			}
			// update notification state

			// currently only one message being sent at a time
			notification.UpdateReceipt(db, res.Data[0].Id)
		}

		return nil
	}
	// will never reach here
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
