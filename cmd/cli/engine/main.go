package main

import (
	"fmt"

	loadsim "github.com/dave-malone/aws-iot-loadsimulator/pkg"
)

const (
	//TODO - externalize this config
	sns_topic_arn string = "arn:aws:sns:us-east-1:068311527115:iot_simulator_notifications"
	one_million   int    = 1000000
	one_thousand  int    = 1000
)

func main() {
	fmt.Println("Running aws-iot-loadsimulator engine")

	config := &loadsim.SnsEventEngineConfig{
		TargetTotalConcurrentThings: one_thousand * 10,
		ClientsPerWorker:            one_thousand,
		MessagesToGeneratePerClient: 10,
		AwsRegion:                   "us-east-1",
		AwsSnsTopicARN:              sns_topic_arn,
		SecondsBetweenEachEvent:     10,
	}

	engine := loadsim.NewSnsEventEngine(config)
	if err := engine.GenerateEvents(); err != nil {
		fmt.Printf("Failed to generate events: %v", err)
		return
	}

	fmt.Println("Simulation requests generated")
}
