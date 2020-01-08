import React from 'react';
import logo from './logo.svg';
import './App.css';

export default class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      totalThings: 10,
      messagesPerThing: 10,
      secondsBetweenMqttMessages: 10,
      clientsPerWorker: 1000,
      secondsBetweenSnsMessages: 5
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(event) {
    this.setState({[event.target.name]: event.target.value});
    console.log(JSON.stringify(this.state))
  }

  handleSubmit(event) {
    alert('A name was submitted: ' + this.state.value);
    event.preventDefault();
  }

  render() {
    let { totalThings, clientsPerWorker, messagesPerThing, secondsBetweenSnsMessages, secondsBetweenMqttMessages} = this.state

    let totalWorkers = totalThings / clientsPerWorker
    totalWorkers = totalWorkers < 1 ? 1 : totalWorkers

    let totalRampTime = secondsBetweenSnsMessages * totalWorkers

    let durationPerWorker = messagesPerThing * secondsBetweenMqttMessages
    durationPerWorker = durationPerWorker <= 900 ? durationPerWorker : 900

    let mqttMessageTotal = (durationPerWorker / secondsBetweenMqttMessages) * totalThings
    let estimatedDuration = durationPerWorker + (secondsBetweenSnsMessages * (totalWorkers - 1))

    return (
      <div>
        <form onSubmit={this.handleSubmit}>
          <label>
            <input type="number" name="totalThings" value={totalThings} onChange={this.handleChange} />
            Total number of Simulated Things
          </label>
          <br />
          <label>
            <input type="number" name="messagesPerThing" value={messagesPerThing} onChange={this.handleChange} />
            Quantity of MQTT Messages to publish per Thing
          </label>
          <br />
          <label>
            <input type="number" name="secondsBetweenMqttMessages" value={secondsBetweenMqttMessages} onChange={this.handleChange} />
            Seconds between MQTT messages to be published per Thing
          </label>
          <br />
          <input type="submit" value="Begin Simulation" />
        </form>

        <h2>Simulation Engine Resource Estimates:</h2>
        <ul>
          <li>Total number of SNS messages: {totalWorkers.toLocaleString()}</li>
          <li>Seconds between SNS messages (ramp time): {secondsBetweenSnsMessages}</li>
          <li>Total ramp time: {totalRampTime} seconds</li>
          <li>Total number of iot-simulation-worker Lambda functions: {totalWorkers.toLocaleString()}</li>
          <li>Total number of concurrent MQTT clients: {Number(totalThings).toLocaleString()}</li>
          <li>Total number of MQTT messages: {mqttMessageTotal.toLocaleString()}</li>
          <li>iot-simulation-worker Lambda function duration: {durationPerWorker.toLocaleString()} seconds</li>
          <li>Estimated, end-to-end Duration of simulation: {estimatedDuration.toLocaleString()} seconds</li>
        </ul>
      </div>
    );
  }
}
