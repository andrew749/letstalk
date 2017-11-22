import React, { Component } from 'react';
import { AppRegistry, Text, View, FlatList, StyleSheet } from 'react-native';
import { connect } from 'react-redux';

import DefaultListElement from './DefaultListElement';
import MessageData from '../models/message-data';

class MessagesList extends Component {
  render() {
    return(
      <View style={styles.container}>
        <FlatList
          data={this.props.threads}
          keyExtractor={ (item) => item.getId() }

          renderItem={(
            {item}) =>
              <DefaultListElement
                onPress={() => { this.props.navigation.navigate('MessageThread', {'name': item.getUserName()}); } }
                data={ item }
                style={ styles.item }>
              </DefaultListElement>
          }
        />
      </View>
    );
  }
}

export default connect(({conversationsReducer}) => conversationsReducer, {})(MessagesList);

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
