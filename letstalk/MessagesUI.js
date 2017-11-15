import React, { Component } from 'react';
import { AppRegistry, Text, View, FlatList, StyleSheet } from 'react-native';
import DefaultListElement from './ListElement';
import MessageData from './messagedata';

export default class MessagesUI extends Component {
  render() {
    return(
      <View style={styles.container}>
        <FlatList
          data={[
            new MessageData(1, "Andrew")
          ]}
          keyExtractor={ (item) => item.getId() }
          renderItem={(
            {item}) =>
              <DefaultListElement
                data={item}
                style={styles.item}>
              </DefaultListElement>
          }
        />
      </View>
    );
  }
}

const styles = StyleSheet.create({
  container: {
   flex: 1,
   paddingTop: 22
  },
  item: {
    padding: 10,
    fontSize: 18,
    height: 44,
  },
});
