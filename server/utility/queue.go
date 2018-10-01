package utility

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/romana/rlog"
)

type SQSMock struct {
	listeners map[string][]func(*events.SQSEvent)
}

type SQSQueue interface {
	SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}

func (s SQSMock) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	queueUrl := input.QueueUrl
	message := input.MessageBody
	var handlers []func(*events.SQSEvent)
	var ok bool
	if handlers, ok = s.listeners[*queueUrl]; !ok {
		rlog.Warnf("No subscribers for queue %s", queueUrl)
	}

	for _, handler := range handlers {
		rlog.Info("Sending message to handler")
		handler(messageToSQSEvent(message))
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

func (s SQSMock) SubscribeListener(url string, handler func(*events.SQSEvent)) {
	_, ok := s.listeners[url]
	if !ok {
		s.listeners[url] = make([]func(*events.SQSEvent), 0, 1)

	}
	s.listeners[url] = append(s.listeners[url], handler)
}

func CreateMockSQSClient() SQSMock {
	rlog.Info("Initializing mock SQS")
	return SQSMock{
		listeners: make(map[string][]func(*events.SQSEvent)),
	}
}
