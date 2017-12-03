import React, { Component } from 'react';
import { ScrollView, AppRegistry, Text, TextInput, View, FlatList, StyleSheet } from 'react-native';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux'

import MessageData from '../models/message-data';
import { fetchMessages } from '../redux/thread/actions';

interface Props {
  fetchMessages: (userId: string) => any;
  messages: MessageData[];
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

class MElementData {
    message: string;
    ours: boolean;
    constructor(message: string, ours: boolean) {
        this.message = message;
        this.ours = ours;
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
        <FlatList
        data = { [ new MElementData("Test", true), new MElementData("Test theirs", false) ] }
        renderItem = {(obj) =>
            <View style = {[styles.message, obj.item.ours ? styles.ourMessage : styles.theirMessage ]} >
                <Text>{obj.item.message}</Text>
            </View>
        }
        />
        <TextInput
            style={ styles.textInput }
            value={ placeholderText }
            editable={ true }/>
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
  message: {
    padding:5,
    paddingRight: 10,
    paddingLeft: 10,
    borderRadius: 5,
    backgroundColor: "#01F"
  },
  theirMessage: {
    alignSelf: "flex-start",
    backgroundColor: "#00B0FF",
    marginLeft: 10
  },
  ourMessage: {
    alignSelf: "flex-end",
    backgroundColor: "#E0E0E0",
    marginRight: 10
  },
  item: {
    padding: 10,
    fontSize: 18,
    height: 44,
  },
  textInput: {
    marginBottom:10
  }
});
