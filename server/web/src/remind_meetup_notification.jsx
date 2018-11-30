import ReactDOM from 'react-dom';
import React from 'react';
import getProperty from './context.jsx';
import './scss/notification.scss';

function getIntro(connectionName, foodSuggestion) {
  return 'An important part of any great mentorship is in-person interaction! ' +
    'If you haven\'t met up with ' + connectionName + ', you should set something up! ' +
    foodSuggestion + ' Some good topics of conversation include:';
}

const defaultProfilePic =
  'https://t4.ftcdn.net/jpg/02/15/84/43/240_F_215844325_ttX9YiIIyeaR7Ne6EaLLjMAmy4GvPC69.jpg';

const mentorFoodSuggestion =
  'Take your mentee out for Waterloo specialities such as bubble tea or shawarma.';
const mentorTopics = [
  'University lifehacks',
  'Life outside of school',
  'Memorable experiences so far',
  'Advice on dealing with stress',
];
const mentorOutro =
  'Chances are you\'ve been in their shoes, why not help them out? ' +
  'What would you have wished you knew in your first years at Waterloo?';

const menteeFoodSuggestion =
  'Ask your mentor to take you to a new food or bubble tea joint.';
const menteeTopics = [
  'Favourite memory from your first years at Waterloo',
  'Life outside of school',
  'Memorable experiences so far',
  'Why you chose your program',
];
const menteeOutro =
  'Chances are your mentor has been in your shoes, here\'s an opportunity to look into the future!';

class RemindMeetupNotification extends React.Component {
  render() {
    const userType = getProperty('userType');
    const connectionUserId = getProperty('connectionUserId');
    const connectionFirstName = getProperty('connectionFirstName');
    const connectionLastName = getProperty('connectionLastName');
    const connectionName = connectionFirstName + ' ' + connectionLastName;
    const connectionProfilePic = getProperty('connectionProfilePic') || defaultProfilePic;
    const link = 'hive://MatchProfile/' + connectionUserId;

    let foodSuggestion = '';
    let outro = '';
    let topics = [];
    if (userType === 'MENTOR') {
      foodSuggestion = mentorFoodSuggestion;
      outro = mentorOutro;
      topics = mentorTopics;
    } else {
      foodSuggestion = menteeFoodSuggestion;
      outro = menteeOutro;
      topics = menteeTopics;
    }

    return (
      <div className="notification-background">
        <img className="profile-pic" alt="" border="0" height="100" src={ connectionProfilePic } />
        <div className="notification-container">
          <div className="content">
            <p className="title">Meet up with { connectionFirstName }</p>
            <p className="message">
              { getIntro(connectionName, foodSuggestion) }
              <ul>
                { topics.map(topic => <li>{ topic }</li>) }
              </ul>
              { outro }
            </p>
            <a className="cta" href={link}>Contact them now!</a>
          </div>
        </div>
    </div>
    );
  }
}

ReactDOM.render(<RemindMeetupNotification />, document.getElementById('content'));
