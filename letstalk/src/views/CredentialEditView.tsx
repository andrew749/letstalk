import React, { Component } from 'react';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import {
  ActivityIndicator,
  Button as ReactNativeButton,
  Linking,
  ScrollView,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from 'react-native';
import {
  NavigationScreenProp,
  NavigationScreenDetails,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import { MaterialIcons } from '@expo/vector-icons';

import { RootState } from '../redux';
import { State as CredentialsState, fetchCredentials } from '../redux/credentials/reducer';
import { ActionTypes } from '../redux/credentials/actions';
import { ActionButton, Card, Header, Loading } from '../components';

interface DispatchActions {
  fetchCredentials: ActionCreator<ThunkAction<Promise<ActionTypes>, CredentialsState, void>>;
}

interface Props extends CredentialsState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class CredentialEditView extends Component<Props> {
  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: 'Edit Credentials',
  })

  constructor(props: Props) {
    super(props);

    this.load = this.load.bind(this);
    this.renderHome = this.renderHome.bind(this);
  }

  async componentDidMount() {
    this.load();
  }

  private async load() {
    await this.props.fetchCredentials();
  }

  private renderHome() {
    return <Text>yo</Text>;
  }

  render() {
    const {
      state,
      errorMsg,
    } = this.props.fetchState;
    return (
      <Loading
        state={state}
        load={this.load}
        errorMsg={errorMsg}
        renderBody={this.renderHome}
      />
    );
  }
}

export default connect(
  ({ credentials }: RootState) => credentials,
  { fetchCredentials }
)(CredentialEditView);
