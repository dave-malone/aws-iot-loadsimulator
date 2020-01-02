package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	loadsim "github.com/dave-malone/aws-iot-loadsimulator/pkg"
)

const (
	//TODO - externalize config
	certPath             = "./certs/golang_thing.cert.pem"
	privateKeyPath       = "./certs/golang_thing.private.key"
	rootCAPath           = "./certs/root-CA.crt"
	host                 = "a1tq0bx5we8tnk-ats.iot.us-east-1.amazonaws.com"
	port                 = 8883
	clientIDPrefix       = "golang_thing"
	topicPrefix          = "golang_simulator"
	maxConcurrentClients = 1000
)

func main() {
	lambda.Start(requestHandler)
}

func requestHandler(ctx context.Context, snsEvent events.SNSEvent) (string, error) {
	if len(snsEvent.Records) != 1 {
		return "", fmt.Errorf("snsEvent.Records expected to be exactly 1. Length was %d", len(snsEvent.Records))
	}

	snsRecord := snsEvent.Records[0].SNS
	var event loadsim.SimulationRequest

	log.Printf("sns message body: %s\n", snsRecord.Message)

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
		return "", fmt.Errorf("Failed to run %d concurrent clients: %v", maxConcurrentClients, err)
	}

	fmt.Println(result)
	return result, nil
}
