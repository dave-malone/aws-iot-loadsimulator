package awsiotloadsimulator

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
)

type DeviceRegistryConfig struct {
	AwsRegion            string
	ThingTypeName        string
	ThingNamePrefix      string
	MaxRequestsPerSecond int
	TotalNumberOfThings  int
}

type DeviceRegistry struct {
	DeviceRegistryConfig
	client *iot.IoT
}

func NewDeviceRegistry(config *DeviceRegistryConfig) *DeviceRegistry {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.AwsRegion),
	}))

	client := iot.New(sess)

	dr := &DeviceRegistry{
		DeviceRegistryConfig: *config,
		client:               client,
	}

	if len(dr.AwsRegion) == 0 {
		dr.AwsRegion = "us-east-1"
	}

	if len(dr.ThingNamePrefix) == 0 {
		dr.ThingNamePrefix = "golang_thing"
	}

	if len(dr.ThingTypeName) == 0 {
		dr.ThingTypeName = "simulated-thing"
	}

	if dr.MaxRequestsPerSecond == 0 {
		dr.MaxRequestsPerSecond = 15
	}

	return dr
}

func (dr *DeviceRegistry) Initialize() error {
	fmt.Println("Initializing Thing Registry")

	if err := dr.CreateThingType(dr.ThingTypeName); err != nil {
		fmt.Println(err.Error())
	}

	ConcurrentWorkerExecutor(dr.TotalNumberOfThings, dr.MaxRequestsPerSecond, func(thingId int) error {
		thingName := fmt.Sprintf("%s-%d", dr.ThingNamePrefix, thingId)

		if err := dr.CreateThing(thingName, dr.ThingTypeName); err != nil {
			return fmt.Errorf("Unable to create thing %s: %v\n", thingName, err)
		}

		return nil
	})

	return nil
}

func (dr *DeviceRegistry) Cleanup() error {
	fmt.Println("Cleaning Up Thing Registry")

	// if err := dr.DeprecateThingType(dr.ThingTypeName); err != nil {
	// 	fmt.Println(err.Error())
	// }

	ConcurrentWorkerExecutor(dr.TotalNumberOfThings, dr.MaxRequestsPerSecond, func(thingId int) error {
		thingName := fmt.Sprintf("%s-%d", dr.ThingNamePrefix, thingId)

		if err := dr.DeleteThing(thingName); err != nil {
			return fmt.Errorf("Unable to delete thing %s: %v\n", thingName, err)
		}

		return nil
	})

	// if err := dr.DeleteThingType(dr.ThingTypeName); err != nil {
	// 	fmt.Println(err.Error())
	// }

	return nil
}

func (dr *DeviceRegistry) CreateThing(thingName string, thingTypeName string) error {
	fmt.Printf("Creating Thing %s - Thing Type: %s\n", thingName, thingTypeName)
	params := &iot.CreateThingInput{
		ThingName:     aws.String(thingName),
		ThingTypeName: aws.String(thingTypeName),
	}

	req, resp := dr.client.CreateThingRequest(params)
	if err := req.Send(); err != nil {
		return fmt.Errorf("CreateThingRequest error: %v\nResponse: %v\n", err, resp)
	}

	fmt.Printf("CreateThing result: %v\n", resp)

	return nil
}

func (dr *DeviceRegistry) CreateThingType(thingTypeName string) error {
	fmt.Printf("Creating Thing Type %s\n", thingTypeName)

	params := &iot.CreateThingTypeInput{
		ThingTypeName: aws.String(thingTypeName),
	}

	req, resp := dr.client.CreateThingTypeRequest(params)
	if err := req.Send(); err != nil { // resp is now filled
		return fmt.Errorf("Failed to create thing type: %v\nCreateThingType response: %v\n", err, resp)
	}

	fmt.Printf("CreateThingType result: %v\n", resp)

	return nil
}

func (dr *DeviceRegistry) DeleteThing(thingName string) error {
	fmt.Printf("Deleting Thing %s\n", thingName)

	params := &iot.DeleteThingInput{
		ThingName: aws.String(thingName),
	}

	req, resp := dr.client.DeleteThingRequest(params)
	if err := req.Send(); err != nil {
		return fmt.Errorf("DeleteThingRequest error: %v\nDeleteThing response: %v\n", err, resp)
	}

	fmt.Printf("DeleteThing result: %v\n", resp)

	return nil
}

func (dr *DeviceRegistry) DeleteThingType(thingTypeName string) error {
	fmt.Printf("Deleting Thing Type %s\n", thingTypeName)

	deleteThingTypeInput := &iot.DeleteThingTypeInput{
		ThingTypeName: aws.String(thingTypeName),
	}

	req, resp := dr.client.DeleteThingTypeRequest(deleteThingTypeInput)
	if err := req.Send(); err != nil { // resp is now filled
		return fmt.Errorf("Failed to delete thing type: %v\nDeleteThingType response: %v\n", err, resp)
	}

	fmt.Printf("DeleteThingType result: %v\n", resp)

	return nil
}

func (dr *DeviceRegistry) DeprecateThingType(thingTypeName string) error {
	fmt.Printf("Deprecating Thing Type %s\n", thingTypeName)
	params := &iot.DeprecateThingTypeInput{
		ThingTypeName: aws.String(thingTypeName),
	}

	req, resp := dr.client.DeprecateThingTypeRequest(params)
	if err := req.Send(); err != nil { // resp is now filled
		return fmt.Errorf("Failed to deprecate thing type: %v\nDepcrecateThingType response: %v\n", err, resp)
	}

	fmt.Printf("DeprecateThingType result: %v\n", resp)

	return nil
}
