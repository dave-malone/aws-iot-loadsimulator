package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	loadsim "github.com/dave-malone/aws-iot-loadsimulator/pkg"
)

const (
	sns_topic_arn string = "arn:aws:sns:us-east-1:068311527115:iot_simulator_notifications"
	one_million   int    = 1000000
	one           int    = 1
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

	clientCount := 1000
	//maxClients := one_million
	maxClients := one

	for i := 0; i < (maxClients / clientCount); i++ {
		simRequest := &loadsim.SimulationRequest{
			StartClientNumber: (i * clientCount),
			ClientCount:       clientCount,
			MessagesPerClient: 250,
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

		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}

	return fmt.Sprintf("Simulation requests generated"), nil
}
