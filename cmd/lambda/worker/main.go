package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	loadsim "github.com/dave-malone/aws-iot-loadsimulator/pkg"
)

const (
	certPath       = "./certs/golang_thing.cert.pem"
	privateKeyPath = "./certs/golang_thing.private.key"
	rootCAPath     = "./certs/root-CA.crt"
	port           = 8883
	//TODO - these should be part of the loadsim.SimulationRequest, as well as the client message body
	clientIDPrefix = "golang_thing"
	topicPrefix    = "golang_simulator"
)

func main() {
	lambda.Start(requestHandler)
}

func requestHandler(ctx context.Context, snsEvent events.SNSEvent) (string, error) {
	log.Printf("Received SNS Event: %v\n", snsEvent)

	host := os.Getenv("MQTT_HOST")
	if len(host) == 0 {
		return "", fmt.Errorf("Environment variable MQTT_HOST must be set")
	}

	if len(snsEvent.Records) != 1 {
		return "", fmt.Errorf("snsEvent.Records expected to be exactly 1. Length was %d", len(snsEvent.Records))
	}

	snsRecord := snsEvent.Records[0].SNS
	log.Printf("sns message body: %s\n", snsRecord.Message)

	var event loadsim.SimulationRequest
	if err := json.Unmarshal([]byte(snsRecord.Message), &event); err != nil {
		return "", fmt.Errorf("Failed to unmarshall sns message body: %v", err)
	}

	config := &loadsim.WorkerConfig{
		CertificatePath:                certPath,
		PrivateKeyPath:                 privateKeyPath,
		RootCAPath:                     rootCAPath,
		MqttHost:                       host,
		MqttPort:                       port,
		MaxConnectionRequestsPerSecond: 100,
		ClientIDPrefix:                 clientIDPrefix,
		TopicPrefix:                    topicPrefix,
	}

	worker, err := loadsim.NewWorker(config)
	if err != nil {
		return "", fmt.Errorf("Failed to initialize worker: %v", err)
	}

	result, err := worker.RunConcurrentlyPublishingClients(&event)
	if err != nil {
		return "", fmt.Errorf("Failed to run %d concurrent clients: %v", event.ClientCount, err)
	}

	fmt.Println(result)
	return result, nil
}
