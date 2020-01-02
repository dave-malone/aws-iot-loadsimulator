package main

import (
	"flag"
	"fmt"

	loadsim "github.com/dave-malone/aws-iot-loadsimulator/pkg"
)

var mode = flag.String("mode", "init", "Executes this program in one of two modes: 'init' or 'cleanup'")
var total_number_of_things = flag.Int("totalThings", 100, "[Optional] Total Number of things to generate in the thing registry")
var max_requests_per_second = flag.Int("maxRequestsPerSecond", 15, "[Optional] Maximum number of IoT API requests per second")
var aws_region = flag.String("region", "us-east-1", "[Optional] set the target AWS region")

func main() {
	flag.Parse()

	fmt.Printf("%s AWS IoT Device Registry for Device Simulation\n", *mode)

	config := &loadsim.DeviceRegistryConfig{
		AwsRegion:            *aws_region,
		ThingNamePrefix:      "golang_thing",
		ThingTypeName:        "simulated-thing",
		MaxRequestsPerSecond: *max_requests_per_second,
		TotalNumberOfThings:  *total_number_of_things,
	}

	deviceRegistry := loadsim.NewDeviceRegistry(config)

	switch *mode {
	case "init":
		if err := deviceRegistry.Initialize(); err != nil {
			fmt.Printf("Failed to initialize device registry: %v\n", err)
			return
		}

		fmt.Println("Device Registry Initialized")
	case "cleanup":
		if err := deviceRegistry.Cleanup(); err != nil {
			fmt.Printf("Failed to cleanup device registry: %v\n", err)
			return
		}

		fmt.Println("Device Registry Cleaned Up")
	default:
		fmt.Printf("%s is not a valid mode; please use either 'init' or 'cleanup'\n", *mode)
	}

}
