import React, { Component } from 'react';
import {
  Text,
  View,
} from 'react-native';

interface Props {}

export default class AchievementsView extends Component<Props> {
  static navigationOptions = {
    headerTitle: 'Achievements',
  }

  render() {
    return (
      <View><Text>AchievementsView</Text></View>
    );
  }
}
