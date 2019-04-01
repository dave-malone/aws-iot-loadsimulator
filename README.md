# aws-iot-loadsimulator

A Golang based client simulator for AWS IoT Core intended to be run on Lambda.

Basic architecture:

engine -> SNS -> worker

The source for each of these resides under cmd/lambda. 

Project structured according to https://github.com/golang-standards/project-layout
