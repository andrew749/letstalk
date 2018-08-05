package queue

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/sqs"
)

func CreateNewQueue(sess *sqs.SQS, queueId *string) (*sqs.CreateQueueOutput, error) {
	yes := "yes"
	attributes := map[string]*string{
		"ContentBasedDeduplication": &yes,
	}

	queueInput := &sqs.CreateQueueInput{
		Attributes: attributes,
		QueueName:  queueId,
	}
	output, err := sess.CreateQueue(queueInput)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func GetQueueLocation(sess *sqs.SQS, queueName *string) (*string, error) {
	input := &sqs.GetQueueUrlInput{
		QueueName: queueName,
	}

	out, err := sess.GetQueueUrl(input)
	if err != nil {
		return nil, err
	}

	return out.QueueUrl, nil
}

func AddNewMessage(
	sess *sqs.SQS,
	queueId string,
	queueUrl string,
	payload interface{},
) (*sqs.SendMessageOutput, error) {
	marshalledPayload, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	stringPayload := string(marshalledPayload)
	input := &sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: &stringPayload,
	}
	output, err := sess.SendMessage(input)
	if err != nil {
		return nil, err
	}
	return output, nil
}
