import React, { Component } from 'react';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import { View, Text, ScrollView, StyleSheet } from 'react-native';
import {
  NavigationScreenProp,
  NavigationScreenDetails,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import { ThunkAction } from 'redux-thunk';

import { AnalyticsHelper } from '../services/analytics';
import { RootState } from '../redux';
import {
  State as NotificationsState,
  fetchNewestNotifications,
} from '../redux/notifications/reducer';
import { ActionTypes as NotificationsActionTypes } from '../redux/notifications/actions';
import Loading from './Loading';
import { headerStyle } from './TopHeader';
import { Notification } from '../models/notification';

interface DispatchActions {
  fetchNewestNotifications:
    ActionCreator<ThunkAction<Promise<NotificationsActionTypes>, NotificationsState, void>>;
}

interface Props extends NotificationsState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class NotificationView extends Component<Props> {
  NOTIFICATIONS_VIEW_IDENTIFIER = "HomeView";

  static navigationOptions = {
    headerTitle: 'Notifications',
    headerStyle,
  }

  constructor(props: Props) {
    super(props);

    this.load = this.load.bind(this);
    this.renderBody = this.renderBody.bind(this);
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.NOTIFICATIONS_VIEW_IDENTIFIER);
    });

    this.load();
  }

  async load() {
    await this.props.fetchNewestNotifications();
  }

  renderNotification(notification: Notification) {
    const {
      notificationId,
      type,
      state,
    } = notification;

    return (
      <View style={styles.notifContainer}>
        <View style={styles.leftContainer}>
          <Text key={notification.notificationId}>{notification.type}</Text>;
        </View>
        <View style={styles.rightContainer}>
          <Text key={notification.notificationId}>{notification.type}</Text>;
        </View>
      </View>
    );
  }

  renderBody() {
    const notifs = this.props.notifications.map(this.renderNotification).toJS();
    return <ScrollView>{notifs}</ScrollView>;
  }

  render() {
    const {
      state,
      errorMsg,
      errorType,
    } = this.props.fetchState;

    return (
      <Loading
        state={state}
        errorMsg={errorMsg}
        errorType={errorType}
        load={this.load}
        renderBody={this.renderBody}
        navigation={this.props.navigation}
      />
    );
  }
}

export default connect(({ notifications }: RootState) => notifications,
  { fetchNewestNotifications })(NotificationView);

const styles = StyleSheet.create({
  notifContainer: {
    backgroundColor: 'white',
    padding: 10,
    flexDirection: 'row',
    borderBottomWidth: 0.5,
    borderColor: 'gray',
  },
  leftContainer: {
    flex: 1,
    backgroundColor: 'red',
  },
  rightContainer: {
    flex: 4,
    backgroundColor: 'green',
  },
});
