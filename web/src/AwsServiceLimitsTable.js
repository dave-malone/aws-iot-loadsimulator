import React from 'react'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import CardHeader from '@material-ui/core/CardHeader'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableContainer from '@material-ui/core/TableContainer'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Paper from '@material-ui/core/Paper'

export default class AwsServiceLimitsTable extends React.Component {

  render() {
    const { serviceLimits } = this.props
    // console.log(`service limits: ${serviceLimits}`)
    return (
      <Card>
        <CardHeader title="Relevant AWS Limits" />
        <CardContent>
          <TableContainer component={Paper}>
            <Table size="small" aria-label="a dense table">
              <TableHead>
                <TableRow>
                  <TableCell>Service</TableCell>
                  <TableCell>Resource</TableCell>
                  <TableCell>Default Limit</TableCell>
                  <TableCell>Adjustable</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
              {serviceLimits.map((row, i) => (
                <TableRow key={i} style={{backgroundColor: row.highlight ? 'yellow' : ''}}>
                  <TableCell>{row.service}</TableCell>
                  <TableCell>{row.resource}</TableCell>
                  <TableCell>{row.limit}</TableCell>
                  <TableCell>{row.adjustable}</TableCell>
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
