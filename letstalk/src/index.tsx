import React, { Component } from 'react';
import { Button, StyleSheet, Text, View } from 'react-native';
import { Provider } from 'react-redux';
import { combineReducers, compose, createStore, applyMiddleware } from 'redux';
import { MaterialIcons } from '@expo/vector-icons';
import createLogger from 'redux-logger';
import thunk from 'redux-thunk';

import appReducer from './redux';
import auth from './services/auth';
import AchievementsView from './views/AchievementsView';
import EventsView from './views/EventsView';
import HomeView from './views/HomeView';
import LoginView from './views/LoginView';
import ProfileView from './views/ProfileView';

import { StackNavigator, TabNavigator } from 'react-navigation';

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

const MainView = TabNavigator({
  Home: {
    screen: StackNavigator({
      Home: {
        screen: HomeView,
      },
      Profile: {
        screen: ProfileView,
      },
    }, {
      headerMode: 'float',
    }),
    navigationOptions: {
      tabBarLabel:"Home",
      tabBarIcon: ({ tintColor }: TabBarIcon) =>
        <MaterialIcons name="home" size={32} color={tintColor} />,
    },
  },
  Events: {
    screen: StackNavigator({
      Events: {
        screen: EventsView,
      },
    }),
    navigationOptions: {
      tabBarLabel:"Events",
      tabBarIcon: ({ tintColor }: TabBarIcon) =>
        <MaterialIcons name="event" size={32} color={tintColor} />,
    },
  },
  Achievements: {
    screen: StackNavigator({
      Achievements: {
        screen: AchievementsView,
      },
    }),
    navigationOptions: {
      tabBarLabel:"Achievements",
      tabBarIcon: ({ tintColor }: TabBarIcon) =>
        <MaterialIcons name="done" size={32} color={tintColor} />,
    },
  },
});

const createAppNavigation = (loggedIn: boolean) => StackNavigator({
  Login: {
    screen: LoginView,
  },
  Main: {
    screen: MainView,
    navigationOptions: {
      header: null,
    },
  },
}, {
  initialRouteName: loggedIn ? "Main" : "Login",
  headerMode: 'float',
});

const store = createStore(appReducer, applyMiddleware(thunk));

interface AppState {
  loggedIn: null | boolean;
}

class App extends React.Component<void, AppState> {
  constructor(props: void) {
    super(props);
    this.state = {
      loggedIn: null,
    };
  }

  async componentWillMount() {
    const sessionToken = await auth.getSessionToken();
    this.setState({ loggedIn: sessionToken !== null });
  }

  render() {
    const { loggedIn } = this.state;
    if (loggedIn === null) return <View><Text>Splash</Text></View>;

    const AppNavigation = createAppNavigation(loggedIn)
    return (
      <Provider store={store}>
        <AppNavigation />
      </Provider>
    );
  }
}

export default App;

