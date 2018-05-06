package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"letstalk/server/notifications"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, notification notifications.Notification) error {
	return notificationsSendNotification(notification)
}

func main() {
	lambda.Start(HandleRequest)
}
