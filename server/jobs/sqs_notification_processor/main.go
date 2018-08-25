package main

import (
	"context"
	"fmt"

	sqs_notification_processor "letstalk/server/jobs/sqs_notification_processor/src"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// HandleRequest Handle the message data passed to the lambda from sqs
func HandleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	return sqs_notification_processor.SendNotificationLambda(sqsEvent)
}

func testNotification(id uint) error {
	return sqs_notification_processor.SendNotificationLambda(events.SQSEvent{
		Records: []events.SQSMessage{
			events.SQSMessage{
				Body: fmt.Sprintf("{\"id\": %d}", id),
			},
		},
	})
}

func main() {
	lambda.Start(HandleRequest)
	// test the job locally
	// if err := testNotification(49); err != nil {
	// 	panic(err)
	// }
}
