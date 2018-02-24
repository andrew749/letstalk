package aws_utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/sns"
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
