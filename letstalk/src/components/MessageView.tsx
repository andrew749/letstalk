import React, { Component } from 'react';
import { ScrollView, AppRegistry, Text, TextInput, View, FlatList, StyleSheet } from 'react-native';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux'

import MessageData from '../models/message-data';
import { fetchMessages } from '../redux/thread/actions';

interface Props {
  fetchMessages: (userId: string) => any;
  messages: any[];
  navigation: any;
}

function mapStateToProps(state: any): any {
  return { 
    messages: state.threadReducer.messages
  }
}

function mapDispatchToProps(dispatch: any): any {
  console.log(dispatch);
  return {
    fetchMessages: (userId: string) => {
      dispatch(fetchMessages(userId));
    }
  }
}

class MessageView extends Component<Props> {
  componentDidMount() {
    this.props.fetchMessages(this.props.navigation.state.params.name);
  }

  render() {
    const { params } = this.props.navigation.state;
    const placeholderText = `Send a message to ${ params.name }`;

    return(
      <View style={ styles.container }>
        <ScrollView></ScrollView>
        <TextInput value={ placeholderText } />
      </View>
    );
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(MessageView);

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
