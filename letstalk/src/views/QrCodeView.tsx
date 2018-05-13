import React, { Component } from 'react';
import {
  Text,
  View,
} from 'react-native';

interface Props {
  userId: number,
}

export default class QrCodeView extends Component<Props> {
  static navigationOptions = {
    headerTitle: 'QrCode',
  }

  render() {
    return (
      <View><Text>QrCodeView for {this.props.userId}</Text></View>
    );
  }
}
