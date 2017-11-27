import React, { Component } from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { Provider } from 'react-redux';
import { combineReducers, compose, createStore, applyMiddleware } from 'redux';
import createLogger from 'redux-logger';
import thunk from 'redux-thunk';

import rootReducer from './redux';

import MessagesList from './components/MessagesList';
import MessageView from './components/MessageView';

import { StackNavigator } from 'react-navigation';

const ConversationsScene = ({ navigation, threads }: any) => (
  <MessagesList threads= { threads} navigation={ navigation }/>
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
  Home: {
    screen: ConversationsScene,
    navigationOptions: {
      headerTitle: 'Conversations'
    }
  },
  MessageThread: {screen: MessageView}
});

const store = createStore(rootReducer, applyMiddleware(thunk));

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

