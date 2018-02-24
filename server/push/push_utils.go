package push

import (
	"letstalk/server/aws_utils"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/romana/rlog"
)

func pushNotification(message string, subject string, arn string) error {
	client, err := aws_utils.GetSNSServiceClient()
	if err != nil {
		return err
	}

	publishOutput, err := client.Publish(
		&sns.PublishInput{
			Message:  &message,
			Subject:  &subject,
			TopicArn: &arn,
		},
	)
	if err != nil {
		return err
	}

	rlog.Debug("Sent message with id:", publishOutput.MessageId)
	return nil
}

func createEmailSubscriber(snsClient *sns.SNS, emailAddress string, topicArn string) error {
	protocol := "email"
	rlog.Debug(topicArn)
	subscription, err := snsClient.Subscribe(
		&sns.SubscribeInput{
			Endpoint: &emailAddress,
			Protocol: &protocol,
			TopicArn: &topicArn,
		},
	)

	if err != nil {
		return err
	}

	rlog.Debug("Succesfully subscribed: ", emailAddress, " ", subscription.SubscriptionArn)

	return nil
}

func SendDiagnosticEmailNotification(
	snsClient *sns.SNS,
	emailAddress string,
	message string,
	subject string,
	topicArn string,
) error {
	err := createEmailSubscriber(snsClient, emailAddress, topicArn)
	if err != nil {
		return err
	}
	// pushNotification(message, subject, topicArn)
	return nil
}
