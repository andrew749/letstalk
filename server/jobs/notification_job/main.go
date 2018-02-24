package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"letstalk/server/core/notifications"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

const (
	EXPO_HOST = "https://exp.host"
	API_URL   = "/--/api/v2"
	PUSH_API  = "/push/send"
)

func CreateAndSendNotification(deviceToken string, message string, title string) error {
	notification := notifications.Notification{
		To:    fmt.Sprintf("ExponentPushToken[%s]", deviceToken),
		Title: title,
		Body:  message,
	}

	return SendNotification(notification)
}

func SendNotification(notification notifications.Notification) error {

	marshalledNotification, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s%s", EXPO_HOST, API_URL, PUSH_API),
		bytes.NewBuffer(marshalledNotification),
	)
	client := &http.Client{}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept-encoding", "gzip, deflate")
	log.Print("Sending Notification to Expo")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	log.Print(bodyString)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("Successfully sent notification to client: %s", notification.To)
	return nil
}

func HandleRequest(ctx context.Context, notification notifications.Notification) error {
	return SendNotification(notification)
}

func main() {
	lambda.Start(HandleRequest)
}
