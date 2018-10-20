import ReactDOM from 'react-dom';
import React from 'react';
import CookieAwareComponent from './cookie_aware_component.jsx'
import {withCookies} from 'react-cookie';
import NotificationPreviewPane from './notification_preview_pane.jsx';
import NotificationForm from './notification_form.jsx';
import './scss/notification_console.scss';

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
      <div className="notification-console-container">
        <NotificationForm
          onUpdateData={this.onUpdateData}/>
        <NotificationPreviewPane
          sessionId={cookies.get('sessionId')}
          template={this.state.notificationState.template}
          data={this.state.notificationState.data}/>
      </div>
    );
  }
}

const CookieAwareNotificationConsole = CookieAwareComponent(withCookies(NotificationConsole));

ReactDOM.render(<CookieAwareNotificationConsole/>, document.getElementById('content'));
