#!/usr/bin/env bash
set -x
RUN_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )"/.. && pwd )"

if [ -z "${1}" ]; then
  echo "1st arg is the Thing Name, and may not be empty"
  exit 0
fi

THING_NAME="${1}"
THING_TYPE_NAME=simulated-thing

if ! aws --output text iot list-thing-types --query 'thingTypes[].thingTypeName' | grep -q "${THING_TYPE_NAME}"; then
  echo "Thing Type with name \"${THING_TYPE_NAME}\" does not exist. Creating AWS resources"
  aws iot create-thing-type --thing-type-name "${THING_TYPE_NAME}"
fi

if ! aws iot describe-thing --thing-name "${THING_NAME}" > /dev/null 2>&1; then
  echo "Thing with name \"${THING_NAME}\" does not exist. Creating AWS resources"

  ACCOUNT_ID=$(aws --output text sts get-caller-identity --query 'Account')
  AWS_REGION=$(aws configure get region)

  echo "Creating Thing in IoT Thing Registry"
  aws --output text iot create-thing \
    --thing-name "${THING_NAME}" \
    --thing-type-name ${THING_TYPE_NAME} \
    --query 'thingArn'

  echo "Creating Keys and Certificate"
  CERTIFICATE_ARN=$(aws --output text iot create-keys-and-certificate \
    --set-as-active \
    --certificate-pem-outfile "${PROJECT_DIR}/certs/golang_thing.cert.pem" \
    --public-key-outfile "${PROJECT_DIR}/certs/golang_thing.public.key" \
    --private-key-outfile "${PROJECT_DIR}/certs/golang_thing.private.key" \
    --query 'certificateArn')

  sed -e 's/{{account-id}}/'"$ACCOUNT_ID"'/g ; s/{{aws-region}}/'"$AWS_REGION"'/g' \
    "${RUN_DIR}"/thing-policy.template > "${RUN_DIR}"/thing-policy.json

  echo "Creating Thing Policy"
  aws iot create-policy \
    --policy-name "${THING_NAME}-Policy" \
    --policy-document file://"${RUN_DIR}"/thing-policy.json

  echo "Attaching Thing Policy to Thing Certificate"
  aws iot attach-policy \
    --policy-name "${THING_NAME}-Policy" \
    --target "${CERTIFICATE_ARN}"

  echo "Attaching Thing to Thing Certificate (Principal)"
  aws iot attach-thing-principal \
    --thing-name "${THING_NAME}" \
    --principal "${CERTIFICATE_ARN}"
else
  echo "Thing with name \"${THING_NAME}\" already exists - skipping resource creation"
fi
