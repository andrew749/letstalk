import React, { Component } from 'react';
import { Button, StyleSheet, Text, View } from 'react-native';
import { Provider } from 'react-redux';
import { combineReducers, compose, createStore, applyMiddleware } from 'redux';
import createLogger from 'redux-logger';
import thunk from 'redux-thunk';

import appReducer from './redux';

import MessagesList from './components/MessagesList';
import MessageView from './components/MessageView';
import LoginView from './views/LoginView';
import ProfileView from './views/ProfileView';

import { StackNavigator, TabNavigator } from 'react-navigation';

// TODO: Investigate typing with navigation
const ConversationsScene = ({ navigation }: any) => (
  <MessagesList navigation={ navigation }/>
);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f1f',
    alignItems: 'center',
    justifyContent: 'center'
  }
});

const MainPage = TabNavigator({
  TabItem1: {
    screen: ProfileView,
    navigationOptions: {
      tabBarLabel:"Tab 1",
    },
  },
  TabItem2: {
    screen: ProfileView,
    navigationOptions: {
      tabBarLabel:"Tab 1",
    },
  },
});

const AppNavigation = StackNavigator({
  Login: {
    screen: LoginView,
  },
  Profile: {
    screen: MainPage,
    navigationOptions: {
      header: null,
    },
  },
  Home: {
    screen: ConversationsScene,
    navigationOptions: {
      headerTitle: 'Conversations'
    }
  },
  MessageThread: {screen: MessageView}
});

const store = createStore(appReducer, applyMiddleware(thunk));

class App extends React.Component {
  render() {
    return (
      <Provider store={store}>
        <AppNavigation />
      </Provider>
    );
  }
}

export default App;

