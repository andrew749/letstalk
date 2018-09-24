import React, { Component } from 'react';
import { View, StyleSheet, Text } from 'react-native';
import {
  NavigationScreenDetails,
  NavigationScreenProp,
  NavigationStackAction,
  NavigationActions,
} from 'react-navigation';
import QRCode from "react-native-qrcode";
import { headerStyle, headerTitleStyle, headerTintColor } from './TopHeader';

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction & { secret: string }>;
}

export default class QrCodeView extends Component<Props> {
  static navigationOptions = {
    headerTitle: 'QR Code',
    headerStyle, 
    headerTitleStyle, 
    headerTintColor 
  }

  render() {
    const secret = this.props.navigation.getParam('secret', null);
    const qrCodeElement = !secret ? <Text>No QR code available...</Text> : (
      <QRCode
        value={secret}
        size={150}
        bgColor='black'
        fgColor='white'
      />
    );
    return (
      <View style={styles.container}>
        { qrCodeElement }
      </View>
    );
  }
}

const styles = StyleSheet.create({
  container: {
    width: '100%',
    height: '100%',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: 'white',
  },
});
