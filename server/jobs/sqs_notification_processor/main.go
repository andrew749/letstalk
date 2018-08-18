package main

import (
	"context"
	"encoding/json"
	"letstalk/server/data"
	"letstalk/server/notifications"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jinzhu/gorm"
	"github.com/romana/rlog"
)

var db *gorm.DB

// HandleRequest Handle the message data passed to the lambda from sqs
func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	if db == nil {
		// TODO: establish connection
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

		var sendNotification = &notifications.Notification{}
		sendNotification = sendNotification.FromNotificationDataModel(notification)

		_, err = notifications.SendNotification(*sendNotification)
		if err != nil {
			rlog.Error(err)
			return err
		}
		return nil
	}
	// will never reach here
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
