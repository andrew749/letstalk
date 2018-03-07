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

import appReducer from './redux';
import auth from './services/auth';
import AchievementsView from './views/AchievementsView';
import EventsView from './views/EventsView';
import HomeView from './views/HomeView';
import LoginView from './views/LoginView';
import ProfileView from './views/ProfileView';
import SignupView from './views/SignupView';
import OnboardingView from './views/OnboardingView';

interface TabBarIcon {
  tintColor: (string | null),
  focused: boolean,
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f1f',
    alignItems: 'center',
    justifyContent: 'center'
  }
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
  Onboarding: {
    screen: OnboardingView,
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
      title: 'Got a notification',
      message: '' + notification.data.yo,
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
        </View>
      </Provider>
    );
  }
}

export default App;
