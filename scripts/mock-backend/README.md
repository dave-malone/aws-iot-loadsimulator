# simulation mock-backend
This directory holds some basic CloudFormation and Lambda functions to create a 
basic backend to run the iot-loadsimulator against. The purpose of this set of 
infrastructure is to demonstrate a set of backend processing/storage for iot data. 

The backend consists of three parts, 1) and SQS queue where the IoT rules engine 
will push MQTT messages, 2) a lambda function that will process and enrich SQS 
messages, 3) an RDS database where the lambda function will persist enriched data. 

This is modeled around a simplified customer use case. It is designed to show 
message processing latencies/timings across the backend infrastructure. 

## creating the infrastructure
It is expected that the reader had an active AWS account and has installed and 
configured the AWS cli. 

Prior to running the create script please review the variables in the create script. 
```$bash
PROJECT="aws-iot-loadsimulator-mockbackend"
VPC_CIDR="10.0.0.0/16"
PUBLICSUBNET1_CIDR="10.0.1.0/24"
PUBLICSUBNET2_CIDR="10.0.2.0/24"
PRIVATESUBNET1_CIDR="10.0.3.0/24"
DATABASEACCESS_CIDR="131.125.33.66/32" # IP/CIDR of your db client
DB_USERNAME="iotdba"
DB_PASSWORD="S0M3THING_SUP3R_S3CR3T"
DB_NAME="iotdb"
TABLE_NAME="iotdata"
```

After reviewing the variables, to create the infrastructure, run the create script:
```$bash
./create-infrastructure.sh
```

