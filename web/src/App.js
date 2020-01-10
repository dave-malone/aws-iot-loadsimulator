import React from 'react'
import Grid from '@material-ui/core/Grid'
import ResourceUsageEstimator from './ResourceUsageEstimator'
import SimulationRequestForm from './SimulationRequestForm'
import AwsServiceLimitsTable from './AwsServiceLimitsTable'

const serviceLimitsRef = [
  {id: 0, highlight: false, service: 'AWS IoT Core', limit: '500,000', adjustable: 'Yes', resource: 'Maximum concurrent client connections per account'},
  {id: 1, highlight: false, service: 'AWS IoT Core', limit: '500', adjustable: 'Yes', resource: 'Connect requests per second per account'},
  {id: 2, highlight: false, service: 'AWS IoT Core', limit: '20,000', adjustable: 'Yes', resource: 'Inbound publish requests per second per account'},
  {id: 3, highlight: false, service: 'AWS IoT Core', limit: '128 KB', adjustable: 'No', resource: 'Message Size'},
  {id: 4, highlight: false, service: 'AWS IoT Rules Engine', limit: '20,000', adjustable: 'Yes', resource: 'Inbound publish requests per second per account'},
  {id: 5, highlight: false, service: 'AWS Lambda', limit: '1,000', adjustable: 'Yes', resource: 'Concurrent executions'},
  {id: 6, highlight: false, service: 'AWS Lambda', limit: '900 seconds (15 minutes)', adjustable: 'No', resource: 'Function timeout'},
  {id: 7, highlight: false, service: 'AWS Lambda', limit: '1,024', adjustable: 'No', resource: 'File descriptors'},
  {id: 8, highlight: false, service: 'AWS Lambda', limit: '1,024', adjustable: 'No', resource: 'Execution processes/threads'},
  {id: 9, highlight: false, service: 'Amazon SNS', limit: '30,000 / second', adjustable: 'Yes', resource: 'Publish (US East Region)'}
]

const controlsRef = {
  totalThings: 1000,
  durationPerWorker: 900,
  secondsBetweenMqttMessages: 5,
  clientsPerWorker: 500,
  secondsBetweenSnsMessages: 5,
}

const resourceUsageEstimateRef = [
  { key: "totalSnsMessages", associatedServiceLimit: serviceLimitsRef[9], label: "SNS Messages", value: 0},
  { key: "totalWorkers", associatedServiceLimit: serviceLimitsRef[5], label: "Worker functions", value: 0},
  { key: "secondsBetweenSnsMessages", associatedServiceLimit: null, label: "Seconds between SNS messages (ramp time)", value: 0},
  { key: "totalRampTime", associatedServiceLimit: null, label: "Total Ramp time (seconds)", value: 0},
  { key: "totalThings", associatedServiceLimit: serviceLimitsRef[0], label: "Simulated Things", value: 0},
  { key: "messagesPerThing", associatedServiceLimit: null, label: "Messages Per Thing", value: 0},
  { key: "mqttMessageTotal", associatedServiceLimit: null, label: "Total MQTT Messages", value: 0},
  { key: "messagesPerSecond", associatedServiceLimit: serviceLimitsRef[2], label: "MQTT Messages / Second", value: 0},
  { key: "durationPerWorker", associatedServiceLimit: null, label: "Worker function duration (seconds)", value: 0},
  { key: "estimatedDuration", associatedServiceLimit: null, label: "End-to-end simulation duration (seconds)", value: 0},
]

export default class App extends React.Component {
  constructor(props) {
    super(props)

    this.state = {
      controls: Object.assign({}, controlsRef),
      resourceUsageEstimate: [...resourceUsageEstimateRef],
      serviceLimits: [...serviceLimitsRef],
    }

    this.handleChange = this.handleChange.bind(this)
    this.handleSubmit = this.handleSubmit.bind(this)
    this.computeEstimates = this.computeEstimates.bind(this)
  }

  componentDidMount(){
    this.handleChange({target:{name:'', value:''}})
  }

  handleChange(event) {
    // console.log(`form input change event for target ${event.target.name} = ${event.target.value}`)

    let updatedControls = Object.assign({},{...this.state.controls}, { [event.target.name]: Number(event.target.value) })
    let updatedResourceUsageEstimate = this.computeEstimates(updatedControls)

    this.setState({resourceUsageEstimate: updatedResourceUsageEstimate, controls: updatedControls})
    // console.log(JSON.stringify(this.state))
  }

  handleSubmit(event) {
    console.log('Submitting simulation request')
    event.preventDefault()
  }

  computeEstimates(controls){
    console.log(`computing estimates using controls ${JSON.stringify(controls)}`)

    let {
      totalThings,
      durationPerWorker,
      clientsPerWorker,
      secondsBetweenMqttMessages,
      secondsBetweenSnsMessages
    } = controls

    let resourceUsageEstimate = [...resourceUsageEstimateRef]

    const findByKey = (key) => {
      let element = resourceUsageEstimate.find( element => element.key === key)
      return element
    }

    const setValue = (key, value, highlightAssociatedServiceLimit) => {
      let element = findByKey(key)
      element.value = value

      if(element.associatedServiceLimit){
        element.associatedServiceLimit.highlight = highlightAssociatedServiceLimit
      }
    }

    let totalWorkers = Math.ceil(totalThings / clientsPerWorker)
    let totalRampTime = secondsBetweenSnsMessages * totalWorkers
    let messagesPerThing = durationPerWorker / secondsBetweenMqttMessages
    let mqttMessageTotal = (durationPerWorker / secondsBetweenMqttMessages) * totalThings
    let estimatedDuration = durationPerWorker + (secondsBetweenSnsMessages * (totalWorkers - 1))
    let messagesPerSecond = Math.ceil(mqttMessageTotal / estimatedDuration)

    setValue("totalWorkers", totalWorkers, (totalWorkers > 1000))
    setValue("totalSnsMessages", totalWorkers)
    setValue("totalRampTime", totalRampTime)
    setValue("messagesPerThing", messagesPerThing)
    setValue("mqttMessageTotal", mqttMessageTotal)
    setValue("estimatedDuration", estimatedDuration)
    setValue("messagesPerSecond", messagesPerSecond, (messagesPerSecond > 20000))
    setValue("totalThings", totalThings, (totalThings > 500000))
    setValue("secondsBetweenSnsMessages", secondsBetweenSnsMessages)
    setValue("durationPerWorker", durationPerWorker)

    return resourceUsageEstimate
  }

  render() {
    return (
      <Grid container spacing={2}>
        <Grid item xs={12}>
          <Grid container justify="center">
            <SimulationRequestForm
              controls={this.state.controls}
              onChangeHandler={this.handleChange}
              onSubmitHandler={this.handleSubmit} />
          </Grid>
        </Grid>
        <Grid item xs={12}>
          <Grid container justify="center" direction="row" wrap="nowrap">
            <Grid item xs={5}>
              <ResourceUsageEstimator resourceUsageEstimate={this.state.resourceUsageEstimate} />
            </Grid>
            <Grid item>
              <AwsServiceLimitsTable serviceLimits={this.state.serviceLimits} />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    )
  }
}
