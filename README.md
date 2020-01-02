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
