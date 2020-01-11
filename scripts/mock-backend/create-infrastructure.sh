#!/usr/bin/env bash

RUN_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ACCOUNT_ID=$(aws --output text sts get-caller-identity --query 'Account')
VPC_STACK_NAME=aws-iot-loadsimulator-mockbackend-vpc
STACK_NAME=aws-iot-loadsimulator-mockbackend
S3_BUCKET="${STACK_NAME}-${ACCOUNT_ID}"

if ! aws s3api head-bucket --bucket "${S3_BUCKET}"; then
  echo "${S3_BUCKET} does not exist. Creating now..."
  aws s3api create-bucket --bucket "${S3_BUCKET}" --acl private
fi

PROJECT="${STACK_NAME}"
VPC_CIDR="10.0.0.0/16"
PUBLICSUBNET1_CIDR="10.0.1.0/24"
PUBLICSUBNET2_CIDR="10.0.2.0/24"
PRIVATESUBNET1_CIDR="10.0.3.0/24"
DATABASEACCESS_CIDR="72.21.196.0/24"
DB_USERNAME="iotdba"
DB_PASSWORD="wcqvo2dAaE7kdy3yYyPC"
DB_NAME="iotdb"
TABLE_NAME="iotdata"

pushd "${RUN_DIR}"/lambda || exit 127
./zip-lambda.sh
popd || exit 128

aws cloudformation create-stack \
  --stack-name ${VPC_STACK_NAME} \
  --template-body file://"${RUN_DIR}"/vpc.json \
  --parameters ParameterKey=Project,ParameterValue=${PROJECT} \
ParameterKey=VpcCIDR,ParameterValue=${VPC_CIDR} \
ParameterKey=PublicSubnet1CIDR,ParameterValue=${PUBLICSUBNET1_CIDR} \
ParameterKey=PublicSubnet2CIDR,ParameterValue=${PUBLICSUBNET2_CIDR} \
ParameterKey=PrivateSubnet1CIDR,ParameterValue=${PRIVATESUBNET1_CIDR} \
ParameterKey=DatabaseAccessCIDR,ParameterValue=${DATABASEACCESS_CIDR}

until aws --output text cloudformation describe-stacks --query "Stacks[?StackName=='${VPC_STACK_NAME}'].StackStatus" | grep -q COMPLETE; do
    printf '.'
    sleep 5
done
echo

aws cloudformation create-stack \
  --stack-name ${STACK_NAME} \
  --template-body file://"${RUN_DIR}"/infrastructure.json \
  --capabilities CAPABILITY_NAMED_IAM \
  --parameters ParameterKey=Project,ParameterValue=${PROJECT} \
ParameterKey=CodeBucketName,ParameterValue="${S3_BUCKET}" \
ParameterKey=Username,ParameterValue=${DB_USERNAME} \
ParameterKey=Password,ParameterValue=${DB_PASSWORD} \
ParameterKey=DatabaseName,ParameterValue=${DB_NAME} \
ParameterKey=SQSTableName,ParameterValue=${TABLE_NAME}







