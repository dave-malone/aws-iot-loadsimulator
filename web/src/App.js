import React from 'react';

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

    let messagesPerSecond = mqttMessageTotal / estimatedDuration
    messagesPerSecond = Math.ceil(messagesPerSecond)

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
          <li>Estimated MQTT messages per second: {messagesPerSecond.toLocaleString()}</li>
          <li>iot-simulation-worker Lambda function duration: {durationPerWorker.toLocaleString()} seconds</li>
          <li>Estimated, end-to-end Duration of simulation: {estimatedDuration.toLocaleString()} seconds</li>
        </ul>
        <p>
          <i>** PLEASE NOTE ** since the maximum execution duration of Lambda today is 15 minutes, the maximum number of MQTT
          messages that can possibly be generated per iot-simulation-worker Lambda is based on the following formula:</i>

          <div>
            <code>
              let durationPerWorker = messagesPerThing * secondsBetweenMqttMessages<br />
              durationPerWorker = durationPerWorker &lt;= 900 ? durationPerWorker : 900<br />
              let mqttMessageTotal = (durationPerWorker / secondsBetweenMqttMessages) * totalThings
            </code>
          </div>
        </p>
      </div>
    );
  }
}
