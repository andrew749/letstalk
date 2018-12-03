import ReactDOM from 'react-dom';
import React from 'react';
import getProperty from './context.jsx';
import './scss/notification.scss';
import { identifyUser, trackNotificationOpened, trackLinkClicked } from './metrics/mixpanel';

const NOTIFICATION_NAME = 'notif-with-quote';

class NotificationWithQuote extends React.Component {
  componentDidMount() {
    const user = getProperty('user');
    const link = getProperty('link');
    const cta = getProperty('cta');
    identifyUser(user.UserId);

    const properties = {
      userId: user.UserId,
      firstName: user.FirstName,
      lastName: user.LastName,
      link: link,
      cta: cta,
    };

    trackNotificationOpened(NOTIFICATION_NAME, properties);
    trackLinkClicked(NOTIFICATION_NAME, '#cta', 'cta', properties);
  }

  render() {
    const user = getProperty('user');
    const quotes = getProperty('quotes');
    const body = getProperty('body');
    const link = getProperty('link');
    const cta = getProperty('cta');
    const quote = quotes[Math.floor(Math.random()*quotes.length)];
    return (
      <div className="notification-background">
        <a href="https://hiveapp.org"><img alt="" border="0" height="163" src="https://marketing-image-production.s3.amazonaws.com/uploads/0f48be323dc73ee6e531f2fa9752b4c35be39478b7c83f7a9747262b9f60abb856ad92305e465504297eebf994da0407361d56b2c35d07ef72ca884f167b4de0.png" width="100"/></a>
        <div className="notification-container">
          <div className="content">
            <p className="title">Hey {user.FirstName}!</p>

            <p className="quote">
            “{quote.body}”
            </p>

            <p className="author">
              ~{quote.author}
            </p>

            <p className="message">
              {body}
            </p>

            <a id="cta" className="cta" href={link}>
              {cta}
            </a>
          </div>
        </div>
    </div>
    );
  }
}

ReactDOM.render(<NotificationWithQuote />, document.getElementById('content'));
