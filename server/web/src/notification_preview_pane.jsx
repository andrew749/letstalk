import {withCookies} from 'react-cookie';
import React from 'react';
import {serverUrl} from './config.js';

const PREVIEW_PANE_WIDTH = "320";
const PREVIEW_PANE_HEIGHT = "480";

class NotificationPreviewPane extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      htmlContent:  null
    }
  }
  fetchContent(sessionId, templatePath, templateMetadata) {
    fetch(`${serverUrl}/v1/echo_notification`, {
      method: 'POST',
      headers: {'sessionId': sessionId},
      body: JSON.stringify({
        templateLink: templatePath,
        data: templateMetadata
      })
    })
    .then((response) => response.text())
    .then((data) => {
      // handle success response
      this.setState({
        htmlContent: data
      });
    }).catch(err => {
      console.warn(err);
    });
  }

  componentDidUpdate(prevProps) {
    if (this.props.templatePath !== prevProps.templatePath || this.props.templateMetadata !== prevProps.templateMetadata) {
      const {templatePath, templateMetadata, sessionId} = this.props;
      this.fetchContent(sessionId, templatePath, templateMetadata);
    }
  }

  render() {
    return (
      <iframe
        width={PREVIEW_PANE_WIDTH}
        height={PREVIEW_PANE_HEIGHT}
        className="notification-render-preview"
        srcDoc={this.state.htmlContent}>
      </iframe>
    );
  }
}

const CookieNotificationPreviewPane = withCookies(NotificationPreviewPane);

export default CookieNotificationPreviewPane;
