package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	loadsim "github.com/dave-malone/aws-iot-loadsimulator/pkg"
)

func main() {
	lambda.Start(requestHandler)
}

func requestHandler(ctx context.Context, config loadsim.SnsEventEngineConfig) (string, error) {
	fmt.Printf("Received event: %v\n", config)

	snsTopicArn := os.Getenv("SNS_TOPIC_ARN")
	if len(snsTopicArn) == 0 {
		return "", fmt.Errorf("Environment variable SNS_TOPIC_ARN not set")
	}

	config.AwsRegion = os.Getenv("AWS_REGION")
	config.AwsSnsTopicARN = snsTopicArn

	if config.TargetTotalConcurrentThings < 1 {
		config.TargetTotalConcurrentThings = 1
	}

	if config.ClientsPerWorker < 1 {
		config.ClientsPerWorker = 1
	}

	if config.MessagesToGeneratePerClient < 1 {
		config.MessagesToGeneratePerClient = 10
	}

	if config.SecondsBetweenEachEvent < 1 {
		config.SecondsBetweenEachEvent = 5
	}

	if config.SecondsBetweenMessages < 1 {
		config.SecondsBetweenMessages = 10
	}

	fmt.Printf("Generating SNS Notifications using the following configuration: %v\n", config)

	engine := loadsim.NewSnsEventEngine(&config)
	notificationCount, err := engine.GenerateEvents()
	if err != nil {
		return "", fmt.Errorf("Failed to generate events: %v", err)
	}

	return fmt.Sprintf("%d Simulation requests generated", notificationCount), nil
}
