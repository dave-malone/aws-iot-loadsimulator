import React, { Component } from 'react'
import './App.css'
import { device } from 'aws-iot-device-sdk'

const client = device({
  protocol: 'wss',
  clientId: 'aws-iot-loadsimulator-ui',
  host: 'a1tq0bx5we8tnk-ats.iot.us-east-1.amazonaws.com',
  accessKeyId: 'AKIAQ7Z5PF3FS4CD3YW7',
  secretKey: '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'
})

class App extends Component {

  constructor(props) {
    super(props)
    this.state = {
      connected_device_count: 0,
      last_updated: new Date(),
      error: '',
      connected_devices: [],
      disconnected_devices: []
    }
  }

  componentWillMount(){
    client.subscribe('golang_simulator_connected_events')
    client.on('message', (topic, payload) => {
      console.log('message', topic, payload.toString())

      let event = JSON.parse(payload.toString())

      if(event.eventType === 'connected') {
        this.setState((state, props) => {
          const newState = Object.assign(state, { last_updated: new Date(), error: '' })
          newState.connected_device_count = newState.connected_devices.push(event)
          //console.log(JSON.stringify(newState))
          return newState
        })
      } else {
        this.setState((state, props) => {
          const newState = Object.assign(state, { last_updated: new Date(), error: '' })
          newState.connected_device_count = newState.connected_device_count - 1
          //newState.connected_devices = state.connected_devices.filter(existingEvent => existingEvent.clientId !== event.clientId)

          //console.log(JSON.stringify(newState))
          return newState
        })
      }
    })
  }

  render() {
    return (
      <div>
        { this.state.error &&
          <div>{this.state.error}</div>
        }

        <h3>Connected Clients ({this.state.connected_device_count})</h3>
        <i>Last Updated: {this.state.last_updated.toString()}</i>
        <table>
          <thead>
            <tr>
              <th>Client ID</th>
              <th>Timestamp</th>
              <th>Event Type</th>
            </tr>
          </thead>
          <tbody>
          { this.state.connected_devices.map((event, i) => {
            return (
              <tr key={i}>
                <td>{event.clientId}</td>
                <td>{event.timestamp}</td>
                <td>{event.eventType}</td>
              </tr>
            )
          })}
          </tbody>
        </table>
      </div>
    )
  }
}

export default App
