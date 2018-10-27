import ReactDOM from 'react-dom';
import React from 'react';
import getProperty from './context.jsx';
import './scss/notification.scss';

class ExploreNotification extends React.Component {
    render() {
      const user = getProperty('user');
      return (
        <div className="notification-container">
          <div className="content">
            <p>Hey {user.FirstName}!</p>

            <p>
            Looks like you've been using Hive for a few weeks.
            </p>

            <p>
            Did you know that 90% of X does Y.
            </p>
          </div>
        </div>
      );
    }
}

ReactDOM.render(<ExploreNotification/>, document.getElementById('content'));
