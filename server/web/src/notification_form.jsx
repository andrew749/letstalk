import {withCookies} from 'react-cookie';
import React from 'react';
import locale from 'react-json-editor-ajrm/locale/en';
import JSONInput from 'react-json-editor-ajrm';
import { Button, FormGroup, FormControl, ControlLabel, Alert } from "react-bootstrap";
import {serverUrl} from './config.js'

class NotificationForm extends React.Component {
  constructor(props) {
    super(props);

    const {cookies} = props;

    this.state = {
      notification:{
        title: "",
        message: "",
        template: "sample_template.html",
        group: "",
        thumbnail: "",
        data: {title:"Hello World", body:"Body", thumbnail:"https://www.dike.lib.ia.us/images/sample-1.jpg/image"},
        deeplink: "",
      },
      sessionId: cookies.get('sessionId')
    }

    this.handleChange = this.handleChange.bind(this);
    this.handleChangeJsonBox = this.handleChangeJsonBox.bind(this);
    this.onSubmit = this.onSubmit.bind(this);
  }

  validateForm() {
    const notification = this.state.notification;
    return notification.title.length > 0 &&
    notification.message.length > 0 &&
    notification.template.length > 0 &&
    notification.group.length > 0 &&
    notification.deeplink.length > 0;
  }

  handleChange(event) {
    const value = event.target.value;
    const id = event.target.id;
    this.setState(function (state, props){
      state.notification[id] = value;
      this.props.onUpdateData(state.notification);
      return state;
    });
  }

  handleChangeJsonBox(data) {
    this.setState(function (state, props){
      state.notification.data = data.jsObject;
      this.props.onUpdateData(state.notification);
      return state;
    });
  }

  onSubmit(event) {
    const {cookies} = this.props;
    event.preventDefault();
    // send to api server
    const notifiction = this.state.notification;
    fetch(`${serverUrl}/admin/campaign`, {
      method: 'POST',
      body: JSON.stringify({
        title: notification.title,
        message: notification.message,
        template: notification.template,
        group: notification.group,
        thumbnail: notification.thumbnail,
        data: notification.data,
        deeplink: notification.deeplink
      })
    })
    .then(response => response.json())
    .then((data) => {
      // handle success response
    }).catch(err => {
    });

  }

  render() {
    let alert;
    if (this.state.submitState) {
      if (this.state.submitState === "SUCCESS") {
        alert = (<Alert bsStyle="success">Succesfully sent campaign.</Alert>);
      } else if (this.state.submitState === "ERROR") {
        alert = (<Alert bsStyle="danger">Failed to send campaign because {this.state.err}</Alert>)
      }
    }
    return (
      <div className="notification-form">
        <h2>Notification Management Console</h2>
        <form onSubmit={this.onSubmit}>
          <FormGroup controlId="title" bsSize="large">
            <ControlLabel>Title</ControlLabel>
            <FormControl
              autoFocus
              value={this.state.title}
              onChange={this.handleChange}
            />
          </FormGroup>
          <FormGroup controlId="message" bsSize="large">
            <ControlLabel>Message</ControlLabel>
            <FormControl
              value={this.state.message}
              onChange={this.handleChange}
            />
          </FormGroup>
          <FormGroup controlId="template" bsSize="large">
            <ControlLabel>Template</ControlLabel>
            <FormControl
              value={this.state.template}
              onChange={this.handleChange}
            />
          </FormGroup>
          <FormGroup controlId="group" bsSize="large">
            <ControlLabel>Group</ControlLabel>
            <FormControl
              value={this.state.group}
              onChange={this.handleChange}
            />
          </FormGroup>
          <FormGroup controlId="thumbnail" bsSize="large">
            <ControlLabel>Thumbnail</ControlLabel>
            <FormControl
              value={this.state.thumbnail}
              onChange={this.handleChange}
            />
          </FormGroup>
          <FormGroup controlId="data" bsSize="large">
            <ControlLabel>Data</ControlLabel>
            <JSONInput
              id='notification-content-editor'
              placeholder={{title:"Hello World", body:"Body", thumbnail:"https://www.dike.lib.ia.us/images/sample-1.jpg/image"}}
              locale={ locale }
              onChange={this.handleChangeJsonBox}
              value={this.state.data}
              height='200px'
            />
          </FormGroup>
          <FormGroup controlId="deeplink" bsSize="large">
            <ControlLabel>Deeplink</ControlLabel>
            <FormControl
              value={this.state.deeplink}
              onChange={this.handleChange}
            />
          </FormGroup>
          <Button
            block
            bsSize="large"
            disabled={!this.validateForm()}
            type="submit">
            Send
          </Button>
        </form>
        <div className="message-container">
          {alert}
        </div>
      </div>
    );
  }
}

const CookieNotificationForm = withCookies(NotificationForm);

export default CookieNotificationForm;
