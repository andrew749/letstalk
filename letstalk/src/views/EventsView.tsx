import React, { Component } from 'react';
import {
  Text,
  View,
} from 'react-native';

interface Props {}

export default class EventsView extends Component<Props> {
  static navigationOptions = {
    headerTitle: 'Events',
  }

  render() {
    return (
      <View><Text>EventsView</Text></View>
    );
  }
}
