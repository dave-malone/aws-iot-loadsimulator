import React from 'react'
import TextField from '@material-ui/core/TextField'
import { withStyles } from '@material-ui/core/styles'
import Button from '@material-ui/core/Button'

const styles = theme => ({
  root: {
    '& .MuiTextField-root': {
      margin: theme.spacing(1),
      width: 200,
    },
  },
})

class SimulationRequestForm extends React.Component {
  render() {
    let {
      classes,
      controls,
      onChangeHandler,
      onSubmitHandler
    } = this.props

    // console.log(`SimulationRequestForm rendering: ${JSON.stringify(this.props)}`)

    return (
      <form onSubmit={onSubmitHandler} className={classes.root}>
        <TextField
          id="totalThings"
          name="totalThings"
          type="number"
          value={controls.totalThings}
          onChange={onChangeHandler}
          label="Total Things"
          required
          variant="outlined"
        />
        <TextField
          id="durationPerWorker"
          name="durationPerWorker"
          type="number"
          min="30"
          max="900"
          value={controls.durationPerWorker}
          onChange={onChangeHandler}
          label="Duration"
          helperText="In Seconds"
          required
          variant="outlined"
        />
        <TextField
          id="secondsBetweenMqttMessages"
          name="secondsBetweenMqttMessages"
          type="number"
          value={controls.secondsBetweenMqttMessages}
          onChange={onChangeHandler}
          label="Sleep time"
          helperText="Seconds between messages"
          required
          variant="outlined"
        />
        <Button
          type="submit"
          color="primary"
          variant="contained"
          onClick={onSubmitHandler}
          style={{ marginTop: '1.4em', marginLeft: '1.2em' }}>
          Begin Simulation
        </Button>
      </form>
    )
  }
}

export default withStyles(styles)(SimulationRequestForm)
