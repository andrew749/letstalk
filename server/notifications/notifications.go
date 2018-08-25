package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/romana/rlog"
)

const (
	EXPO_HOST           = "https://exp.host"
	API_URL             = "/--/api/v2"
	PUSH_API            = "/push/send"
	PUSH_RECEIPT_STATUS = "/push/getReceipts"
	OK_STATUS           = "ok"
	ERROR_STATUS        = "error"
)

type ExpoNotification struct {
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

type ExpoNotificationStatusDetails struct {
	Error string `json:"error"`
}

type ExpoNotificationStatusResponse struct {
	Id      string                         `json:"id,omitempty"`
	Status  string                         `json:"status"`
	Message *string                        `json:"message,omitempty"`
	Details *ExpoNotificationStatusDetails `json:"details,omitempty"`
}

type ExpoNotificationStatus struct {
	Data map[string]ExpoNotificationStatusResponse `json:"data"`
}

type ExpoNotificationSendResponse struct {
	Data []ExpoNotificationStatusResponse `json:"data"`
}

type ExpoNotificationStatusRequest struct {
	Ids []string `json:"ids"`
}

// SendNotifications Send a notification to the expo api and serialize response
func SendNotifications(notifications []ExpoNotification) (*ExpoNotificationSendResponse, error) {
	marshalledNotification, err := json.Marshal(notifications)
	rlog.Debugf("Marshalled notification into payload: %s\n", marshalledNotification)
	if err != nil {
		return nil, err
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
	rlog.Debug("Sending Notification to Expo")
	resp, err := client.Do(req)

	if err != nil {
		rlog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		rlog.Error(err)
		return nil, err
	}
	rlog.Debugf("Successfully sent notification to clients\n")

	var res ExpoNotificationSendResponse
	err = json.Unmarshal(bodyBytes, &res)
	rlog.Debugf("Got response from expo: %v", string(bodyBytes))

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetNotificationStatus Get the status on expo for the notification wrt it being delivered to apple or google.
func GetNotificationStatus(notificationIds []string) (*ExpoNotificationStatus, error) {
	reqBody, err := json.Marshal(&ExpoNotificationStatusRequest{notificationIds})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s%s", EXPO_HOST, API_URL, PUSH_RECEIPT_STATUS),
		bytes.NewBuffer(reqBody),
	)
	client := &http.Client{}
	req.Header.Add("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// cleanup
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var res ExpoNotificationStatus
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &res, nil
}
