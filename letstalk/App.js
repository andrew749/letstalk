import React from 'react';
import { StyleSheet, Text, View } from 'react-native';

import MessagesList from './messages-list';
import MessageData from './message-data';
import MessageView from './message-view';

import { StackNavigator } from 'react-navigation';


var data = [
  new MessageData(1, "Andrew"),
  new MessageData(2, "Adam")
];

const ConversationsScene = ({ navigation }) => (
  <MessagesList data={data} navigation={navigation}/>
);

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f1f',
    alignItems: 'center',
    justifyContent: 'center'
  }
});

const App = StackNavigator({
  Home: {
    screen: ConversationsScene,
    navigationOptions: {
      headerTitle: 'Conversations'
    }
  },
  MessageThread: {screen: MessageView}
});

export default App;

