import React, { Component } from 'react';
import { AppRegistry, Text, View, FlatList, StyleSheet } from 'react-native';
import { connect } from 'react-redux';

import DefaultListElement from './DefaultListElement';
import MessageData from '../models/message-data';
import { RootState } from '../redux';
import { State as ConversationState } from '../redux/conversation/reducer';

interface Props extends ConversationState {
  navigation: any;
}

function mapStateToProps (state: RootState): ConversationState {
  return state.conversation;
};

class MessagesList extends Component<Props> {
  render() {
    return(
      <View style={styles.container}>
        <FlatList
          data={this.props.threads.toArray()}
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

export default connect(mapStateToProps)(MessagesList);

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
