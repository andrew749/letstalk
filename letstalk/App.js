import React from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { Provider } from 'react-redux';

import store from './store';
import MessagesList from './components/MessagesList';
import MessageView from './components/MessageView';

import { StackNavigator } from 'react-navigation';

const ConversationsScene = ({ navigation }) => (
  <MessagesList navigation={navigation}/>
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

