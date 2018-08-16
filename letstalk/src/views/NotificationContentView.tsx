import {WebView} from 'react-native';
import React from 'react';
import {Notification} from '../models/notification';
import { NavigationScreenProp, NavigationStackAction } from 'react-navigation';
import { headerStyle } from './TopHeader';
import { View, Text } from 'react-native';
import {BASE_URL, NOTIFICATION_PAGE_ROUTE} from '../services/constants';
import auth from '../services/auth';

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  notification: Notification;
}

export const NotificationContentViewRoute = "NotificationContent";
interface State {
  sessionId: string;
}

class NotificationContentView extends React.Component<Props, State> {
    static navigationOptions = () => ({
      headerStyle,
    })
    private notificationId: number;
    private sessionId: string;

    constructor(props: Props) {
      super(props);
      //@ts-ignore
      this.notificationId = (this.props.notification && this.props.notification.notificationId) || this.props.navigation.getParam("notificationId", undefined);
      this.state = {sessionId: undefined}
    }
    async componentDidMount() {
      const token = await auth.getSessionToken();
      this.setState({sessionId: token});
    }

    private getNotificationPage(notificationId: number): string {
      return `${BASE_URL}${NOTIFICATION_PAGE_ROUTE}?notificationId=${notificationId}`;
    }

    render() {
      const url = this.getNotificationPage(this.notificationId);
      if (this.state.sessionId !== undefined && this.state.sessionId !== null) {
        return <WebView source={{
          uri: url,
          headers: {"sessionId": this.state.sessionId}
        }} />;
      }
      return <View/>;
    }
}

export default NotificationContentView;
