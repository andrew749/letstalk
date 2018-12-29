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

type LocalQueueImpl struct {
	listeners  map[string][]func(*events.SQSEvent) error
	eventQueue chan queueMessage
	doneQueue  chan bool
}

type SQSQueue interface {
	SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
}

// WaitForQueueDone Wait until a queue is done processing (has been closed)
func (s LocalQueueImpl) WaitForQueueDone() {
	// block on channel
	<-s.doneQueue
	s.doneQueue <- true
}

func (s LocalQueueImpl) CloseQueue() {
	close(s.eventQueue)
}

func (s LocalQueueImpl) QueueProcessor() {
	for {
		select {
		case m, more := <-s.eventQueue:
			var handlers []func(*events.SQSEvent) error
			var ok bool
			if handlers, ok = s.listeners[m.dest]; !ok {
				rlog.Warnf("No subscribers for queue %s", m.dest)
			}

			time.Sleep(time.Duration(1) * time.Second)

			for _, handler := range handlers {
				rlog.Info("Sending message to handler")
				if err := handler(messageToSQSEvent(m.payload.(*string))); err != nil {
					rlog.Errorf("Error processing event: %+v", err)
					// requeue message
					m.retry -= 1
					if m.retry >= 0 {
						rlog.Debug("Retrying event")
						s.eventQueue <- m
					}
				}
			}

			// tell people that we are done.
			if !more {
				// exit goroutine
				rlog.Debug("Queue processor exiting.")
				s.doneQueue <- true
				return
			}
			break
		}
	}
	rlog.Debug("Queue processor exiting.")
}

func (s LocalQueueImpl) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
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

func (s LocalQueueImpl) SubscribeListener(url string, handler func(*events.SQSEvent) error) {
	_, ok := s.listeners[url]
	if !ok {
		s.listeners[url] = make([]func(*events.SQSEvent) error, 0, 1)

	}
	s.listeners[url] = append(s.listeners[url], handler)
	rlog.Infof("Subscribing listener for url %s", url)
}

func CreateLocalSQSClient() LocalQueueImpl {
	rlog.Info("Initializing local SQS")
	return LocalQueueImpl{
		listeners:  make(map[string][]func(*events.SQSEvent) error),
		eventQueue: make(chan queueMessage, 10),
		doneQueue:  make(chan bool, 1),
	}
}
