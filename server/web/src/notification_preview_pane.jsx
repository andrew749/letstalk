import {withCookies} from 'react-cookie';
import React from 'react';
import {serverUrl} from './config.js';

class NotificationPreviewPane extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      htmlContent:  null
    }
  }
  fetchContent(sessionId, template, data) {
    fetch(`${serverUrl}/echo_notification`, {
      method: 'POST',
      headers: {'sessionId': sessionId},
      body: JSON.stringify({
        templateLink: template,
        data: data
      })
    })
    .then((response) => response.text())
    .then((data) => {
      // handle success response
      this.setState({
        htmlContent: data
      });
    }).catch(err => {
      console.log(err);
    });
  }

  componentDidUpdate(prevProps) {
    console.log("updateing")
    if (this.props.template !== prevProps.template || this.props.data !== prevProps.data) {
      console.log("can update")
      const {template, data, sessionId} = this.props;
      this.fetchContent(sessionId, template, data);
    }
  }

  render() {
    return <iframe className="notification-render-preview"srcDoc={this.state.htmlContent}></iframe>;
  }
}

const CookieNotificationPreviewPane = withCookies(NotificationPreviewPane);

export default CookieNotificationPreviewPane;
