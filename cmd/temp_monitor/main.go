package main

import (
	"temp_monitor/internal/app/temp_monitor"

	"github.com/aws/aws-lambda-go/lambda"
)

var lambdaStart = lambda.Start

func main() {
	lambdaStart(temp_monitor.HandleRequest)
}
