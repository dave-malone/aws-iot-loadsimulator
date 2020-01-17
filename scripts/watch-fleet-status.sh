#!/bin/bash

DEVICE_STATUS=$1

if [ -z "$DEVICE_STATUS" ]; then
  DEVICE_STATUS=connected
fi

if [ "$DEVICE_STATUS" == "connected" ]; then
  watch 'aws dynamodb query \
    --table-name iot-simulator-fleet-status \
    --index-name deviceStatus-index \
    --select COUNT \
    --key-condition-expression "deviceStatus = :status" \
    --expression-attribute-values "{\":status\": {\"S\":\"connected\"}}"'
else
  watch 'aws dynamodb query \
    --table-name iot-simulator-fleet-status \
    --index-name deviceStatus-index \
    --select COUNT \
    --key-condition-expression "deviceStatus = :status" \
    --expression-attribute-values "{\":status\": {\"S\":\"disconnected\"}}"'
fi
