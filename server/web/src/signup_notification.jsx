import ReactDOM from 'react-dom';
import React from 'react';
import getProperty from './context.jsx';
import './scss/notification.scss';

export default class SignupNotification extends React.Component {
    render() {
      const title = getProperty('title');
      const body = getProperty('body');
      const thumbnail = getProperty('thumbnail');
      return (
        <div className="notification-background">
          <div className="notification-container">
          <div className="content">
            <p className="title">Welcome to Hive!</p>
            <p>Welcome to Hive, a platform for connecting bright and promising students to one another! Having been in your shoes, we know the struggles of needing advice and not knowing where to turn. Mentorship is key in navigating through the challenging, fast-paced university environment and beyond.</p>
            <br />
            <p>Through interacting with a mentor, you can learn from others' experience. Conversely, being a mentor to someone is extremely rewarding, allowing you to help out others who are facing the same challenges that you struggled with in the past.</p>
            <br />
            <p>In our increasingly online world, it's all the more important to maintain our real physical connections; there's something indispensable about socializing in person. Even if it's just a excuse to get a drink in the middle of a busy day, we'd love to see you take a leap of faith and see what you uncover with your connections!</p>
            <br />
            <p>- The Hive Team</p>
            </div>
          </div>
        </div>
      );
    }
}

ReactDOM.render(<SignupNotification />, document.getElementById('content'));
