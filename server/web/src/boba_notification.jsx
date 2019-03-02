import ReactDOM from 'react-dom';
import React from 'react';
import getProperty from './context.jsx';
import './scss/notification.scss';
import { identifyUser, trackNotificationOpened, trackLinkClicked } from './metrics/mixpanel';

const NOTIFICATION_NAME = (s) => `generic-notification`;

const mentorTitle = genericOutro("mentee");
const menteeTitle = genericOutro("mentor");

class GenericMessageNotification extends React.Component {

  componentDidMount() {
    const user = getProperty('user');
    // whether the user is a mentor or mentee
    const userType = getProperty('userType');
    identifyUser(user.UserId);
    trackNotificationOpened(NOTIFICATION_NAME);
  }

  render() {
    let title;
    if (userType === 'MENTOR') {
      title = mentorTitle;
    } else {
      title = menteeTitle;
    }

    const message = getProperty('message');
    return (
      <div className="notification-background">
        <div className="notification-container">
          <div className="content">
            <p className="title">{title}</p>
            <p className="message">{message}</p>
          </div>
        </div>
    </div>
    );
  }
}

ReactDOM.render(<GenericMessageNotification />, document.getElementById('content'));
