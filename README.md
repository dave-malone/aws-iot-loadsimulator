# aws-iot-loadsimulator

A Golang based client simulator for AWS IoT Core intended to be run on Lambda.

Basic architecture:

engine -> SNS -> worker

The source for each of these resides under [cmd/lambda](cmd/lambda).

## Setup

Golang 1.12
AWS Account
AWS CLI


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

## Scratch

```bash
THING_TYPE_NAME="simulated-thing"

aws iot create-thing-type --thing-type-name $THING_TYPE_NAME

# setup

END=50
for ((i=0;i<=END;i++)); do
    THING_NAME="golang_thing-$i"
    echo "creating thing $THING_NAME"
    aws iot create-thing --thing-name $THING_NAME --thing-type-name $THING_TYPE_NAME
done

# cleanup

END=50
for ((i=0;i<=END;i++)); do
    THING_NAME="golang_thing-$i"
    echo "deleting thing $THING_NAME"
    aws iot delete-thing --thing-name $THING_NAME
done


aws iot delete-thing-type --thing-type-name $THING_TYPE_NAME
```
