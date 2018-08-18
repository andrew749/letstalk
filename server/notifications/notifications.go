package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"letstalk/server/data"
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

type NotificationStatusDetails struct {
	Error string `json:"error"`
}

type NotificationStatusResponse struct {
	Id      string                     `json:"id,omitempty"`
	Status  string                     `json:"status"`
	Message *string                    `json:"message,omitempty"`
	Details *NotificationStatusDetails `json:"details,omitempty"`
}

type NotificationStatus struct {
	Data map[string]NotificationStatusResponse `json:"data"`
}

type NotificationSendResponse struct {
	Data []NotificationStatusResponse `json:"-"`
}

type NotificationStatusRequest struct {
	Ids []string `json:"ids"`
}

//UnmarshalJSON Custom unmarshalling since expo api could return an array or single item.
func (s *NotificationSendResponse) UnmarshalJSON(data []byte) error {
	res := struct {
		Data NotificationStatusResponse `json:"data"`
	}{}
	err := json.Unmarshal(data, &res)
	// if we were able to deserialize a single response
	if err == nil {
		s.Data = []NotificationStatusResponse{res.Data}
		return nil
	}
	s.Data = make([]NotificationStatusResponse, 0)

	return json.Unmarshal(data, &s.Data)
}

// FromNotificationDataModel Convert a notification data model to a version that the expo API expects
func (n *Notification) FromNotificationDataModel(orig data.Notification) *Notification {
	n.To = string(orig.UserId)
	n.Title = orig.Message
	n.Data = orig
	return n
}

// SendNotification Send a notification to the expo api and serialize response
func SendNotification(notification Notification) (*NotificationSendResponse, error) {
	marshalledNotification, err := json.Marshal(notification)
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
	rlog.Debug("Successfully sent notification to client: %s", notification.To)

	var res NotificationSendResponse
	err = json.Unmarshal(bodyBytes, &res)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetNotificationStatus Get the status on expo for the notification wrt it being delivered to apple or google.
func GetNotificationStatus(notificationIds []string) (*NotificationStatus, error) {
	reqBody, err := json.Marshal(&NotificationStatusRequest{notificationIds})
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

	var res NotificationStatus
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &res, nil
}
