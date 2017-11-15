import React, { Component } from 'react';
import { AppRegistry, Text, View, FlatList, StyleSheet } from 'react-native';

export default class DefaultListElement extends Component {
  render() {
    return(
      <View>
        <Text style={styles.item}>{this.props.data.getUserName()}</Text>
      </View>
    );
  }
}

const styles = StyleSheet.create({
  item: {
    padding:10,
    fontSize: 30,
    height: 58
  }
});
