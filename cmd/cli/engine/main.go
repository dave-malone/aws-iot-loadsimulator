package main

import (
	"flag"
	"fmt"
	"os"

	loadsim "github.com/dave-malone/aws-iot-loadsimulator/pkg"
)

var total_number_of_things = flag.Int("totalThings", 1000, "[Optional] Total Number of things to generate in the thing registry")
var clients_per_worker = flag.Int("clientsPerWorker", 300, "[Optional] Maximum number of concurrent clients per worker")
var max_requests_per_second = flag.Int("maxRequestsPerSecond", 15, "[Optional] Maximum number of IoT API requests per second")
var aws_region = flag.String("region", "us-east-1", "[Optional] set the target AWS region")

func main() {
	flag.Parse()

	fmt.Println("Running aws-iot-loadsimulator engine")

	sns_topic_arn := os.Getenv("SNS_TOPIC_ARN")
	if len(sns_topic_arn) == 0 {
		fmt.Println("Environment variable SNS_TOPIC_ARN not set")
		return
	}

	config := &loadsim.SnsEventEngineConfig{
		TargetTotalConcurrentThings: *total_number_of_things,
		ClientsPerWorker:            *clients_per_worker,
		MessagesToGeneratePerClient: 10,
		AwsRegion:                   *aws_region,
		AwsSnsTopicARN:              sns_topic_arn,
		SecondsBetweenEachEvent:     *max_requests_per_second,
	}

	engine := loadsim.NewSnsEventEngine(config)
	if err := engine.GenerateEvents(); err != nil {
		fmt.Printf("Failed to generate events: %v", err)
		return
	}

	fmt.Println("Simulation requests generated")
}
