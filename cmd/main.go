package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func hello(_ events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Hello ƛ!",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(hello)
}
