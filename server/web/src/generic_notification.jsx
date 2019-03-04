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

  // dont do this at home kids
  body(isHTML, body) {
    if (!!isHTML) {
      return <div dangerouslySetInnerHTML={{__html: body}}></div>;
    } else {
      return <p className="message">{body}</p>;
    }
  }

  render() {
    const caption = getProperty('caption');
    const bodyContent = getProperty('body')
    const bodyIsHTML = getProperty('bodyIsHTML');
    return (
      <div className="notification-background">
        <div className="notification-container">
          <div className="content">
            <p className="title">{caption}</p>
            {this.body(bodyIsHTML, bodyContent)}
          </div>
        </div>
    </div>
    );
  }
}

ReactDOM.render(<GenericMessageNotification />, document.getElementById('content'));
