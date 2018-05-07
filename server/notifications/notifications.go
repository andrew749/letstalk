package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	EXPO_HOST = "https://exp.host"
	API_URL   = "/--/api/v2"
	PUSH_API  = "/push/send"
)

type Notification struct {
	To    string `json:"to"`
	Title string `json:"title"`
	Body  string `json:"body"`

	// extra stuff
	Data interface{} `json:"data,omitempty"`

	// default to play, nothing to play no sound
	Sound *string `json:"sound,omitempty"`

	// how long to keep message for redelivery
	TTL *int `json:"ttl,omitempty"`

	// unix timestamp for when message should go away
	Expiration *int `json:"expiration,omitempty"`

	// default, normal or high
	Priority *string `json:"priority,omitempty"`

	// unread notification count
	Badge *int `json:"badge,omitempty"`
}

func SendNotification(notification Notification) error {
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

func CreateAndSendNotification(deviceToken string, message string, title string) error {
	data := map[string]string{"message": message, "title": title}
	notification := Notification{
		To:    deviceToken,
		Title: title,
		Body:  message,
		Data:  data,
	}

	return SendNotification(notification)
}
