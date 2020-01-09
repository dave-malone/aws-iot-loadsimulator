#!/usr/bin/env bash
set -x
RUN_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
STACK_NAME=aws-iot-loadsimulator
ACCOUNT_ID=$(aws --output text sts get-caller-identity --query 'Account')
S3_BUCKET="${STACK_NAME}-${ACCOUNT_ID}"
MQTT_HOST=$(aws --output text iot describe-endpoint --endpoint-type iot:Data-ATS --query 'endpointAddress')

if ! aws s3api head-bucket --bucket "${S3_BUCKET}"; then
  echo "${S3_BUCKET} does not exist. Creating now..."
  aws s3api create-bucket --bucket "${S3_BUCKET}" --acl private
fi

"${ROOT_DIR}"/create-iot-thing.sh golang_thing
"${RUN_DIR}"/build-engine-lambda.sh
"${RUN_DIR}"/build-worker-lambda.sh

sam package \
  --template-file "${RUN_DIR}"/template.yml \
  --s3-bucket "${S3_BUCKET}" \
  --output-template-file "${RUN_DIR}"/packaged.yml

sam deploy \
  --template-file "${RUN_DIR}"/packaged.yml \
  --stack-name ${STACK_NAME} \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameter-overrides MqttHost="${MQTT_HOST}"
