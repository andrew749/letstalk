import React, { Component } from 'react';
import { AppRegistry, Text, View, FlatList, StyleSheet, TouchableOpacity } from 'react-native';

export default class DefaultListElement extends Component {
  constructor() {
    super();
  }

  render() {
    return(
      <TouchableOpacity onPress={ this.props.onPress }>
        <View style={ styles.listItem } >
          <Text style={styles.item}>{ this.props.data.getUserName() }</Text>
        </View>
      </TouchableOpacity>
    );
  }
}

const styles = StyleSheet.create({
  listItem: {
  },
  item: {
    padding:10,
    fontSize: 30,
    height: 58
  }
});
