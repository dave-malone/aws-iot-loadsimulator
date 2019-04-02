package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	loadsim "github.com/dave-malone/aws-iot-loadsimulator/pkg"
)

const (
	//TODO - externalize this config
	sns_topic_arn string = "arn:aws:sns:us-east-1:068311527115:iot_simulator_notifications"
	one_million   int    = 1000000
	one_thousand  int    = 1000
)

func main() {
	lambda.Start(requestHandler)
}

func requestHandler(ctx context.Context) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		return "", fmt.Errorf("Failed to initialize New AWS Session: %v", err)
	}

	client := sns.New(sess)

	targetTotalConcurrentThings := one_thousand * 10
	clientsPerWorker := one_thousand
	totalWorkers := targetTotalConcurrentThings / clientsPerWorker

	fmt.Printf("targetTotalConcurrentThings: %d\n", targetTotalConcurrentThings)
	fmt.Printf("clientsPerWorker: %d\n", clientsPerWorker)
	fmt.Printf("totalWorkers: %d\n", totalWorkers)

	for i := 0; i < totalWorkers; i++ {
		simRequest := &loadsim.SimulationRequest{
			StartClientNumber: (i * clientsPerWorker),
			ClientCount:       clientsPerWorker,
			MessagesPerClient: 10,
		}

		messagePayload, err := json.Marshal(simRequest)
		if err != nil {
			return "", fmt.Errorf("Failed to marshall simulation request payload: %v", err)
		}

		input := &sns.PublishInput{
			Message:  aws.String(string(messagePayload)),
			TopicArn: aws.String(sns_topic_arn),
		}

		result, err := client.Publish(input)
		if err != nil {
			return "", fmt.Errorf("Publish error: %v", err)
		}

		fmt.Printf("SNS publish result: %v\n", result)

		time.Sleep(time.Duration(10) * time.Second)
	}

	return fmt.Sprintf("%d Simulation requests generated", totalWorkers), nil
}
