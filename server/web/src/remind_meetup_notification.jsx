import ReactDOM from 'react-dom';
import React from 'react';
import getProperty from './context.jsx';
import './scss/notification.scss';
import { identifyUser, trackNotificationOpened, trackAction, trackLinkClicked } from './metrics/mixpanel';
import { serverUrl } from './config.js'

const hiveDeepLinkRoot = 'hive:/';
//const hiveDeepLinkRoot = 'exp://192.168.0.179:19000/--'; // For dev

const defaultProfilePic =
  'https://t4.ftcdn.net/jpg/02/15/84/43/240_F_215844325_ttX9YiIIyeaR7Ne6EaLLjMAmy4GvPC69.jpg';
const meetupReminderUrl = `${serverUrl}/v1/meetup_reminder`;
const homeViewLink = `${hiveDeepLinkRoot}/Home`;

function switchMatchType(matchType, mentor, mentee, connection) {
  if (matchType === 'MENTOR') {
    return mentor;
  } else if (matchType === 'MENTEE') {
    return mentee;
  } else if (matchType === 'CONNECTION') {
    return connection;
  }
}

const NOTIFICATION_NAME = 'remind-meetup';

const VIEW_DEFAULT = 'default';
const VIEW_DID_MEET = 'did_meet';
const VIEW_DID_NOT_MEET = 'did_not_meet';
const VIEW_DONT_REMIND = 'no_dont_remind';
const VIEW_REMINDER_CONFIRMED = 'reminder_confirmed';

class RemindMeetupNotification extends React.Component {
  /*
   * state: {
   *  view: string
   *  followupIn: int (weeks)
   *  errorMsg: string
   * }
   */

  constructor() {
    super();
    this.state = { view: VIEW_DEFAULT, followupIn: 1, errorMsg: '' };

    this.render = this.render.bind(this);
    this.renderDefaultView = this.renderDefaultView.bind(this);
    this.renderDidMeetView = this.renderDidMeetView.bind(this);
    this.renderDidNotMeetView = this.renderDidNotMeetView.bind(this);
    this.renderReminderConfirmedView = this.renderReminderConfirmedView.bind(this);
    this.renderDontRemindView = this.renderDontRemindView.bind(this);
    this.postMeetupReminder = this.postMeetupReminder.bind(this);
    this.deleteMeetupReminder = this.deleteMeetupReminder.bind(this);
  }

  componentDidMount() {
    const user = getProperty('user');
    identifyUser(user.UserId);
    trackNotificationOpened(NOTIFICATION_NAME);
    trackLinkClicked(NOTIFICATION_NAME, '#view-match-profile', 'view-match-profile');
  }

  postMeetupReminder(userId, matchUserId, reminderTime, final) {
    fetch(meetupReminderUrl, {
      method: 'POST',
      body: JSON.stringify({
        userId,
        matchUserId,
        reminderTime,
      })
    }).catch(err => {
        console.warn(err);
        this.setState({errorMsg: err})
    }).finally(() => final());
  }

  deleteMeetupReminder(userId, matchUserId, final) {
    fetch(meetupReminderUrl, {
      method: 'DELETE',
      body: JSON.stringify({
        userId,
        matchUserId,
      })
    }).catch(err => {
        console.warn(err);
        this.setState({errorMsg: err})
    }).finally(() => final());
  }

  renderDefaultView(matchType, meetupType, matchFirstName, matchLastName) {
    let bodyText = '';
    let title = '';
    let matchFullName = matchFirstName + ' ' + matchLastName;
    if (meetupType === 'INITIAL_MEETING') {
      title = 'Meet up with ' + matchFirstName;
      bodyText = switchMatchType(matchType,
        'Have you met up with your mentor ' + matchFullName + ' yet?',
        'Have you met up with your mentee ' + matchFullName + ' yet?',
        'Have you met up with your connection ' + matchFullName + ' yet?',
      );
    } else if (meetupType === 'FOLLOWUP_MEETING') {
      title = 'Follow up with ' + matchFirstName;
      bodyText = switchMatchType(matchType,
        'You wanted us to remind you to follow up with your mentor, ' +
        matchFullName + '. Did you meet up with them again?',
        'You wanted us to remind you to follow up with your mentee, ' +
        matchFullName + '. Did you meet up with them again?',
        'You wanted us to remind you to follow up with your connection, ' +
        matchFullName + '. Did you meet up with them again?',
      );
    }

    return (
      <div className="content">
        <p className="title">{ title }</p>
        <p className="message">{bodyText}</p>
        <a
          id="btn-did-meet"
          className="btn-primary"
          onClick={() => {
            trackAction(NOTIFICATION_NAME, "click-did-meet");
            this.setState({view:VIEW_DID_MEET});
          }}>
          Yes, we met
        </a>
        <a
          id="btn-did-not-meet"
          className="btn-decline"
          onClick={() => {
            trackAction(NOTIFICATION_NAME, "click-did-not-meet");
            this.setState({view:VIEW_DID_NOT_MEET});
          }}>
          No, not yet
        </a>
      </div>
    );
  }

  renderDidMeetView(userId, matchUserId, matchType, meetupType, matchFirstName) {
    let bodyText = '';
    let title = 'Follow up with ' + matchFirstName;
    if (meetupType === 'INITIAL_MEETING') {
      bodyText = switchMatchType(matchType,
        'Awesome! If you got along well with your mentor you should check in with' +
          ' them again soon, even if you don\'t have specific questions.' +
          ' Get reminded to follow up in:',
        'Awesome! Mentorships only work if you check up on your mentee every once in' +
          ' a while. Get reminded to follow up in:',
        'Great! Meeting up in person strengthens your connections like nothing else.' +
          ' If you like, we can remind you to follow up in:'
      );
    } else if (meetupType === 'FOLLOWUP_MEETING') {
      bodyText = switchMatchType(matchType,
        'Wow, you met with your mentor again! We hope the connection doesn\'t end here -' +
          ' get reminded to follow up in:',
        'You\'re a mentoring superstar! We know there\'s still a ton your mentee can learn' +
          ' from you. Get reminded to follow up again in:',
        'Nice, sounds like you\'re building a strong connection! Don\'t let it end here' +
          ' though, get reminded to follow up again in:',
      );
    }

    let selector_weeks = [1, 2, 3, 4];
    return (
      <div className="content">
        <p className="title">{ title }</p>
        <p className="message">{bodyText}</p>
        <select
          id="select-follow-up-interval"
          className="btn-primary"
          value={this.state.followupIn}
          onChange={(e) => { this.setState({followupIn: e.target.value}); }}>
          { selector_weeks.map((n) => <option key={n} value={n}>{n} Week{n>1?'s':''}</option>) }
        </select>
        <a
          id="btn-follow-up"
          className="btn-primary"
          onClick={() => {
            const {followupIn} = this.state; // weeks
            const millisPerWeek = 7*24*60*60*1000;
            const reminderTime = new Date(new Date().getTime() + (followupIn*millisPerWeek));
            this.postMeetupReminder(userId, matchUserId, reminderTime, () => {
              this.setState({view:VIEW_REMINDER_CONFIRMED})
            });
          }}>
          Submit
        </a>
        <a
          id="btn-dont-follow-up"
          className="btn-decline"
          onClick={() => {
            trackAction(NOTIFICATION_NAME, "click-no-follow-up");
            this.deleteMeetupReminder(userId, matchUserId, () => {
              this.setState({view:VIEW_DONT_REMIND})
            });
          }}>
          Don't follow up
        </a>
      </div>
    );
  }

  renderDidNotMeetView(userId, matchUserId, matchType, meetupType, matchProfileLink, matchFirstName, matchLastName) {
    let bodyText = '';
    let title = '';
    let matchFullName = matchFirstName + ' ' + matchLastName;
    let things_to_talk_about = [];
    if (meetupType === 'INITIAL_MEETING') {
      title = 'Meet up with ' + matchFirstName;
      bodyText = switchMatchType(matchType,
        'No time like the present! Don\'t be afraid to reach out to your mentor first.' +
        ' They signed up for the program and are happy to meet up with you sometime.' +
        ' We\'ll check in again in a few days, consider scheduling a meeting before then.' +
        ' Here are some things you might want to talk about at your first meeting:',
        'Meeting up in person is the cornerstone of a great mentorship. Your mentees often' +
        ' won\'t reach out even if they need to, so it would be amazing if you can contact' +
        ' them first. We\'ll check in with you in a few days, consider grabbing a coffee or' +
        ' bubble tea together before then. Here are some things you might want to talk' +
        ' about at your first meeting:',
        'Meeting up in person is the best way to strengthen new connections. We\'ll check in' +
        ' with you again in a few days, consider grabbing a quick coffee together before then.' +
        ' Or, if you matched on something you both like doing, why not do that?',
      );
      things_to_talk_about = switchMatchType(matchType,
        [
          'Why did you choose your program? What do you hope to get out of it?',
          'Anything stressing you out at the moment?',
          'Are there any hobbies your mentor picked up while at Waterloo?',
        ],
        [
          'How do you manage stress?',
          'What are some of your high points at Waterloo or on co-op?',
          'Are there questions you remember having back in first year?',
          'Ask open-ended questions and take the time to actively listen to your mentee',
        ],
        [],
      );
    } else if (meetupType === 'FOLLOWUP_MEETING') {
      title = 'Follow up with ' + matchFirstName;
      bodyText = switchMatchType(matchType,
        'No worries, we\'ll check back in with you in a bit - try to follow up in the' +
        ' next few days. Here are some things you might want to talk about at your next meeting:',
        'Following up with your mentee and seeing their progress is a big part of what' +
        ' makes mentorship so helpful. We\'ll check back in a few days, consider having your' +
        ' next meeting before then. Here are some things you can do at your next meeting:',
        'We\'re sure your connection wants to meet again just as much as you do! See if they' +
        ' want to get together sometime in the next few days - we\'ll check again then.',
      );
      things_to_talk_about = switchMatchType(matchType,
        [
          'What are your dreams for the future? Any particular companies or industries you\'re' +
          ' interested in working in?',
          'Any places you want to travel? Bucket list items? Interested in going on exchange?',
          'How has your mental wellbeing been recently? Many students go through periods of' +
          ' struggle for various reasons - your mentor has been through this too',
        ],
        [
          'Figure out your mentee\'s goals and help them set realistic steps to reach them',
          'Hold them accountable to their planned steps. Encourage and praise them when they' +
          ' make progress',
          'Connect them with the people or resources that can help them achieve their dreams',
          'Give some examples of how you\'ve overcome tough obstacles, made hard decisions,' +
          ' and grown since you were in your mentee\'s position',
        ],
        [],
      );
    }

    return (
      <div className="content">
        <p className="title">{ title }</p>
        <p className="message">{bodyText}</p>
        {
          things_to_talk_about && <ul>
            { things_to_talk_about.map((s, i) => <li key={i}>{s}</li>)}
          </ul>
        }
        <a
          id="view-match-profile"
          className="btn-primary"
          onClick={() => { document.location.href = matchProfileLink }}>
          View {matchFullName}'s Profile
        </a>
        <a
          id="btn-dont-remind-me"
          className="btn-decline"
          onClick={() => {
            trackAction(NOTIFICATION_NAME, "click-no-reminder");
            this.deleteMeetupReminder(userId, matchUserId, () => {
              this.setState({view:VIEW_DONT_REMIND});
            });
          }}>
          Don't remind me again
        </a>
      </div>
    );
  }

  renderReminderConfirmedView(matchType, matchFirstName, matchLastName) {
    const { followupIn, errorMsg } = this.state;
    const matchFullName = matchFirstName + ' ' + matchLastName;
    const title = 'Reminder confirmed';
    const bodyText = "Ok, we'll remind you to meet up with your " +
      switchMatchType(matchType, "mentor", "mentee", "connection") + " " + matchFullName +
      " in " + followupIn + " week" + (followupIn > 1 ? "s." : ".");

    return (
      <div className="content">
        <p className="title">{ title }</p>
        <p className="message">{bodyText}</p>
        { errorMsg &&
        <p className="error-text">
          Error communicating with server, please try again later.
        </p>}
        <a
          id="btn-close"
          className="btn-decline"
          href={homeViewLink}
          >
          Close
        </a>
      </div>
    );
  }

  renderDontRemindView(matchFirstName, matchLastName) {
    const { errorMsg } = this.state;
    let title = 'Reminder cancelled';
    let matchFullName = matchFirstName + ' ' + matchLastName;
    let bodyText = 'We\'ll no longer remind you to meet up with ' + matchFullName + '.';

    return (
      <div className="content">
        <p className="title">{ title }</p>
        <p className="message">{bodyText}</p>
        { errorMsg &&
        <p className="error-text">
          Error communicating with server, please try again later.<br/>
          Message: {errorMsg}
        </p>}
        <a
          id="btn-close"
          className="btn-decline"
          href={homeViewLink}>
          Close
        </a>
      </div>
    );
  }

  /* Test:
  {
    matchType: 'MENTOR',
    meetupType: 'INITIAL_MEETING',
    matchUserId: 1,
    matchFirstName: 'Foo',
    matchLastName: 'Bar'
  }
   */

  render() {
    const user = getProperty('user');
    const matchType = getProperty('matchType');
    const meetupType = getProperty('meetupType');
    const matchUserId = getProperty('matchUserId');
    const matchFirstName = getProperty('matchFirstName');
    const matchLastName = getProperty('matchLastName');
    const matchProfilePic = getProperty('matchProfilePic') || defaultProfilePic;
    const profileLink = `${hiveDeepLinkRoot}/MatchProfile?userId=${matchUserId}`;

    const { view } = this.state;

    let body = <div/>;
    if (view === VIEW_DEFAULT) {
      body = this.renderDefaultView(matchType, meetupType, matchFirstName, matchLastName);
    } else if (view === VIEW_DID_MEET) {
      body = this.renderDidMeetView(user.UserId, matchUserId, matchType, meetupType, matchFirstName);
    } else if (view === VIEW_DID_NOT_MEET ) {
      body = this.renderDidNotMeetView(user.UserId, matchUserId, matchType, meetupType, profileLink, matchFirstName, matchLastName);
    } else if (view === VIEW_REMINDER_CONFIRMED ) {
      body = this.renderReminderConfirmedView(matchType, matchFirstName, matchLastName)
    } else if (view === VIEW_DONT_REMIND ) {
      body = this.renderDontRemindView(matchFirstName, matchLastName)
    }

    return (
      <div className="notification-background">
        <img className="profile-pic" alt="" border="0" height="100" src={ matchProfilePic } />
        <div className="notification-container">
          { body }
        </div>
    </div>
    );
  }
}

ReactDOM.render(<RemindMeetupNotification/>, document.getElementById('content'));
