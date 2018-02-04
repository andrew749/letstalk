import React, { Component } from 'react';
import { Button, StyleSheet, Text, View } from 'react-native';
import { Provider } from 'react-redux';
import { combineReducers, compose, createStore, applyMiddleware } from 'redux';
import MaterialIcons from 'react-native-vector-icons/MaterialIcons';
import { Icon } from 'react-native-elements';
import createLogger from 'redux-logger';
import thunk from 'redux-thunk';

import appReducer from './redux';
import auth from './services/auth';
import LoginView from './views/LoginView';
import ProfileView from './views/ProfileView';

import { StackNavigator, TabNavigator } from 'react-navigation';

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
    screen: ProfileView,
    navigationOptions: {
      tabBarLabel:"Home",
    },
  },
  Events: {
    screen: ProfileView,
    navigationOptions: {
      tabBarLabel:"Events",
    },
  },
  Achievements: {
    screen: ProfileView,
    navigationOptions: {
      tabBarLabel:"Achievements",
      tabBarIcon: <Icon name='rowing' />,
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
    await MaterialIcons.loadFont();
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

