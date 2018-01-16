import React, { Component } from 'react';
import { ScrollView, AppRegistry, Text, TextInput, View, FlatList, StyleSheet } from 'react-native';
import { connect } from 'react-redux';
import { Dispatch } from 'redux'

import MessageData from '../models/message-data';
import { RootState } from '../redux';
import { State as ThreadState } from '../redux/thread/reducer';

interface Props extends ThreadState {
  navigation: any;
};

class MElementData {
    message: string;
    ours: boolean;
    constructor(message: string, ours: boolean) {
        this.message = message;
        this.ours = ours;
    }

}

class MessageView extends Component<Props> {
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

function mapStateToProps(state: RootState): ThreadState {
  return state.thread;
}

function mapDispatchToProps(dispatch: Dispatch<ThreadState>): any {
  return {};
}

export default connect(mapStateToProps, mapDispatchToProps)(MessageView);

// TODO: Move style into some stylesheet
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
