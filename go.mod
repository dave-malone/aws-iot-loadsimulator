module github.com/dave-malone/aws-iot-loadsimulator/

go 1.12

replace github.com/dave-malone/aws-iot-loadsimulator/pkg => ./pkg

replace github.com/dave-malone/aws-iot-loadsimulator/pkg/mqtt => ./pkg/mqtt

require (
	github.com/aws/aws-lambda-go v1.13.3
	github.com/dave-malone/aws-iot-loadsimulator/pkg v0.0.0-00010101000000-000000000000
)
