import React from 'react'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableContainer from '@material-ui/core/TableContainer'
import TableRow from '@material-ui/core/TableRow'
import Paper from '@material-ui/core/Paper'

export default class ResourceUsageEstimator extends React.Component {
  render() {
    let { resourceUsageEstimate } = this.props

    // console.log(`resourceUsageEstimate: ${JSON.stringify(resourceUsageEstimate)}`)
    return (
      <Card>
        <CardHeader title="AWS Resource Usage Estimates" />
        <CardContent>
          <TableContainer component={Paper}>
            <Table size="small">
              <TableBody>
                {resourceUsageEstimate.map((row, i) => (
                  <TableRow key={i} style={{backgroundColor: row.associatedServiceLimit ? (row.associatedServiceLimit.highlight ? 'yellow' : '') : ''}}>
                    <TableCell>{row.label}</TableCell>
                    <TableCell>{row.value.toLocaleString()}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>
    )
  }
}
