package main

import (
	"context"
	"encoding/json"
	"letstalk/server/notifications"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/romana/rlog"
)

// HandleRequest Handle the message data passed to the lambda from sqs
func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	// TODO: handle error
	rlog.Printf("Received message %#v\n", sqsEvent)
	var notification notifications.Notification
	records := sqsEvent.Records
	// Only handles one record since each sqs message only contains at most notification.
	for _, record := range records {
		err := json.Unmarshal([]byte(record.Body), &notification)
		if err != nil {
			return err
		}
		err = notifications.SendNotification(notification)
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
