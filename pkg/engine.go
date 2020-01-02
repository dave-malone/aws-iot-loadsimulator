package awsiotloadsimulator

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SnsEventEngineConfig struct {
	TargetTotalConcurrentThings int
	ClientsPerWorker            int
	MessagesToGeneratePerClient int
	AwsRegion                   string
	AwsSnsTopicARN              string
	SecondsBetweenEachEvent     time.Duration
}

type SnsEventEngine struct {
	SnsEventEngineConfig
}

func (e *SnsEventEngine) GenerateEvents() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(e.AwsRegion),
	})

	if err != nil {
		return fmt.Errorf("Failed to initialize New AWS Session: %v", err)
	}

	client := sns.New(sess)

	targetTotalConcurrentThings := e.TargetTotalConcurrentThings
	clientsPerWorker := e.ClientsPerWorker
	totalWorkers := targetTotalConcurrentThings / clientsPerWorker

	fmt.Printf("targetTotalConcurrentThings: %d\n", targetTotalConcurrentThings)
	fmt.Printf("clientsPerWorker: %d\n", clientsPerWorker)
	fmt.Printf("totalWorkers: %d\n", totalWorkers)

	executionDuration := ConcurrentWorkerExecutor(totalWorkers, e.SecondsBetweenEachEvent, func(clientId int) error {
		simRequest := &SimulationRequest{
			ClientId:    clientId,
			ClientCount: clientsPerWorker,
		}

		messagePayload, err := json.Marshal(simRequest)
		if err != nil {
			return fmt.Errorf("Failed to marshall simulation request payload: %v", err)
		}

		input := &sns.PublishInput{
			Message:  aws.String(string(messagePayload)),
			TopicArn: aws.String(e.AwsSnsTopicARN),
		}

		result, err := client.Publish(input)
		if err != nil {
			return fmt.Errorf("SNS Publish error: %v", err)
		}

		fmt.Printf("Simulation Request: %v\nSNS publish result: %v\n", simRequest, result)

		return nil
	})

	fmt.Printf("Simulation Requests generated. Total Execution Time: %v\n", executionDuration)

	return nil
}

func NewSnsEventEngine(config *SnsEventEngineConfig) *SnsEventEngine {
	return &SnsEventEngine{
		SnsEventEngineConfig: *config,
	}
}
