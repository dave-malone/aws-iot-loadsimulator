#!/bin/bash

set -x

STACK_NAME=aws-iot-loadsimulator
ACCOUNT_ID=$(aws sts get-caller-identity | jq -r '.Account')
S3_BUCKET="${STACK_NAME}-${ACCOUNT_ID}"

if aws s3 ls "s3://$S3_BUCKET" 2>&1 | grep -q 'NoSuchBucket'
then
  echo "$S3_BUCKET does not exist. Creating now..."
  aws s3 mb s3://$S3_BUCKET
fi

./scripts/build-engine-lambda.sh
./scripts/build-worker-lambda.sh

aws s3 cp ./build/engine-handler.zip s3://$S3_BUCKET/engine-handler.zip
aws s3 cp ./build/worker-handler.zip s3://$S3_BUCKET/worker-handler.zip

aws cloudformation create-stack \
  --stack-name $STACK_NAME \
  --template-body file://$(pwd)/scripts/cloudformation.yaml \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameters ParameterKey=S3BucketName,ParameterValue=$S3_BUCKET

aws cloudformation wait stack-create-complete --stack-name $STACK_NAME