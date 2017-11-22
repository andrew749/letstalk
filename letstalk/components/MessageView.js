import React, { Component } from 'react';
import { ScrollView, AppRegistry, Text, TextInput, View, FlatList, StyleSheet } from 'react-native';
import { connect } from 'react-redux';

import MessageData from '../models/message-data';
import { fetchMessages } from '../state/thread';

class MessageView extends Component {
  componentDidMount() {
    this.props.fetchMessages(this.props.navigation.state.params.name);
  }

  render() {
    const { params } = this.props.navigation.state;
    const placeholderText = `Send a message to ${ params.name }`;

    return(
      <View style={ styles.container }>
        <ScrollView></ScrollView>
        <TextInput value={placeholderText} />
      </View>
    );
  }
}

export default connect(({threadReducer}) => threadReducer, { fetchMessages })(MessageView);

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
  textInput: {
    padding:20
  }
});
