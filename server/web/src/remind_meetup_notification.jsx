import ReactDOM from 'react-dom';
import React from 'react';
import getProperty from './context.jsx';
import './scss/notification.scss';

const intro = 'An important part of any great mentorship is in-person interaction!';

class RemindMeetupNotification extends React.Component {
  render() {
    const userType = getProperty('userType');
    const connectionUserId = getProperty('connectionUserId');
    const connectionFirstName = getProperty('connectionFirstName');
    const connectionLastName = getProperty('connectionLastName');
    const connectionName = connectionFirstName + ' ' + connectionLastName;
    const connectionProfilePic = getProperty('connectionProfilePic');
    let message = intro + ' If you haven\'t met up with ' + connectionName + ', you should set ' +
        'something up!';
    if (userType === 'MENTOR') {
      message = message + ' Maybe take your mentee out for bubble tea or shawarma. Some good topics of ' +
        'conversation include: university tips and tricks, things you do outside of school and ' +
        'crazy stories from your university experience so far. Make sure to ask if they need any ' +
        'advice - chances are, you\'ve already went through a lot of what they are going through ' +
        'right now.';
    } else {
      message = message + ' Suggest to your mentor that you want to check out a new food place. ' +
        'Some good topics of conversation include: favourite memory from first year, things you ' +
        'do outside of school and why you chose your program. Make sure to ask lots of questions - ' +
        'chances are, your mentor has went through a lot of what you are currently going through.';
    }
    const link = 'hive://MatchProfile/' + connectionUserId;
    return (
      <div className="notification-background">
        <img className="profile-pic" alt="" border="0" height="100" src={ connectionProfilePic } />
        <div className="notification-container">
          <div className="content">
            <p className="title">Meet up with { connectionFirstName }</p>
            <p className="message">{ message }</p>
            <a className="cta" href={link}>View Profile</a>
          </div>
        </div>
    </div>
    );
  }
}

ReactDOM.render(<RemindMeetupNotification />, document.getElementById('content'));
