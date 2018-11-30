import React, { Component, ReactNode } from 'react';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import {
  Dimensions,
  ListView,
  View,
  Text,
  Image,
  RefreshControl,
  RefreshControlProps,
  StyleSheet,
  TouchableOpacity,
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
import Immutable from 'immutable';

import { Button } from '../components';
import { AnalyticsHelper } from '../services/analytics';
import { RootState } from '../redux';
import { errorToast } from '../redux/toast';
import {
  State as NotificationsState,
  fetchNewestNotifications,
  fetchAdditionalNotifications,
  updateNotificationState,
} from '../redux/notifications/reducer';
import { ActionTypes as NotificationsActionTypes } from '../redux/notifications/actions';
import Loading from './Loading';
import { headerStyle, headerTitleStyle, headerTintColor } from './TopHeader';
import {NewMatchNotification, Notification} from '../models/notification';
import Colors from '../services/colors';
import { ViewStyle } from 'react-native';
import { TextStyle } from 'react-native';
import { Linking } from 'expo';
import ProfileAvatar from "../components/ProfileAvatar";

interface DispatchActions {
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  fetchAdditionalNotifications:
    ActionCreator<ThunkAction<Promise<NotificationsActionTypes>, NotificationsState, void>>;
  fetchNewestNotifications:
    ActionCreator<ThunkAction<Promise<NotificationsActionTypes>, NotificationsState, void>>;
  updateNotificationState:
    ActionCreator<ThunkAction<Promise<NotificationsActionTypes>, NotificationsState, void>>;
}

interface Props extends NotificationsState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

interface State {
  refreshing: boolean;
}

type Row = "HEADER" | Notification | "LOAD_MORE" | "NO_MORE";

const ICON_SIZE = 48;

class NotificationView extends Component<Props, State> {
  NOTIFICATIONS_VIEW_IDENTIFIER = "NotificationView";

  static navigationOptions = {
    headerTitle: 'Notifications',
    headerStyle,
    headerTitleStyle,
    headerTintColor
  }

  constructor(props: Props) {
    super(props);

    this.state = { refreshing: false };

    this.load = this.load.bind(this);
    this.renderBody = this.renderBody.bind(this);
    this.renderRow = this.renderRow.bind(this);
    this.renderNotification = this.renderNotification.bind(this);
    this.onRefresh = this.onRefresh.bind(this);
    this.navigateHome = this.navigateHome.bind(this);
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

  private async navigateHome() {
    await this.props.navigation.dispatch(NavigationActions.reset({
      index: 0,
      actions: [NavigationActions.navigate({ routeName: 'Tabbed' })]
    }));
    await this.props.navigation.navigate('Home');
  }

  private renderNotification(notification: Notification) {

    let notifText: ReactNode = null;
    let icon = <MaterialIcons size={ICON_SIZE} name='face'/>;
    let onPressAction: () => void = this.navigateHome;
    const {
        notificationId,
        state,
        timestamp,
        type,
        thumbnail,
        message,
        link,
    } = notification;

    // update action to use deeplink
    if (link !== null && link !== undefined) {
      let { path, queryParams } = Linking.parse(link);
      onPressAction = this.props.navigation.navigate.bind(this, path, queryParams);
    }

    notifText = <Text>{message}</Text>;
    if (thumbnail) {
      icon = <Image style={styles.notifImageStyle} source={{uri: thumbnail}}/>;
    } else {
      icon = <MaterialIcons size={ICON_SIZE} name='message'/>;
    }

    // TODO: enforce stricter typing with metadata
    switch (notification.type) {
      case 'NEW_CREDENTIAL_MATCH':
        const { data } = notification;
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
        icon = <MaterialIcons size={ICON_SIZE} name='account-circle'/>;
        break;
      case 'CONNECTION_REQUESTED':
      case 'CONNECTION_ACCEPTED':
        const { data: {connUserId } } = notification;
        if (connUserId) {
          icon = <ProfileAvatar userId={connUserId} medium/>
        } else {
          icon = <MaterialIcons size={ICON_SIZE} name='account-circle'/>;
        }
        break;
      case 'NEW_MATCH':
        const { data: {matchUserId} } = notification;
        if (matchUserId) {
          icon = <ProfileAvatar userId={matchUserId} medium/>
        } else {
          icon = <MaterialIcons size={ICON_SIZE} name='account-circle'/>;
        }
        break;
      case 'ADHOC_NOTIFICATION':
        break;
      default:
        break;
    }

    const containerStyle = [styles.notifContainer];
    if (state === 'UNREAD') {
      containerStyle.push(styles.notifContainerUnread);
    }

    let onPress = (async () => {
      await onPressAction();
      if (state === 'UNREAD') {
        await this.props.updateNotificationState(Immutable.List([notificationId]), 'READ');
      }
    }).bind(this);

    return (
      <TouchableOpacity key={notification.notificationId} style={containerStyle} onPress={onPress}>
        <View style={styles.leftContainer}>
          {icon}
        </View>
        <View style={styles.rightContainer}>
          {notifText}
          <Text style={styles.ago}>{Moment(timestamp).fromNow()}</Text>
        </View>
      </TouchableOpacity>
    );
  }

  private renderRow(r: Row) {
    switch (r) {
      case "HEADER":
        const getUnread = () => this.props.notifications.filter(n => n.state === 'UNREAD');
        const hasUnread = !getUnread().isEmpty();
        const onPress = async () => {
          const unreadNotifIds = getUnread().map(n => n.notificationId).toList();
          try {
            await this.props.updateNotificationState(unreadNotifIds, 'READ');
          } catch (e) {
            await this.props.errorToast(e.errorMsg);
          }
        }
        const button = hasUnread ? (
          <Button
            buttonStyle={styles.markAllReadButton}
            onPress={onPress}
            title="Mark all as read"
          />
        ) : null;
        return <View style={styles.header}>{ button }</View>;
      case "LOAD_MORE":
        return (
          <View style={styles.loadingMoreContainer}>
            <Text style={styles.loadingMore}>Loading more...</Text>
          </View>
        );
      case "NO_MORE":
        return (
          <View style={styles.loadingMoreContainer}>
            <Text style={styles.loadingMore}>No more notifications</Text>
          </View>
        );
      default:
        return this.renderNotification(r);
    }
  }

  renderBody() {
    const ds = new ListView.DataSource({rowHasChanged: (n1: Notification, n2: Notification) => {
      return n1.notificationId !== n2.notificationId;
    }});
    const footer: Row = this.props.hasAll ? 'NO_MORE' : 'LOAD_MORE';
    const data = ds.cloneWithRows(
      ['HEADER'].concat(this.props.notifications.toJS().concat([footer])),
    );
    const onEndReached = (async () => {
      if (!this.props.hasAll && this.props.notifications.size > 0) {
        const oldestNotifId = this.props.notifications.last().notificationId;
        try {
          await this.props.fetchAdditionalNotifications(oldestNotifId);
        } catch (e) {
          await this.props.errorToast(e.errorMsg);
        }
      }
    }).bind(this);

    return (
      <ListView
        refreshControl={
          <RefreshControl
            refreshing={this.state.refreshing}
            onRefresh={this.onRefresh}
          /> as React.ReactElement<RefreshControlProps>
        }
        dataSource={data}
        renderRow={this.renderRow}
        onEndReached={onEndReached}
        onEndReachedThreshold={40}
      />
    );
  }

  render() {
    const {
      state,
      errorMsg,
      errorType,
    } = this.props.fetchState;
    // If `this.state.refreshing` is true, it means that we are reloading data using the pull
    // down, which means that we want to still display the ListView.
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
  {
    fetchAdditionalNotifications,
    fetchNewestNotifications,
    updateNotificationState,
    errorToast,
  })(NotificationView);

const styles = StyleSheet.create({
  ago: {
    paddingTop: 5,
    color: Colors.DARK_GRAY,
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
  notifImageStyle: {
    width: ICON_SIZE,
    height: ICON_SIZE,
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
  loadingMoreContainer: {
    height: 80,
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingMore: {
    color: 'gray',
  },
  header: {
    justifyContent: 'center',
    alignItems: 'center',
    borderBottomWidth: 0.5,
    borderColor: Colors.HIVE_SUBDUED,
  },
  markAllReadButton: {
    width: 200,
    margin: 10,
  },
});
