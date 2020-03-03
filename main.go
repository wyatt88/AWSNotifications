package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	//"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/awserr"
	//"github.com/aws/aws-sdk-go/aws/session"
)


func handler(ctx context.Context, event events.SNSEvent) {
	if len(event.Records) > 0 {
		for _, record := range event.Records {
			fmt.Print(record.SNS.Message)
		}
	}
}

func main() {
	lambda.Start(handler)
}