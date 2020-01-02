package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	loadsim "github.com/dave-malone/aws-iot-loadsimulator/pkg"
)

const (
	one_million  int = 1000000
	one_thousand int = 1000
)

func main() {
	lambda.Start(requestHandler)
}

func requestHandler(ctx context.Context) (string, error) {
	sns_topic_arn := os.Getenv("SNS_TOPIC_ARN")
	if len(sns_topic_arn) == 0 {
		return "", fmt.Errorf("Environment variable SNS_TOPIC_ARN not set")
	}

	config := &loadsim.SnsEventEngineConfig{
		TargetTotalConcurrentThings: 131072,
		ClientsPerWorker:            1,
		MessagesToGeneratePerClient: 10,
		AwsRegion:                   os.Getenv("AWS_REGION"),
		AwsSnsTopicARN:              sns_topic_arn,
		SecondsBetweenEachEvent:     0,
	}

	engine := loadsim.NewSnsEventEngine(config)
	if err := engine.GenerateEvents(); err != nil {
		return "", fmt.Errorf("Failed to generate events: %v", err)
	}

	return fmt.Sprintf("%d Simulation requests generated", config.TargetTotalConcurrentThings), nil
}
