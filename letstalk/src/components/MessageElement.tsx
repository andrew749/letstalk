import React, { Component } from 'react';
import { AppRegistry, Text, View, StyleSheet } from 'react-native';

export default class MessageElement extends Component {
  render() {
    return(
      <View>
        <Text>Message Text</Text>
        <Text>Sender Image</Text>
      </View>
    );
  }
}

const styles = StyleSheet.create({
  item: {
   flex: 1,
   flexDirection: 'row',
   paddingTop: 22
  }
});
