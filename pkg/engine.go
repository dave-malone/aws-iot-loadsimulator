package awsiotloadsimulator

import (
	"encoding/json"
	"fmt"

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
	SecondsBetweenEachEvent     int
}

type SnsEventEngine struct {
	config *SnsEventEngineConfig
}

func (e *SnsEventEngine) GenerateEvents() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(e.config.AwsRegion),
	})

	if err != nil {
		return fmt.Errorf("Failed to initialize New AWS Session: %v", err)
	}

	client := sns.New(sess)

	targetTotalConcurrentThings := e.config.TargetTotalConcurrentThings
	clientsPerWorker := e.config.ClientsPerWorker
	totalWorkers := targetTotalConcurrentThings / clientsPerWorker

	fmt.Printf("targetTotalConcurrentThings: %d\n", targetTotalConcurrentThings)
	fmt.Printf("clientsPerWorker: %d\n", clientsPerWorker)
	fmt.Printf("totalWorkers: %d\n", totalWorkers)

	ConcurrentWorkerExecutor(totalWorkers, e.config.SecondsBetweenEachEvent, func(clientId int) error {
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
			TopicArn: aws.String(e.config.AwsSnsTopicARN),
		}

		result, err := client.Publish(input)
		if err != nil {
			return fmt.Errorf("SNS Publish error: %v", err)
		}

		fmt.Printf("Simulation Request: %v\nSNS publish result: %v\n", simRequest, result)

		return nil
	})

	return nil
}

func NewSnsEventEngine(config *SnsEventEngineConfig) *SnsEventEngine {
	return &SnsEventEngine{
		config: config,
	}
}
