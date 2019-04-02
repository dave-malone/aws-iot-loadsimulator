package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	loadsim "github.com/dave-malone/aws-iot-loadsimulator/pkg"
	mqtt "github.com/dave-malone/aws-iot-loadsimulator/pkg/mqtt"
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

	concurrencyLimit := maxConcurrentClients
	clientNumberMax := event.StartClientNumber + event.ClientCount
	messagesToPublishPerClient := event.MessagesPerClient

	result, err := runConcurrentlyPublishingClients(event.StartClientNumber, clientNumberMax, messagesToPublishPerClient, concurrencyLimit)
	if err != nil {
		return "", fmt.Errorf("Failed to run %d concurrent clients: %v", concurrencyLimit, err)
	}

	return result, nil
}

func runConcurrentlyPublishingClients(startClientNumber, clientNumberMax, messagesToPublishPerClient, concurrencyLimit int) (string, error) {
	tlsConfig, err := mqtt.NewTlsConfig(certPath, privateKeyPath, rootCAPath)
	if err != nil {
		return "", fmt.Errorf("Failed to initialize tls.Config: %v", err)
	}

	var wg sync.WaitGroup
	sem := make(chan int, concurrencyLimit)
	start := time.Now()

	for clientNumber := startClientNumber; clientNumber < clientNumberMax; clientNumber++ {
		wg.Add(1)
		go func(clientNumber int, messagesToPublishPerClient int, tlsConfig *tls.Config) {
			sem <- 1
			go func(clientNumber int, messagesToPublishPerClient int, tlsConfig *tls.Config) {
				defer wg.Done()
				if err := publishMessages(clientNumber, messagesToPublishPerClient, tlsConfig); err != nil {
					fmt.Printf("Failed to publish messages: %v\n", err)
				}
				<-sem
			}(clientNumber, messagesToPublishPerClient, tlsConfig)
		}(clientNumber, messagesToPublishPerClient, tlsConfig)
	}

	wg.Wait()
	elapsed := time.Since(start)
	return fmt.Sprintf("simulation complete; execution time: %s", elapsed), nil
}

func publishMessages(clientNumber int, messageCount int, tlsConfig *tls.Config) error {
	clientID := fmt.Sprintf("%s-%d", clientIDPrefix, clientNumber)
	//TODO - allow for topic to be injected with support for placeholders
	topic := fmt.Sprintf("%s/%s", topicPrefix, clientID)

	mqttClient := mqtt.NewClient(host, port, clientID, tlsConfig)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("[%s] Failed to get connection token: %v", clientID, token.Error())
	}

	log.Printf("[%s] Connected\n", clientID)

	for messageNumber := 1; messageNumber <= messageCount; messageNumber++ {
		//TODO - allow for the message to be injected from SimulationRequest with support for dynamic ranges
		payload := map[string]interface{}{
			"message_number": messageNumber,
			"client_id":      clientID,
		}

		if err := mqttClient.PublishAsJson(payload, topic, 0); err != nil {
			log.Printf("[%s] %v", clientID, err)
		}

		log.Printf("[%s] Successfully published message %v\n", clientID, payload)
		time.Sleep(time.Duration(90) * time.Second)
	}

	mqttClient.Disconnect(1000)
	return nil
}
