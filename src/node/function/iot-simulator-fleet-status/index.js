const AWS = require('aws-sdk')

AWS.config.region = process.env.AWS_REGION
const DDB_TABLE_NAME = process.env.DDB_TABLE_NAME
const dynamodb = new AWS.DynamoDB()


function updateDeviceStatus(event) {
  /*
    Event payload. See: https://docs.aws.amazon.com/iot/latest/developerguide/life-cycle-events.html#connect-disconnect

    {
        "clientId": "186b5",
        "timestamp": 1573002340451,
        "eventType": "disconnected",
        "sessionIdentifier": "a4666d2a7d844ae4ac5d7b38c9cb7967",
        "principalIdentifier": "12345678901234567890123456789012",
        "clientInitiatedDisconnect": true,
        "disconnectReason": "CLIENT_INITIATED_DISCONNECT",
        "versionNumber": 0
    }
  */
  return new Promise((resolve, reject) => {
    var params = {
      TableName: DDB_TABLE_NAME,
      Item: {
          'clientId': {S: event.clientId},
          'deviceStatus': {S: event.eventType},
          'clientInitiatedDisconnect': {S: `${event.clientInitiatedDisconnect}`},
          'disconnectReason': {S: `${event.disconnectReason}` },
          'timestamp': {N: `${event.timestamp}` }
      }
    }

    // See: https://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/DynamoDB.html#putItem-property
    dynamodb.putItem(params, function(err, data) {
      if (err){
        console.log(`Failed to putItem in DynamoDB: ${err}`, err.stack)
        reject(`Failed to putItem in DynamoDB: ${err.errorMessage}`)
      } else{
        console.log(`Successfully putItem in DynamoDB`)
        resolve(`Successfully putItem in DynamoDB`)
      }
    })
  })
}


exports.handler = async (event) => {
    console.log(`Received event: ${JSON.stringify(event)}`)

    await updateDeviceStatus(event)
        .then((result) => {
            return {
                statusCode: 200,
                body: result,
            }
        }).catch((err) => {
            return {
                statusCode: 200,
                body: err,
            }
        })
};
