import React, { Component, ReactNode } from 'react';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import {
  View,
  Text,
  RefreshControl,
  RefreshControlProps,
  ScrollView,
  StyleSheet,
} from 'react-native';
import {
  NavigationScreenProp,
  NavigationScreenDetails,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import { ThunkAction } from 'redux-thunk';
import Moment from 'moment';
import { MaterialIcons } from '@expo/vector-icons';

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
import Colors from '../services/colors';

interface DispatchActions {
  fetchNewestNotifications:
    ActionCreator<ThunkAction<Promise<NotificationsActionTypes>, NotificationsState, void>>;
}

interface Props extends NotificationsState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

interface State {
  refreshing: boolean;
}

class NotificationView extends Component<Props, State> {
  NOTIFICATIONS_VIEW_IDENTIFIER = "HomeView";

  static navigationOptions = {
    headerTitle: 'Notifications',
    headerStyle,
  }

  constructor(props: Props) {
    super(props);

    this.state = { refreshing: false };

    this.load = this.load.bind(this);
    this.renderBody = this.renderBody.bind(this);
    this.onRefresh = this.onRefresh.bind(this);
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.NOTIFICATIONS_VIEW_IDENTIFIER);
    });

    this.load();
  }

  private async load() {
    await this.props.fetchNewestNotifications();
  }

  private async onRefresh() {
    this.setState({refreshing: true});
    await this.load();
    this.setState({refreshing: false});
  }

  renderNotification(notification: Notification) {
    const {
      notificationId,
      state,
      data,
      createdAt,
    } = notification;

    let notifText: ReactNode = null;
    let icon = 'face';

    switch (notification.type) {
      case 'NEW_CREDENTIAL_MATCH':
        const pronoun = data.side === 'ASKER' ? 'You' : 'They';
        notifText = (
          <Text>
            {`You were matched with `}
            <Text style={{fontWeight: 'bold'}}>
              {`${data.userName}`}
            </Text>
            {`! ${pronoun} requested the credential `}
            <Text style={{fontWeight: 'bold'}}>
              {`"${data.credentialName}"`}
            </Text>
          </Text>
        );
        icon = 'people';
        break;
      default:
        // Ensure exhaustiveness of select
        const _: never = notification.type;
    }

    const containerStyle = [styles.notifContainer];
    if (state === 'UNREAD') containerStyle.push(styles.notifContainerUnread);

    return (
      <View key={notification.notificationId} style={containerStyle}>
        <View style={styles.leftContainer}>
          <MaterialIcons name={icon} size={48} />
        </View>
        <View style={styles.rightContainer}>
          {notifText}
          <Text style={styles.ago}>{Moment(createdAt).fromNow()}</Text>
        </View>
      </View>
    );
  }

  renderBody() {
    const notifs = this.props.notifications.map(this.renderNotification).toJS();
    return (
      <ScrollView
        refreshControl={
          <RefreshControl
            refreshing={this.state.refreshing}
            onRefresh={this.onRefresh}
          /> as React.ReactElement<RefreshControlProps>
        }
      >
        <View style={styles.container}>
          {notifs}
        </View>
      </ScrollView>
    );
  }

  render() {
    const {
      state,
      errorMsg,
      errorType,
    } = this.props.fetchState;
    // If `this.state.refreshing` is true, it means that we are reloading data using the pull
    // down, which means that we want to still display the ScrollView.
    return (
      <Loading
        state={this.state.refreshing ? 'success' : state}
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
  ago: {
    paddingTop: 5,
    color: Colors.DARK_GRAY,
  },
  container: {
    borderTopWidth: 0.5,
    borderColor: Colors.HIVE_SUBDUED,
  },
  notifContainer: {
    backgroundColor: Colors.WHITE,
    padding: 10,
    flexDirection: 'row',
    borderBottomWidth: 0.5,
    borderColor: Colors.HIVE_SUBDUED,
  },
  notifContainerUnread: {
    backgroundColor: Colors.HIVE_PRIMARY_LIGHT,
  },
  leftContainer: {
    paddingRight: 10,
    justifyContent: 'center',
    alignItems: 'center',
    flex: 1,
  },
  rightContainer: {
    flex: 5,
  },
});
