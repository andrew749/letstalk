import React, { Component } from 'react';
import {
  Button as ReactNativeButton,
  Text,
  View,
} from 'react-native';
import {
  NavigationScreenProp,
  NavigationScreenDetails,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

export default class HomeView extends Component<Props> {
  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: 'Home',
    headerRight: <ReactNativeButton title="Profile"
      onPress={() => navigation.navigate('Profile')} />,
  })

  render() {
    return (
      <View>
      <Text>HomeScreen</Text></View>
    );
  }
}
