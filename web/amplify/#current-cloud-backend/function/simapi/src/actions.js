const AWS = require('aws-sdk')

AWS.config.region = process.env.AWS_REGION

const iot = new AWS.Iot({
  region: process.env.AWS_REGION
})

function fetchConnectedDeviceCount(){
  return new Promise((resolve, reject) => {
    var params = {
      queryString: 'connectivity.connected:true',
    }

    iot.getStatistics(params, (err, data) => {
      if (err){
        console.log(`Failed to getStatistics: ${err.errorMessage}`, err.stack)
        reject(`Failed to getStatistics: ${err.errorMessage}`)
      } else {
        console.log(`getStatistics result: ${JSON.stringify(data)}`)
        resolve(data)
      }
    })
  })
}

module.exports = {
  fetchConnectedDeviceCount,
}
