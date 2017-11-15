import React from 'react';
import { StyleSheet, Text, View } from 'react-native';

import MessagesList from './messages-list';
import MessageData from './message-data';
import MessageView from './message-view';

import { StackNavigator } from 'react-navigation';


const data = [
  new MessageData(1, "Andrew"),
  new MessageData(2, "Adam")
];

export default class MainView extends React.Component {
  render() {
    return (
      <MessagesList data={data}/>
    );
  }
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f1f',
    alignItems: 'center',
    justifyContent: 'center',
  },
});

const App = StackNavigator({
  Home: { screen: MainView },
  MessageThread: {screen: MessageView}
});

