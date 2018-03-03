package main

import (
	"letstalk/server/aws_utils"
	"letstalk/server/queue"
	"log"
)

func main() {
	sqs, err := aws_utils.GetSQSServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	message := struct{ Message string }{"andrews test message"}
	recipient := "recipient"
	queueName := "andrews_test_queue"
	queueUrl := "https://sqs.us-east-1.amazonaws.com/947945882937/andrews_test_queue"

	// queueOut, err := queue.CreateNewQueue(sqs, &queueName)
	// if err != nil {
	// log.Fatal(err)
	// }
	_, err = queue.AddNewMessage(sqs, &queueName, &queueUrl, &recipient, message)
	if err != nil {
		log.Fatal(err)
	}
}
