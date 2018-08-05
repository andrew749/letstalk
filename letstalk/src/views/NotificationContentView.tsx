import {WebView} from 'react-native';
import React from 'react';
import {Notification} from '../models/notification';
import { NavigationScreenProp, NavigationStackAction } from 'react-navigation';
import { headerStyle } from './TopHeader';
import { View, Text } from 'react-native';

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  notification: Notification;
}

interface NotificationContent {
  page: string;
}

interface State extends NotificationContent {
}

export const NotificationContentViewRoute = "NotificationContent";

class NotificationContentView extends React.Component<Props, State> {
    static navigationOptions = () => ({
      headerStyle,
    })

    constructor(props: Props) {
      super(props);
      //@ts-ignore
      const content = this.props.notification || this.props.navigation.getParam("notification", undefined);
    }

    async componentDidMount() {
      // load page
      // const content = await notificationAPIService.getNotificationContent();
      // this.setState({
      //   ...content,
      // })
    }

    render() {
      if (!this.state.page) {
          return  (
            <View>
              <Text>{"Loading"}</Text>
            </View>
          );
      }
      return <WebView html={this.state.page} />;
    }
}

export default NotificationContentView;
