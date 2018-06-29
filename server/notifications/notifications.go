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

type NotifType string

const (
	NOTIF_TYPE_REQUEST_TO_MATCH NotifType = "REQUEST_TO_MATCH"
	NOTIF_TYPE_NEW_MATCH        NotifType = "NEW_MATCH"
	NOTIF_TYPE_MATCH_VERIFIED   NotifType = "MATCH_VERIFIED"
)

func CreateAndSendNotificationWithData(
	deviceToken string,
	message string,
	title string,
	tpe NotifType,
	extraData map[string]interface{},
) error {
	data := map[string]interface{}{"message": message, "title": title, "type": string(tpe)}
	for key, value := range extraData {
		data[key] = value
	}

	notification := Notification{
		To:    deviceToken,
		Title: title,
		Body:  message,
		Data:  data,
	}

	return SendNotification(notification)
}

func CreateAndSendNotification(
	deviceToken string,
	message string,
	title string,
	tpe NotifType,
) error {
	return CreateAndSendNotificationWithData(
		deviceToken,
		message,
		title,
		tpe,
		make(map[string]interface{}),
	)
}

// SPECIFIC NOTIFICATION MESSAGES

type RequestToMatchSide string

const (
	REQUEST_TO_MATCH_SIDE_ASKER    RequestToMatchSide = "ASKER"
	REQUEST_TO_MATCH_SIDE_ANSWERER RequestToMatchSide = "ANSWERER"
)

func RequestToMatchNotification(
	deviceToken string,
	side RequestToMatchSide,
	requestId uint,
	name string,
) error {
	var (
		extraData map[string]interface{} = map[string]interface{}{"side": side, "requestId": requestId}
		title     string                 = "You got a match!"
		message   string                 = fmt.Sprintf("You got matched for \"%s\"", name)
	)
	return CreateAndSendNotificationWithData(
		deviceToken,
		message,
		title,
		NOTIF_TYPE_REQUEST_TO_MATCH,
		extraData,
	)
}

func NewMatchNotification(deviceToken string, message string) error {
	title := "You got a match!"
	return CreateAndSendNotificationWithData(
		deviceToken,
		message,
		title,
		NOTIF_TYPE_NEW_MATCH,
		nil,
	)
}

func NewMentorNotification(deviceToken string) error {
	return NewMatchNotification(deviceToken, "You were matched with a new mentor.")
}

func NewMenteeNotification(deviceToken string) error {
	return NewMatchNotification(deviceToken, "You were matched with a new mentee.")
}

func MatchVerifiedNotification(deviceToken string, userName string) error {
	title := "You verified a match!"
	message := fmt.Sprintf("Your match with %s is now verified.", userName)
	return CreateAndSendNotificationWithData(
		deviceToken,
		message,
		title,
		NOTIF_TYPE_MATCH_VERIFIED,
		nil,
	)
}
