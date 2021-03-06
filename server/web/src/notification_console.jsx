import ReactDOM from 'react-dom';
import React from 'react';
import CookieAwareComponent from './cookie_aware_component.jsx'
import {withCookies} from 'react-cookie';
import NotificationPreviewPane from './notification_preview_pane.jsx';
import NotificationForm from './notification_form.jsx';
import './scss/notification_console.scss';
import {HiveToolTemplate} from './hive_tool_template.jsx';

class NotificationConsole extends React.Component {
  constructor(props) {
    super(props)
    this.onUpdateData = this.onUpdateData.bind(this);
    this.state = {
      // to be sent to server to echo back with information
      notificationState: {}
    }
  }

  onUpdateData(newNotificationData) {
    this.setState({notificationState: newNotificationData});
  }

  render() {
    const {cookies} = this.props;
    return (
    <HiveToolTemplate>
      <div className="notification-console-container">
        <NotificationForm
          onUpdateData={this.onUpdateData}/>
        <NotificationPreviewPane
          sessionId={cookies.get('sessionId')}
          templatePath={this.state.notificationState.templatePath}
          templateMetadata={this.state.notificationState.templateMetadata}/>
      </div>
      </HiveToolTemplate>
    );
  }
}

const CookieAwareNotificationConsole = CookieAwareComponent(withCookies(NotificationConsole));

ReactDOM.render(<CookieAwareNotificationConsole/>, document.getElementById('content'));
