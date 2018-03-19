package aws_utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"letstalk/server/core/secrets"
)

func getDefaultConfiguration() *aws.Config {
	return &aws.Config{Region: aws.String(endpoints.UsEast1RegionID)}
}

func getDefaultSession() (*session.Session, error) {
	sess, err := session.NewSession(getDefaultConfiguration())
	// currently just use the default credentials
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func getAWSConfig() (*aws.Config, error) {
	creds := credentials.NewStaticCredentials(
		secrets.GetSecrets().DefaultAccessKeyID,
		secrets.GetSecrets().DefaultAccessKeySecret,
		"",
	)
	_, err := creds.Get()

	if err != nil {
		return nil, err
	}
	config := &aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String("s3.amazonaws.com"),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
		//LogLevel        :0,
	}
	return config, nil
}

func GetSNSServiceClient() (*sns.SNS, error) {
	sess, err := getDefaultSession()
	if err != nil {
		return nil, err
	}
	snsClient := sns.New(sess)
	return snsClient, nil
}

func GetLambdaServiceClient() (*lambda.Lambda, error) {
	sess, err := getDefaultSession()
	if err != nil {
		return nil, err
	}

	lambdaClient := lambda.New(sess)
	return lambdaClient, nil
}

func GetSQSServiceClient() (*sqs.SQS, error) {
	sess, err := getDefaultSession()
	if err != nil {
		return nil, err
	}
	sqsClient := sqs.New(sess)
	return sqsClient, nil
}

func GetS3ServiceClient() (*s3.S3, error) {
	sess, err := getDefaultSession()
	if err != nil {
		return nil, err
	}
	config, err := getAWSConfig()

	if err != nil {
		return nil, err
	}

	s3Client := s3.New(sess, config)

	return s3Client, nil
}
