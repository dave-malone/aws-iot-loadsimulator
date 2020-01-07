!#/bin/bash

set -x

if [ -z "$1" ]; then
  echo "1st arg is the Thing Name, and may not be empty"
  exit 0
fi

THING_NAME=$1
THING_TYPE_NAME=simulated-thing

if aws iot describe-thing --thing-name $THING_NAME 2>&1 | grep -q 'ResourceNotFoundException'
then
  echo "Thing with name '$THING_NAME' does not exist. Creating AWS resources"

  ACCOUNT_ID=$(aws sts get-caller-identity | jq -r '.Account')
  AWS_REGION=$(aws configure get region)

  echo "Creating Thing in IoT Thing Registry"
  aws iot create-thing \
    --thing-name $THING_NAME \
    --thing-type-name $THING_TYPE_NAME \
    | THING_ARN=$(jq -r '.thingArn')

  echo "Creating Keys and Certificate"
  CERTIFICATE_ARN=$(aws iot create-keys-and-certificate \
    --set-as-active \
    --certificate-pem-outfile "./certs/golang_thing.cert.pem" \
    --public-key-outfile "./certs/golang_thing.public.key" \
    --private-key-outfile "./certs/golang_thing.private.key" \
    | jq -r '.certificateArn')

  sed -e 's/{{account-id}}/'"$ACCOUNT_ID"'/g ; s/{{aws-region}}/'"$AWS_REGION"'/g' \
    ./scripts/thing-policy.template > ./scripts/thing-policy.json

  echo "Creating Thing Policy"
  aws iot create-policy \
    --policy-name "${THING_NAME}-Policy" \
    --policy-document file://scripts/thing-policy.json

  echo "Attaching Thing Policy to Thing Certificate"
  aws iot attach-policy \
    --policy-name "${THING_NAME}-Policy" \
    --target $CERTIFICATE_ARN

  echo "Attaching Thing to Thing Certificate (Principal)"
  aws iot attach-thing-principal \
    --thing-name $THING_NAME \
    --principal $CERTIFICATE_ARN
else
  echo "Thing with name '$THING_NAME' already exists - skipping resource creation"
fi
