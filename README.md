# aws-iot-loadsimulator

A Golang based client simulator for AWS IoT Core intended to be run on Lambda.

Basic architecture:

engine -> SNS -> worker

The source for each of these resides under [cmd/lambda](cmd/lambda).

## Setup

Golang 1.12
AWS Account
AWS CLI
[SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)


This project uses [Go Modules](https://blog.golang.org/using-go-modules)

Project structured according to https://github.com/golang-standards/project-layout


## Local Build & Test with CLI Interfaces

Device Registry:

```bash
go run cmd/cli/registry/main.go -mode init -total-things 1000
```

Simulation Engine:

```bash
go run cmd/cli/engine/main.go \
  -sns-topic-arn arn:aws:sns:us-east-1:068311527115:iot_simulator_notifications
```

Simulation Worker:

```bash
go run cmd/cli/worker/main.go -max-clients 100 -seconds-between-messages 10 -total-messages-per-client 5
```

## Deploy Simulator Lambda Functions

```bash
./scripts/deploy-sam.sh
```

## TODO

* Build simple UI to kick off device simulation, view device simulation stats
* Add CW Metric Dashboards into CW templates
* Can't use fleet indexing metrics without things being registered in the device registry
* Current issue with populating device registry is the max rate at which we can call CreateThing (default 15 TPS); this makes populating the device registry for large scale simulations potentially lengthy (100,000 items created using CreateThing at that rate would take about 1.85 hours to complete).
* Could rely on [AWS IoT Lifecycle Events](https://docs.aws.amazon.com/iot/latest/developerguide/life-cycle-events.html) to get fleet online / offline state, but would need to offload to something that can handle atomic update requests - i.e. Redis
* As part of the demo, can we illustrate the use of Fine-Grained Logging on a select group of things? https://docs.aws.amazon.com/iot/latest/developerguide/cloud-watch-logs.html#configure-logging
* Demonstrate blue/green deployment change to a Rule 
