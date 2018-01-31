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

import { StackNavigator } from 'react-navigation';

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

const AppNavigation = StackNavigator({
  Profile: {
    screen: ProfileView,
  },
  Login: {
    screen: LoginView,
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

