package main

import (
	"fmt"

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
	fmt.Println("Running aws-iot-loadsimulator worker")

	event := &loadsim.SimulationRequest{
		ClientCount:       1000,
		MessagesPerClient: 10,
		StartClientNumber: 0,
	}

	config := &loadsim.Config{
		CertificatePath:      certPath,
		PrivateKeyPath:       privateKeyPath,
		RootCAPath:           rootCAPath,
		MqttHost:             host,
		MqttPort:             port,
		MaxConcurrentClients: maxConcurrentClients,
		ClientIDPrefix:       clientIDPrefix,
		TopicPrefix:          topicPrefix,
	}

	worker, err := loadsim.NewWorker(config)
	if err != nil {
		fmt.Printf("Failed to initialize worker: %v", err)
		return
	}

	result, err := worker.RunConcurrentlyPublishingClients(event)
	if err != nil {
		fmt.Printf("Failed to run %d concurrent clients: %v", maxConcurrentClients, err)
		return
	}

	fmt.Println(result)

}
