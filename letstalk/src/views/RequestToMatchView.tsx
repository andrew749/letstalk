import React, { Component } from 'react';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import {
  ActivityIndicator,
  Button as ReactNativeButton,
  Linking,
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
import Immutable from 'immutable';

import { RootState } from '../redux';
import {
  State as CredentialOptionsState,
  fetchCredentialOptions,
} from '../redux/credential-options/reducer';
import { ActionTypes } from '../redux/credential-options/actions';
import { ActionButton, Card, FilterListModal, Header, Loading } from '../components';
import { CredentialFilterableElement } from '../models/credential';

interface DispatchActions {
  fetchCredentialOptions: ActionCreator<
    ThunkAction<Promise<ActionTypes>, CredentialOptionsState, void>>;
}

interface Props extends CredentialOptionsState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class RequestToMatchView extends Component<Props> {
  static navigationOptions = () => ({
    headerTitle: 'Request To Match',
  })

  constructor(props: Props) {
    super(props);

    this.load = this.load.bind(this);
    this.renderSearch = this.renderSearch.bind(this);
  }

  async componentDidMount() {
    this.load();
  }

  private async load() {
    await this.props.fetchCredentialOptions();
  }

  private async onSelect(elem: CredentialFilterableElement): Promise<void> {
    console.log(elem);
  }

  private renderSearch() {
    return (
      <FilterListModal
        data={this.props.credentialElements}
        onSelect={this.onSelect}
        placeholder="Find someone who is a..."
      />
    );
  }

  render() {
    const {
      state,
      errorMsg,
    } = this.props.fetchState;
    return (
      <Loading
        state={state}
        errorMsg={errorMsg}
        load={this.load}
        renderBody={this.renderSearch}
      />
    );
  }
}

export default connect(
  ({ credentialOptions }: RootState) => credentialOptions, { fetchCredentialOptions }
)(RequestToMatchView);

const styles = StyleSheet.create({
  centeredContainer: {
    flex: 1,
    justifyContent: 'center'
  },
  headline: {
    fontWeight: 'bold',
    fontSize: 18,
    textAlign: 'center',
    marginBottom: 10,
  },
})
