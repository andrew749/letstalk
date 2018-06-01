import React, { Component } from 'react';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import { ThunkAction } from 'redux-thunk';

import { RootState } from '../redux';
import { errorToast } from '../redux/toast';
import { combineFetchStates } from '../redux/actions';
import {
  State as CredentialRequestsState,
  addCredentialRequest,
  fetchCredentialRequests,
} from '../redux/credential-requests/reducer';
import {
  State as CredentialOptionsState,
  fetchCredentialOptions,
} from '../redux/credential-options/reducer';
import { ActionTypes as CredentialsActionTypes } from '../redux/credentials/actions';
import { ActionTypes as CredentialRequestsActionTypes } from '../redux/credential-requests/actions';
import { ActionTypes as CredentialOptionsActionTypes } from '../redux/credential-options/actions';


import {
  Dimensions,
  StyleSheet,
  View,
  Platform,
} from 'react-native';
import {
  NavigationScreenProp,
  NavigationScreenDetails,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import {
  ActionButton,
  Card,
  FilterableElement,
  FilterListModal,
  TextHeader,
} from '../components';
import { Credential } from '../models/credential';

import Colors from '../services/colors';

const SCREEN_WIDTH = Dimensions.get('window').width

interface DispatchActions {
  addCredentialRequest: ActionCreator<
    ThunkAction<Promise<CredentialRequestsActionTypes>, CredentialRequestsState, void>>;
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  fetchCredentialRequests: ActionCreator<
    ThunkAction<Promise<CredentialRequestsActionTypes>, CredentialRequestsState, void>>;
  fetchCredentialOptions: ActionCreator<
    ThunkAction<Promise<CredentialOptionsActionTypes>, CredentialOptionsState, void>>;
}

interface Props extends DispatchActions {
  credentialOptions: CredentialOptionsState;
  credentialRequests: CredentialRequestsState;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}


class Header extends Component<Props> {
  constructor(props: Props) {
      super(props);

      this.load = this.load.bind(this);
      this.onReqSelect = this.onReqSelect.bind(this);
  }

  async componentDidMount() {
      this.load();
  }

  private async load() {
      await Promise.all([
          this.props.fetchCredentialRequests(),
          this.props.fetchCredentialOptions(),
      ]);
  }

  private async onReqSelect(elem: FilterableElement): Promise<void> {
    try {
      await this.props.addCredentialRequest({ id: elem.id, name: elem.value });
    } catch (e) {
      await this.props.errorToast(e.message);
    }
  }

  render() {
      const { credentialOptions } = this.props.credentialOptions;
      return (
        <View>
          <View style={styles.header}/>
      
          <View style={styles.topContainer}>
            <FilterListModal
              data={credentialOptions.map(cred => { return { id: cred.id, value: cred.name }}).toList()}
              onSelect={this.onReqSelect}
              placeholder="Find someone who is a..."
            />
          </View>
      </View>
    );
  }
}

export default connect(
  ({ credentialRequests, credentialOptions }: RootState) => {
    return { credentialOptions, credentialRequests };
  }, {
    addCredentialRequest,
    fetchCredentialOptions,
    fetchCredentialRequests,
    errorToast,
  })(Header);

const styles = StyleSheet.create({
  header: {
    height: Platform.OS == "ios" ? 20 : 0,
    backgroundColor: Colors.HIVE_PRIMARY,
  },
  topContainer: {
    width: SCREEN_WIDTH,
    backgroundColor: '#FFC107',
  },
})