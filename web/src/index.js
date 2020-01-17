import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import * as serviceWorker from './serviceWorker';
import Amplify from 'aws-amplify';
// import config from './amplify-config'
import config from './aws-exports'
import { withAuthenticator } from 'aws-amplify-react';
Amplify.configure(config);
const AppWithAuth = withAuthenticator(App, { includeGreetings: true });
ReactDOM.render(<AppWithAuth />, document.getElementById('root'));

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
