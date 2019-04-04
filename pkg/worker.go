package awsiotloadsimulator

import (
	"crypto/tls"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dave-malone/aws-iot-loadsimulator/pkg/mqtt"
)

type Config struct {
	CertificatePath      string
	PrivateKeyPath       string
	RootCAPath           string
	MqttHost             string
	MqttPort             int
	MaxConcurrentClients int
	ClientIDPrefix       string
	TopicPrefix          string
}

type Worker struct {
	config    *Config
	tlsConfig *tls.Config
}

func NewWorker(config *Config) (*Worker, error) {
	tlsConfig, err := mqtt.NewTlsConfig(config.CertificatePath, config.PrivateKeyPath, config.RootCAPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize tls.Config: %v", err)
	}

	worker := &Worker{
		config:    config,
		tlsConfig: tlsConfig,
	}

	return worker, nil
}

func (w *Worker) RunConcurrentlyPublishingClients(simReq *SimulationRequest) (string, error) {
	var wg sync.WaitGroup
	sem := make(chan int, w.config.MaxConcurrentClients)
	start := time.Now()

	clientNumberMax := simReq.StartClientNumber + simReq.ClientCount

	for clientNumber := simReq.StartClientNumber; clientNumber < clientNumberMax; clientNumber++ {
		wg.Add(1)
		go func(clientNumber int, simReq *SimulationRequest) {
			sem <- 1
			go func(clientNumber int, simReq *SimulationRequest) {
				defer wg.Done()
				if err := w.publishMessages(clientNumber, simReq); err != nil {
					fmt.Printf("Failed to publish messages: %v\n", err)
				}
				<-sem
			}(clientNumber, simReq)
		}(clientNumber, simReq)
	}

	wg.Wait()
	elapsed := time.Since(start)
	return fmt.Sprintf("simulation complete; execution time: %s", elapsed), nil
}

func (w *Worker) publishMessages(clientNumber int, simReq *SimulationRequest) error {
	clientID := fmt.Sprintf("%s-%d", w.config.ClientIDPrefix, clientNumber)
	//TODO - allow for topic to be injected with support for placeholders
	topic := fmt.Sprintf("%s/%s", w.config.TopicPrefix, clientID)

	mqttClient := mqtt.NewClient(w.config.MqttHost, w.config.MqttPort, clientID, w.tlsConfig)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("[%s] Failed to get connection token: %v", clientID, token.Error())
	}

	log.Printf("[%s] Connected\n", clientID)

	for messageNumber := 1; messageNumber <= simReq.MessagesPerClient; messageNumber++ {
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
