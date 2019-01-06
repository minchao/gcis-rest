package main

import (
	"github.com/minchao/gcis-rest/internal/app/company"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(company.Handler)
}
