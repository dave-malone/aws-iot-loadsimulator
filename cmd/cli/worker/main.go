package main

import (
	"flag"
	"fmt"

	loadsim "github.com/dave-malone/aws-iot-loadsimulator/pkg"
)

var (
	certPath               = flag.String("cert", "./certs/golang_thing.cert.pem", "path to the certificate")
	privateKeyPath         = flag.String("private-key", "./certs/golang_thing.private.key", "path to the private key")
	rootCAPath             = flag.String("root-ca", "./certs/root-CA.crt", "path to the root certificate authority")
	host                   = flag.String("host", "", "mqtt host number")
	port                   = flag.Int("port", 8883, "mqtt port number")
	clientIDPrefix         = flag.String("client-id-prefix", "golang_thing", "prefix to give each mqtt client")
	topicPrefix            = flag.String("topic-prefix", "golang_simulator", "prefix to use the mqtt topic used by each client")
	maxConcurrentClients   = flag.Int("max-clients", 10, "maximum number of mqtt clients to run")
	messagesPerClient      = flag.Int("total-messages-per-client", 10, "messages to generate per client")
	secondsBetweenMessages = flag.Int("seconds-between-sns-messages", 1, "number of seconds to wait between publish calls per client")
)

func main() {
	fmt.Println("Running aws-iot-loadsimulator worker")

	flag.Parse()

	if len(*host) == 0 {
		fmt.Println("host flag must be set")
		return
	}

	config := &loadsim.WorkerConfig{
		CertificatePath:                *certPath,
		PrivateKeyPath:                 *privateKeyPath,
		RootCAPath:                     *rootCAPath,
		MqttHost:                       *host,
		MqttPort:                       *port,
		MaxConnectionRequestsPerSecond: 100,
		ClientIDPrefix:                 *clientIDPrefix,
		TopicPrefix:                    *topicPrefix,
	}

	worker, err := loadsim.NewWorker(config)
	if err != nil {
		fmt.Printf("Failed to initialize worker: %v", err)
		return
	}

	event := &loadsim.SimulationRequest{
		ClientCount:            *maxConcurrentClients,
		MessagesPerClient:      *messagesPerClient,
		SecondsBetweenMessages: *secondsBetweenMessages,
		StartClientNumber:      0,
	}

	result, err := worker.RunConcurrentlyPublishingClients(event)
	if err != nil {
		fmt.Printf("Failed to run %d concurrent clients: %v", maxConcurrentClients, err)
		return
	}

	fmt.Println(result)

}
