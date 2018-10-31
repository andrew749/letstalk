import {withCookies} from 'react-cookie';
import React from 'react';
import locale from 'react-json-editor-ajrm/locale/en';
import JSONInput from 'react-json-editor-ajrm';
import { Button, FormGroup, FormControl, ControlLabel, Alert } from "react-bootstrap";
import {serverUrl} from './config.js'

const FORM_ELEMENT_SIZE = 'small';

class NotificationForm extends React.Component {
  constructor(props) {
    super(props);

    const {cookies} = props;

    this.state = {
      notification:{
        title: "Title",
        message: "Message",
        templatePath: "sample_template.html",
        groupId: "1",
        runId: "Notification Campaign Run #1",
        templateMetadata: {
          title:"Hello World",
          body:"Body",
          thumbnail:"https://www.dike.lib.ia.us/images/sample-1.jpg/image"
        }
      },
      sessionId: cookies.get('sessionId')
    }

    props.onUpdateData(this.state.notification);

    this.handleChange = this.handleChange.bind(this);
    this.handleChangeJsonBox = this.handleChangeJsonBox.bind(this);
    this.onSubmit = this.onSubmit.bind(this);
  }

  validateForm() {
    const notification = this.state.notification;
    return (
      notification.title &&
      notification.message &&
      notification.templatePath &&
      notification.groupId &&
      notification.runId
    ) &&
    notification.title.length > 0 &&
    notification.message.length > 0 &&
    notification.templatePath.length > 0 &&
    notification.groupId.length > 0 &&
    notification.runId.length > 0;
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
      state.notification.templateMetadata = data.jsObject;
      this.props.onUpdateData(state.notification);
      return state;
    });
  }

  onSubmit(event) {
    const {cookies} = this.props;
    event.preventDefault();
    // send to api server
    const notification = this.state.notification;
    fetch(`${serverUrl}/admin/campaign`, {
      method: 'POST',
      body: JSON.stringify({
        title: notification.title,
        message: notification.message,
        groupId: notification.groupId,
        runId: notification.runId,
        templatePath: notification.templatePath,
        templateMetadata: notification.templateMetadata
      })
    })
    .then(response => response.json())
    .then((data) => {
      if (data.Error) {
        throw new Error(data.Error.message)
      }
      // handle success response
      console.log(data);
      this.setState({submitState: "SUCCESS"});
      console.log("Successfully started campaign");
    }).catch(err => {
      this.setState({submitState: "ERROR", err: err.message});
      console.warn("Failed to send campaign");
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
          <FormGroup controlId="title" bsSize={FORM_ELEMENT_SIZE}>
            <ControlLabel>Title</ControlLabel>
            <FormControl
              autoFocus
              value={this.state.notification.title}
              onChange={this.handleChange}
            />
          </FormGroup>
          <FormGroup controlId="message" bsSize={FORM_ELEMENT_SIZE}>
            <ControlLabel>Message</ControlLabel>
            <FormControl
              value={this.state.notification.message}
              onChange={this.handleChange}
            />
          </FormGroup>
          <FormGroup controlId="templatePath" bsSize={FORM_ELEMENT_SIZE}>
            <ControlLabel>Template</ControlLabel>
            <FormControl
              value={this.state.notification.templatePath}
              onChange={this.handleChange}
            />
          </FormGroup>
          <FormGroup controlId="groupId" bsSize={FORM_ELEMENT_SIZE}>
            <ControlLabel>Group ID</ControlLabel>
            <FormControl
              value={this.state.notification.groupId}
              onChange={this.handleChange}
            />
          </FormGroup>
          <FormGroup controlId="runId" bsSize={FORM_ELEMENT_SIZE}>
            <ControlLabel>Run ID</ControlLabel>
            <FormControl
              value={this.state.notification.runId}
              onChange={this.handleChange}
            />
          </FormGroup>
          <FormGroup controlId="templateMetadata" bsSize={FORM_ELEMENT_SIZE}>
            <ControlLabel>Data</ControlLabel>
            <JSONInput
              id='notification-content-editor'
              locale={ locale }
              onChange={this.handleChangeJsonBox}
              placeholder={this.state.notification.templateMetadata}
              value={this.state.notification.templateMetadata}
              height='100px'
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
