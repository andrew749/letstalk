import React, { Component } from 'react';
import { Button, StyleSheet, Text, View } from 'react-native';
import { Provider } from 'react-redux';
import { combineReducers, compose, createStore, applyMiddleware } from 'redux';
import { MaterialIcons } from '@expo/vector-icons';
import { Notifications } from 'expo';
import createLogger from 'redux-logger';
import thunk from 'redux-thunk';
import { StackNavigator, TabNavigator } from 'react-navigation';
import Notification from 'react-native-in-app-notification';
import Sentry from 'sentry-expo';
import { Toast } from 'react-native-redux-toast';

import appReducer from './redux';
import auth from './services/auth';
import AchievementsView from './views/AchievementsView';
import CredentialEditView from './views/CredentialEditView';
import EventsView from './views/EventsView';
import HomeView from './views/HomeView';
import LoginView from './views/LoginView';
import ProfileView from './views/ProfileView';
import SignupView from './views/SignupView';
import OnboardingView from './views/OnboardingView';
import RequestToMatchView from './views/RequestToMatchView';

import Colors from './services/colors';
import QrCodeView from "./views/QrCodeView";

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

const createAppNavigation = (loggedIn: boolean) => StackNavigator({
  Login: {
    screen: LoginView,
  },
  Signup: {
    screen: SignupView,
  },
  Home: {
    screen: HomeView,
  },
  Profile: {
    screen: ProfileView,
  },
  RequestToMatch: {
    screen: RequestToMatchView,
  },
  CredentialEdit: {
    screen: CredentialEditView,
  },
  Onboarding: {
    screen: OnboardingView,
  },
  QrCode: {
    screen: QrCodeView,
  },
}, {
  initialRouteName: loggedIn ? "Home" : "Login",
});

const store = createStore(appReducer, applyMiddleware(thunk));

interface AppState {
  loggedIn: null | boolean;
}

type Props = {};

class App extends React.Component<Props, AppState> {
  private notification: any;

  constructor(props: Props) {
    super(props);
    this.state = {
      loggedIn: null,
    };

    this.handleNotification = this.handleNotification.bind(this);
  }

  handleNotification(notification: any) {
    this.notification.show({
      title: notification.data.title,
      message: notification.data.message,
    });
  }

  async componentWillMount() {
    const sessionToken = await auth.getSessionToken();
    this.setState({ loggedIn: sessionToken !== null });
    Notifications.addListener(this.handleNotification);
  }

  render() {
    const { loggedIn } = this.state;
    if (loggedIn === null) return <View><Text>Splash</Text></View>;

    const AppNavigation = createAppNavigation(loggedIn);
    return (
      <Provider store={store}>
        <View style={{ flex: 1 }}>
          <AppNavigation />
          <Notification ref={(ref: any) => { this.notification = ref; }} />
          <Toast messageStyle={styles.toastMessageStyle} />
        </View>
      </Provider>
    );
  }
}

export default App;
