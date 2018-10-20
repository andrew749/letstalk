import ReactDOM from 'react-dom';
import React from 'react';
import { Button, FormGroup, FormControl, ControlLabel, Alert } from "react-bootstrap";
import CookieAwareComponent from './cookie_aware_component.jsx'
import {withCookies} from 'react-cookie';
import getProperty from './context.jsx';

export default class SampleNotification extends React.Component {
    render() {
      const title = getProperty('title');
      const body = getProperty('body');
      const thumbnail = getProperty('thumbnail');
      return (
        <div>
          <h1>{title}</h1>
          <p>{body}</p>
          <img src={thumbnail}/>
        </div>
      );
    }
}

ReactDOM.render(<SampleNotification />, document.getElementById('content'));
