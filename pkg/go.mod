module github.com/dave-malone/aws-iot-loadsimulator/pkg

go 1.12

replace github.com/dave-malone/aws-iot-loadsimulator/pkg/mqtt => ./mqtt

require (
	github.com/aws/aws-sdk-go v1.27.0
	github.com/dave-malone/aws-iot-loadsimulator/pkg/mqtt v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.4.0 // indirect
)
