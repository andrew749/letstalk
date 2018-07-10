import React, { Component } from 'react';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import { ThunkAction } from 'redux-thunk';

import { RootState } from '../redux';
import {
  State as SearchBarState,
  updateFocus,
} from '../redux/search-bar/reducer';
import {
  FilterableElement,
  FilterListModal,
} from '../components';
import { errorToast } from '../redux/toast';
import {
  State as CredentialsState,
  addCredential,
} from '../redux/credentials/reducer';
import {
  State as CredentialRequestsState,
  addCredentialRequest,
} from '../redux/credential-requests/reducer';
import {
  State as CredentialOptionsState,
} from '../redux/credential-options/reducer';
import {
  ActionTypes as SearchBarActionTypes,
  SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS,
  SEARCH_LIST_TYPE_CREDENTIALS,
} from '../redux/search-bar/actions';
import { ActionTypes as CredentialsActionTypes } from '../redux/credentials/actions';
import { ActionTypes as CredentialRequestsActionTypes } from '../redux/credential-requests/actions';
import { ActionTypes as CredentialOptionsActionTypes } from '../redux/credential-options/actions';
import Colors from '../services/colors';
import { AnalyticsActions, logAnalyticsThenExecuteAsync } from '../services/analytics';
import { Alert } from 'react-native';

interface DispatchActions {
  addCredentialRequest: ActionCreator<
    ThunkAction<Promise<CredentialRequestsActionTypes>, CredentialRequestsState, void>>;
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  addCredential: ActionCreator<
    ThunkAction<Promise<CredentialsActionTypes>, CredentialsState, void>>;
  updateFocus: ActionCreator<
    ThunkAction<Promise<SearchBarActionTypes>, SearchBarState, void>>;
}

interface Props extends DispatchActions {
  searchBar: SearchBarState;
  credentialOptions: CredentialOptionsState;
  onSelectSuccess?(): void;
}

class AllFilterableModals extends Component<Props> {

  constructor(props: Props) {
    super(props);

    this.onReqSelect = this.onReqSelect.bind(this);
    this.onCredSelect = this.onCredSelect.bind(this);
    this.onRawCredSelect = this.onRawCredSelect.bind(this);
    this.onRawCredReqSelect = this.onRawCredReqSelect.bind(this);
  }

  private async blurSearchBar() {
    await this.props.updateFocus(false);
  }

  private selectSuccess() {
    if (!!this.props.onSelectSuccess) this.props.onSelectSuccess();
  }

  private async onReqSelect(elem: FilterableElement): Promise<void> {
    try {
      logAnalyticsThenExecuteAsync(
        "Credentials",
        AnalyticsActions.ADD,
        "CredentialRequest",
        1,
        this.props.addCredentialRequest.bind(this, elem.value),
      );
      this.selectSuccess();
    } catch (e) {
      await this.props.errorToast(e.errorMsg);
    }
    await this.blurSearchBar();
  }

  private async onCredSelect(elem: { id: number, value: string }): Promise<void> {
    try {
      logAnalyticsThenExecuteAsync(
        "Credentials",
        AnalyticsActions.SELECT,
        "Credential",
        1,
        this.props.addCredential.bind(this, elem.value),
      );
      this.selectSuccess();
    } catch (e) {
      await this.props.errorToast(e.errorMsg);
    }
    await this.blurSearchBar();
  }

  private async onRawCredSelect(value: string) {
    try {
      logAnalyticsThenExecuteAsync(
        "Credentials",
        AnalyticsActions.SELECT,
        "RawCredential",
        1,
        this.props.addCredential.bind(this, value),
      )
      this.selectSuccess();
    } catch (e) {
      await this.props.errorToast(e.errorMsg);
    }
    await this.blurSearchBar();
  }

  private async confirmationWrapper(
    confirmCallback: (value: FilterableElement) => Promise<void>,
    negativeCallback: () => Promise<void>,
    value: FilterableElement,
  ): Promise<void> {
    Alert.alert(
        'Request a match',
        'Would you like to request to be matched with "' + value.value + '"?',
        [
          {text: 'Cancel', onPress: async () => {await negativeCallback()}, style: 'cancel'},
          {text: 'Request Match', onPress: async () => {await confirmCallback(value)}, style: 'destructive'},
        ],
      );
  }

  private async onRawCredReqSelect(value: string) {
    try {
      logAnalyticsThenExecuteAsync(
        "Credentials",
        AnalyticsActions.SELECT,
        "RawCredentialRequest",
        1,
        this.props.addCredentialRequest.bind(this, value),
      )
      this.selectSuccess();
    } catch (e) {
      await this.props.errorToast(e.errorMsg);
    }
    await this.blurSearchBar();
  }

  render() {
    if (!this.props.searchBar.hasFocus) return null;

    const { credentials } = this.props.credentialOptions;
    switch (this.props.searchBar.listType) {
      case SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS:
        const reqHint = 'Connect with people of a certain background. ' +
          'E.g. graphic designer or tennis player';
        return (
          <FilterListModal
            hint={reqHint}
            curValue={this.props.searchBar.value}
            data={credentials.map(cred => { return { id: cred.id, value: cred.name }}).toList()}
            onSelect={this.confirmationWrapper.bind(this, this.onReqSelect, async () => {})}
            onRawSelect={this.onRawCredReqSelect}
          />
        );
      case SEARCH_LIST_TYPE_CREDENTIALS:
        const hint = 'Let people find you based on ' +
          'your credentials. E.g. dev at Facebook or marathon runner';
        return (
          <FilterListModal
            hint={hint}
            curValue={this.props.searchBar.value}
            data={credentials.map(cred => { return { id: cred.id, value: cred.name }}).toList()}
            onSelect={this.onCredSelect}
            onRawSelect={this.onRawCredSelect}
          />
        );
      default:
        // Ensure exhaustiveness of select
        const _: never = this.props.searchBar.listType;
        return null;
    }
  }
}


export default connect(({ credentialOptions, searchBar }: RootState) => {
  return { credentialOptions, searchBar };
}, {
  addCredential,
  addCredentialRequest,
  errorToast,
  updateFocus,
})(AllFilterableModals);
