import React from 'react'
import ReactDOM from 'react-dom'
import './index.css'
import App from './App'
import Amplify from 'aws-amplify'
import config from './aws-exports'
import { withAuthenticator } from 'aws-amplify-react'

Amplify.configure(config)

const AppWithAuth = withAuthenticator(App, {includeGreetings: true})
ReactDOM.render(<AppWithAuth />, document.getElementById('root'))
