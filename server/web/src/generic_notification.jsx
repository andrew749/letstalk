import ReactDOM from 'react-dom';
import React from 'react';
import getProperty from './context.jsx';
import './scss/notification.scss';
import { identifyUser, trackNotificationOpened, trackLinkClicked } from './metrics/mixpanel';

const NOTIFICATION_NAME = (s) => (!!s) ? `generic-notification-${s}` : `generic-notification`;

class GenericMessageNotification extends React.Component {

  componentDidMount() {
    const user = getProperty('user');
    const notificationTracking = NOTIFICATION_NAME(getProperty('trackingId'));
    identifyUser(user.UserId);
    trackNotificationOpened(notificationTracking);
  }

  render() {
    const caption = getProperty('caption');
    const body = getProperty('body')
    return (
      <div className="notification-background">
        <div className="notification-container">
          <div className="content">
            <p className="title">{caption}</p>
            <p className="message">{body}</p>
          </div>
        </div>
    </div>
    );
  }
}

ReactDOM.render(<GenericMessageNotification />, document.getElementById('content'));
