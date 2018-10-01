package utility

import (
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/romana/rlog"
)

type queueMessage struct {
	dest    string
	payload interface{}
	retry   int
}

type SQSMock struct {
	listeners  map[string][]func(*events.SQSEvent) error
	eventQueue chan queueMessage
}

type SQSQueue interface {
	SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}

func (s SQSMock) QueueProcessor() {
	for {
		select {
		case m := <-s.eventQueue:
			var handlers []func(*events.SQSEvent) error
			var ok bool
			if handlers, ok = s.listeners[m.dest]; !ok {
				rlog.Warnf("No subscribers for queue %s", m.dest)
			}

			time.Sleep(time.Duration(1) * time.Second)

			for _, handler := range handlers {
				rlog.Info("Sending message to handler")
				if err := handler(messageToSQSEvent(m.payload.(*string))); err != nil {
					rlog.Error("Error processing event: %v", err)
					// requeue message
					m.retry -= 1
					if m.retry >= 0 {
						rlog.Debug("Retrying event")
						s.eventQueue <- m
					}
				}
			}
			break
		}
	}
}

func (s SQSMock) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	queueUrl := input.QueueUrl
	message := input.MessageBody
	s.eventQueue <- queueMessage{
		dest:    *queueUrl,
		payload: message,
		retry:   3,
	}

	return nil, nil
}

// TODO make this work for other sqs queues
func messageToSQSEvent(message *string) *events.SQSEvent {
	return &events.SQSEvent{
		Records: []events.SQSMessage{
			events.SQSMessage{
				Body: *message,
			},
		},
	}
}

func (s SQSMock) SubscribeListener(url string, handler func(*events.SQSEvent) error) {
	_, ok := s.listeners[url]
	if !ok {
		s.listeners[url] = make([]func(*events.SQSEvent) error, 0, 1)

	}
	s.listeners[url] = append(s.listeners[url], handler)
	rlog.Infof("Subscribing listener for url %s", url)
}

func CreateMockSQSClient() SQSMock {
	rlog.Info("Initializing mock SQS")
	return SQSMock{
		listeners:  make(map[string][]func(*events.SQSEvent) error),
		eventQueue: make(chan queueMessage, 10),
	}
}
