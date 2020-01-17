/*
Copyright 2017 - 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License"). You may not use this file except in compliance with the License. A copy of the License is located at
    http://aws.amazon.com/apache2.0/
or in the "license" file accompanying this file. This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and limitations under the License.
*/

var express = require('express')
var bodyParser = require('body-parser')
var awsServerlessExpressMiddleware = require('aws-serverless-express/middleware')

const AWS = require('aws-sdk')
AWS.config.region = process.env.AWS_REGION
const lambda = new AWS.Lambda()

function delegateToSimulationEngineFunction(event) {
  console.log(`delegating to simulation engine function\n`)
  return new Promise((resolve, reject) => {
    let params = {
      FunctionName: 'iot-simulator-engine',
      Payload: event.body
    }

    console.log(`sending payload ${JSON.stringify(params)} to iot-simulator-engine\n`)

    lambda.invoke(params, (err, data) => {
      if (err){
        console.log(`Failed to invoke lambda: ${err}\n`, err.stack)
        reject(`Failed to invoke lambda: ${err.errorMessage}`)
      } else{
        console.log(`Successfully invoked lambda: ${data.Payload}\n`)
        resolve(`Successfully invoked lambda: ${data.Payload}`)
      }
    })
  })
}

// declare a new express app
var app = express()
app.use(bodyParser.json())
app.use(awsServerlessExpressMiddleware.eventContext())

// Enable CORS for all methods
app.use(function(req, res, next) {
  res.header("Access-Control-Allow-Origin", "*")
  res.header("Access-Control-Allow-Headers", "*")
  next()
});

app.post('/engine', async function(req, res) {
  let event = req.apiGateway.event
  console.log(`Received event: ${JSON.stringify(event)}`)

  await delegateToSimulationEngineFunction(event)
      .then((result) => {
          res.json({
            success: `post call succeed! - ${result}`,
            body: result
          })
      }).catch((err) => {
        res.json({
          success: `post call failed: ${err}`,
          body: `post call failed: ${err}`,
        })
      })
});


app.listen(3000, function() {
    console.log("App started")
});

// Export the app object. When executing the application local this does nothing. However,
// to port it to AWS Lambda we will create a wrapper around that will load the app from
// this file
module.exports = app
