import {WebView} from 'react-native';
import React from 'react';
import {Notification} from '../models/notification';
import { NavigationScreenProp, NavigationStackAction, NavigationScreenDetails } from 'react-navigation';
import { headerStyle, headerTitleStyle, headerTintColor } from './TopHeader';
import { View, Text } from 'react-native';
import {BASE_URL, NOTIFICATION_PAGE_ROUTE} from '../services/constants';
import auth from '../services/auth';
import {profileService} from '../services/profile-service';

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  notification: Notification;
}

export const NotificationContentViewRoute = "NotificationContent";
interface State {
  sessionId: string;
}

class NotificationContentView extends React.Component<Props, State> {
    static navigationOptions = ({navigation}: NavigationScreenDetails<void>) => {
      const {state} = navigation;
      return {
        title: `${state.params.title || "Notification"}`,
        headerStyle,
        headerTitleStyle,
        headerTintColor
      }
    };
    private notificationId: number;
    private sessionId: string;

    constructor(props: Props) {
      super(props);
      //@ts-ignore
      this.notificationId = (this.props.notification && this.props.notification.notificationId) || this.props.navigation.getParam("notificationId", undefined);
      //@ts-ignore
      const passedTitle = this.props.notification && this.props.notification.title;
      if (passedTitle) {
        //@ts-ignore
        this.props.navigation.setParams({"title": passedTitle});
      }
      this.state = {sessionId: undefined}
    }
    async componentDidMount() {
      const token = await auth.getSessionToken();
      this.setState({sessionId: token});
      await this.load();
    }

    async load() {
      console.log("Loading notification")
      const notification = await profileService.getNotificationForId(this.notificationId);
      console.log(notification);
      //@ts-ignore
      this.props.navigation.setParams({"title": notification.title});
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
