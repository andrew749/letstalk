import React, { Component } from 'react';
import { AppRegistry, Text, View, StyleSheet } from 'react-native';

export default class MessageElement extends Component {
  render() {
    <View>
      <Text>Message Text</Text>
      <Text>Sender Image</Text>
    </View>
  }
}

const styles = StyleSheet.create({
  item: {
   flex: 1,
   flex-direction: row;
   paddingTop: 22
  }
});
