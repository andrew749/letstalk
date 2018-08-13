package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Message string `json:"message"`
}

type Response struct {
	Message string `json:"message"`
}

func HandleRequest(ctx context.Context, data Request) (*Response, error) {
	fmt.Printf("data%#v\n", data)
	response := Response{
		Message: data.Message,
	}
	return &response, nil
}

func main() {
	lambda.Start(HandleRequest)
}
