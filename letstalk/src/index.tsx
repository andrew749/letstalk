import React, { Component } from 'react';
import { Button, StyleSheet, Platform, Text, View } from 'react-native';
import { Provider } from 'react-redux';
import { combineReducers, compose, createStore, applyMiddleware } from 'redux';
import { MaterialIcons } from '@expo/vector-icons';
import { Notifications } from 'expo';
import createLogger from 'redux-logger';
import thunk from 'redux-thunk';
import { StackNavigator, TabNavigator } from 'react-navigation';
import NotificationComponent from 'react-native-in-app-notification';
import Sentry from 'sentry-expo';
import { Toast } from 'react-native-redux-toast';
import { YellowBox } from 'react-native'

import appReducer from './redux';
import auth from './services/auth';
import HomeView from './views/HomeView';
import LoginView from './views/LoginView';
import ProfileView from './views/ProfileView';
import ProfileEditView from './views/ProfileEditView';
import SignupView from './views/SignupView';
import OnboardingView from './views/OnboardingView';
import RequestToMatchView from './views/RequestToMatchView';
import NotificationService, { Notification } from './services/notification-service';

import Colors from './services/colors';

YellowBox.ignoreWarnings(['Warning: isMounted(...) is deprecated']);

Sentry.config('https://444853e4fac84788bbc1247f5c62c82b@sentry.io/1161982').install();

interface TabBarIcon {
  tintColor: (string | null),
  focused: boolean,
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: Colors.HIVE_MAIN_BG,
    alignItems: 'center',
    justifyContent: 'center'
  },
  toastMessageStyle: {
    color: 'white',
  },
});

const createTabView = () => TabNavigator({
  'Home': {
    screen: HomeView,
  },
  'Requests': {
    screen: RequestToMatchView,
  },
  'Profile': {
    screen: ProfileView,
  },
}, {
  tabBarPosition: 'bottom',
  navigationOptions: ({ navigation }) => ({
    tabBarIcon: ({ focused, tintColor }) => {
      const { routeName } = navigation.state;
      let iconName;
      if (routeName === 'Home') {
        iconName = 'home';
      } else if (routeName === 'Profile') {
        iconName = 'person';
      } else if (routeName === 'Requests') {
        iconName = 'supervisor-account';
      }

      return <MaterialIcons name={iconName} size={24} color={tintColor} />;
    },
  }),
  tabBarOptions: {
    showLabel: Platform.IOS ? true: false,
    showIcon: true,
    activeTintColor: Colors.HIVE_MAIN_BG,
    inactiveTintColor: 'gray',
    style: {
      backgroundColor: 'white',
    },
  },
});

const createAppNavigation = (loggedIn: boolean) => StackNavigator({
  Login: {
    screen: LoginView,
  },
  Signup: {
    screen: SignupView,
  },
  Tabbed: {
    screen: createTabView(),
  },
  ProfileEdit: {
    screen: ProfileEditView,
  },
  Onboarding: {
    screen: OnboardingView,
  },
}, {
  initialRouteName: loggedIn ? "Tabbed" : "Login",
});

const store = createStore(appReducer, applyMiddleware(thunk));

interface AppState {
  loggedIn: null | boolean;
}

type Props = {};

class App extends React.Component<Props, AppState> {
  private notificationService: NotificationService;

  constructor(props: Props) {
    super(props);
    this.state = {
      loggedIn: null,
    };

    this.handleNotification = this.handleNotification.bind(this);
    this.notificationService = null;
  }

  async handleNotification(notification: any) {
    await this.notificationService.handleNotification(notification as Notification);
  }

  async componentWillMount() {
    const sessionToken = await auth.getSessionToken();
    this.setState({ loggedIn: sessionToken !== null });
    Notifications.addListener(this.handleNotification);
  }

  render() {
    const { loggedIn } = this.state;

    const AppNavigation = createAppNavigation(loggedIn);
    return (
      <Provider store={store}>
        <View style={{ flex: 1 }}>
          <AppNavigation />
          <NotificationComponent ref={(ref: any) => {
            this.notificationService = new NotificationService(ref, store);
          }} />
          <Toast messageStyle={styles.toastMessageStyle} />
        </View>
      </Provider>
    );
  }
}

export default App;
