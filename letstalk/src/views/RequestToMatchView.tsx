import React, { Component } from 'react';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import {
  ActivityIndicator,
  Button as ReactNativeButton,
  Linking,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
  ScrollView,
} from 'react-native';
import {
  NavigationScreenProp,
  NavigationScreenDetails,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import { MaterialIcons } from '@expo/vector-icons';

import { RootState } from '../redux';
import { errorToast } from '../redux/toast';
import { combineFetchStates } from '../redux/actions';
import {
  State as CredentialRequestsState,
  addCredentialRequest,
  fetchCredentialRequests,
  removeCredentialRequest,
} from '../redux/credential-requests/reducer';
import {
  State as CredentialOptionsState,
  fetchCredentialOptions,
} from '../redux/credential-options/reducer';
import { ActionTypes as CredentialRequestsActionTypes } from '../redux/credential-requests/actions';
import { ActionTypes as CredentialOptionsActionTypes } from '../redux/credential-options/actions';
import { ActionButton, Card, FilterListModal, Header, Loading } from '../components';
import { CredentialFilterableElement } from '../models/credential';

interface DispatchActions {
  addCredentialRequest: ActionCreator<
    ThunkAction<Promise<CredentialRequestsActionTypes>, CredentialRequestsState, void>>;
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  removeCredentialRequest: ActionCreator<
    ThunkAction<Promise<CredentialRequestsActionTypes>, CredentialRequestsState, void>>;
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

class RequestToMatchView extends Component<Props> {
  static navigationOptions = () => ({
    headerTitle: 'Request To Match',
  })

  constructor(props: Props) {
    super(props);

    this.load = this.load.bind(this);
    this.renderBody = this.renderBody.bind(this);
    this.onSelect = this.onSelect.bind(this);
  }

  async componentDidMount() {
    this.load();
  }

  private async load() {
    await Promise.all([this.props.fetchCredentialRequests(), this.props.fetchCredentialOptions()]);
  }

  private renderCredentialRequests() {
    const { credentialRequestsWithState } = this.props.credentialRequests;
    if (credentialRequestsWithState.isEmpty()) {
      return <Text style={styles.credentialRequest}>You don't have any requests</Text>;
    }
    return credentialRequestsWithState.map(credentialWithState => {
      const {
        state,
        credentialRequestId,
        organizationName,
        positionName,
      } = credentialWithState;
      switch (state) {
        case 'normal':
          const name = `${positionName} at ${organizationName}`;
          const onPress = async () => {
            try {
              await this.props.removeCredentialRequest(credentialRequestId);
            } catch(e) {
              await this.props.errorToast(e.message);
            }
          };
          return (
            <Card key={credentialRequestId} style={styles.credentialRequestCard}>
              <Text style={styles.credentialRequest}>{name}</Text>
              <TouchableOpacity onPress={onPress} style={styles.delete}>
                <MaterialIcons name="delete" size={24} />
              </TouchableOpacity>
            </Card>
          );
        case 'deleting':
          return (
            <Card key={credentialRequestId} style={styles.deletingCard}>
              <ActivityIndicator />
            </Card>
          );
        default:
          // Ensure exhaustiveness of select
          const _: never = state;
      }
    });
  }

  private async onSelect(elem: CredentialFilterableElement): Promise<void> {
    try {
      await this.props.addCredentialRequest(elem);
    } catch (e) {
      await this.props.errorToast(e.message);
    }
  }

  private renderBody() {
    const { credentialElements } = this.props.credentialOptions;
    return (
      <ScrollView keyboardShouldPersistTaps={'always'}>
        <FilterListModal
          data={credentialElements}
          onSelect={this.onSelect}
          placeholder="Find someone who is a..."
        />
        <Header>Active Requests</Header>
        <View style={styles.credentialRequestContainer}>
          {this.renderCredentialRequests()}
        </View>
      </ScrollView>
    );
  }

  render() {
    const {
      state,
      errorMsg,
    } = combineFetchStates(
      this.props.credentialRequests.fetchState,
      this.props.credentialOptions.fetchState,
    );
    return (
      <Loading
        state={state}
        errorMsg={errorMsg}
        load={this.load}
        renderBody={this.renderBody}
      />
    );
  }
}

export default connect(
  ({ credentialRequests, credentialOptions }: RootState) => {
    return { credentialOptions, credentialRequests };
  }, {
    addCredentialRequest,
    errorToast,
    fetchCredentialRequests,
    fetchCredentialOptions,
    removeCredentialRequest,
  })(RequestToMatchView);

const styles = StyleSheet.create({
  credentialRequestContainer: {
    flex: 1,
    alignItems: 'center'
  },
  headline: {
    fontWeight: 'bold',
    fontSize: 18,
    textAlign: 'center',
    marginBottom: 10,
  },
  credentialRequestCard: {
    flexDirection: 'row',
  },
  deletingCard: {
    alignItems: 'center',
  },
  credentialRequest: {
    fontWeight: 'bold',
    fontSize: 18,
  },
  // TODO: Make this float right
  delete: {
  },
})
