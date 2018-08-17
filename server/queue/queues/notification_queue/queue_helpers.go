package notification_queue

import (
	"letstalk/server/notifications"
	"letstalk/server/queue"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/romana/rlog"
)

const (
	NotificationQueueID  = "Notifications"
	NotificationQueueUrl = "https://sqs.us-east-1.amazonaws.com/016267150191/Notifications"
)

func PushNotificationToQueue(sqs *sqs.SQS, notification notifications.Notification) error {
	rlog.Debugf("%#v", notification)
	_, err := queue.AddNewMessage(sqs, NotificationQueueID, NotificationQueueUrl, notification)
	return err
}
