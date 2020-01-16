package awsiotloadsimulator

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/dave-malone/aws-iot-loadsimulator/pkg/mqtt"
)

type WorkerConfig struct {
	CertificatePath                string
	PrivateKeyPath                 string
	RootCAPath                     string
	MqttHost                       string
	MqttPort                       int
	MaxConnectionRequestsPerSecond time.Duration
	ClientIDPrefix                 string
	TopicPrefix                    string
}

type Worker struct {
	WorkerConfig
	tlsConfig *tls.Config
}

func NewWorker(config *WorkerConfig) (*Worker, error) {
	tlsConfig, err := mqtt.NewTlsConfig(config.CertificatePath, config.PrivateKeyPath, config.RootCAPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize tls.Config: %v", err)
	}

	worker := &Worker{
		WorkerConfig: *config,
		tlsConfig:    tlsConfig,
	}

	return worker, nil
}

func (w *Worker) RunConcurrentlyPublishingClients(simReq *SimulationRequest) (string, error) {
	fmt.Printf("Initializing %d concurrent mqtt clients\n", simReq.ClientCount)
	executionDuration := ConcurrentWorkerExecutor(simReq.ClientCount, w.MaxConnectionRequestsPerSecond, func(thingId int) error {
		clientId := simReq.StartClientNumber + thingId
		if err := w.publishMessages(clientId, simReq); err != nil {
			fmt.Printf("Failed to publish messages: %v\n", err)
		}

		return nil
	})

	return fmt.Sprintf("Simulation complete. Total Execution time: %v", executionDuration), nil
}

func (w *Worker) publishMessages(clientNumber int, simReq *SimulationRequest) error {
	clientID := fmt.Sprintf("%s-%d", w.ClientIDPrefix, clientNumber)
	//TODO - allow for topic to be injected with support for placeholders
	topic := fmt.Sprintf("%s/%s", w.TopicPrefix, clientID)

	mqttClient := mqtt.NewClient(w.MqttHost, w.MqttPort, clientID, w.tlsConfig)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("[%s] Failed to get connection token: %s", clientID, token.Error().Error())
	}

	log.Printf("[%s] Connected\n", clientID)

	for messageNumber := 1; messageNumber <= simReq.MessagesPerClient; messageNumber++ {
		//TODO - allow for the message to be injected from SimulationRequest with support for dynamic ranges
		payload := map[string]interface{}{
			"message_number": messageNumber,
			"client_id":      clientID,
			"timestamp":      time.Now().Format(time.RFC3339),
		}

		if err := mqttClient.PublishAsJson(payload, topic, 0); err != nil {
			log.Printf("[%s] %v", clientID, err)
		}

		log.Printf("[%s] Successfully published message %v\n", clientID, payload)
		sleepFor := time.Duration(simReq.SecondsBetweenMessages) * time.Second
		log.Printf("[%s] sleeping %.f seconds", clientID, sleepFor.Seconds())
		time.Sleep(sleepFor)
	}

	mqttClient.Disconnect(1000)
	return nil
}
