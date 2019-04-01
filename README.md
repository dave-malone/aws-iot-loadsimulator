# aws-iot-loadsimulator

A Golang based client simulator for AWS IoT Core intended to be run on Lambda.

Basic architecture:

engine -> SNS -> worker

The source for each of these resides under [cmd/lambda](cmd/lambda).

## Setup

This project uses [dep](https://golang.github.io/dep/docs/introduction.html)

Project structured according to https://github.com/golang-standards/project-layout

## TODO

* Script to generate SNS topic
* Script to deploy each Lambda
* Maybe even one script to deploy it all!
