package main

import (
	"context"
	"letstalk/server/notifications"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, notification notifications.Notification) error {
	return notifications.SendNotification(notification)
}

func main() {
	lambda.Start(HandleRequest)
}
