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
	TargetTotalConcurrentThings int `json:"total-things"`
	ClientsPerWorker            int `json:"clients-per-worker"`
	MessagesToGeneratePerClient int `json:"total-messages-per-client"`
	AwsRegion                   string
	AwsSnsTopicARN              string
	SecondsBetweenEachEvent     int `json:"seconds-between-sns-events"`
	SecondsBetweenMessages      int `json:"seconds-between-mqtt-messages"`
}

func (s SnsEventEngineConfig) String() string {
	return fmt.Sprintf(`SnsEventEngineConfig
		TargetTotalConcurrentThings: %d
		MessagesToGeneratePerClient: %d
		SecondsBetweenEachEvent: %d
		SecondsBetweenMessages: %d
		ClientsPerWorker: %d`,
		s.TargetTotalConcurrentThings,
		s.MessagesToGeneratePerClient,
		s.SecondsBetweenEachEvent,
		s.SecondsBetweenMessages,
		s.ClientsPerWorker,
	)
}

type SnsEventEngine struct {
	SnsEventEngineConfig
}

func NewSnsEventEngine(config *SnsEventEngineConfig) *SnsEventEngine {
	return &SnsEventEngine{
		SnsEventEngineConfig: *config,
	}
}

func (e *SnsEventEngine) GenerateEvents() (int, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(e.AwsRegion),
	})

	if err != nil {
		return 0, fmt.Errorf("Failed to initialize New AWS Session: %v", err)
	}

	client := sns.New(sess)

	totalWorkers := e.TargetTotalConcurrentThings / e.ClientsPerWorker

	fmt.Printf("targetTotalConcurrentThings: %d\n", e.TargetTotalConcurrentThings)
	fmt.Printf("clientsPerWorker: %d\n", e.ClientsPerWorker)
	fmt.Printf("totalWorkers: %d\n", totalWorkers)

	start := time.Now()

	for clientId := 0; clientId < totalWorkers; clientId++ {
		simRequest := &SimulationRequest{
			ClientId:               clientId,
			ClientCount:            e.ClientsPerWorker,
			StartClientNumber:      clientId * e.ClientsPerWorker,
			SecondsBetweenMessages: e.SecondsBetweenMessages,
			MessagesPerClient:      e.MessagesToGeneratePerClient,
		}

		messagePayload, err := json.Marshal(simRequest)
		if err != nil {
			return 0, fmt.Errorf("Failed to marshall simulation request payload: %v", err)
		}

		input := &sns.PublishInput{
			Message:  aws.String(string(messagePayload)),
			TopicArn: aws.String(e.AwsSnsTopicARN),
		}

		result, err := client.Publish(input)
		if err != nil {
			fmt.Printf("SNS Publish error: %v\n", err)
		}

		fmt.Printf("Simulation Request: %v\nSNS publish result: %v\n", simRequest, result)
		sleepTime := time.Duration(e.SecondsBetweenEachEvent) * time.Second
		fmt.Printf("Sleeping %.f seconds between publishing SNS notifications\n", sleepTime.Seconds())
		time.Sleep(sleepTime)
	}

	executionDuration := time.Since(start)
	fmt.Printf("Simulation Requests generated. Total Execution Time: %v\n", executionDuration)

	return totalWorkers, nil
}
