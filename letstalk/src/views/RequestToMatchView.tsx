import React, { Component, SFC } from 'react';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import {
  ActivityIndicator,
  Button as ReactNativeButton,
  Dimensions,
  Linking,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
  ScrollView,
  Platform,
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
  State as CredentialsState,
  addCredential,
  fetchCredentials,
  removeCredential,
} from '../redux/credentials/reducer';
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
import { ActionTypes as CredentialsActionTypes } from '../redux/credentials/actions';
import { ActionTypes as CredentialRequestsActionTypes } from '../redux/credential-requests/actions';
import { ActionTypes as CredentialOptionsActionTypes } from '../redux/credential-options/actions';

import {
  ActionButton,
  Card,
  FilterableElement,
  FilterListModal,
  Header,
  TextHeader,
  Loading,
} from '../components';
import { Credential } from '../models/credential';
import Colors from '../services/colors';

const SCREEN_WIDTH = Dimensions.get('window').width;

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
  addCredential: ActionCreator<
    ThunkAction<Promise<CredentialsActionTypes>, CredentialsState, void>>;
  removeCredential: ActionCreator<
    ThunkAction<Promise<CredentialsActionTypes>, CredentialsState, void>>;
  fetchCredentials: ActionCreator<
    ThunkAction<Promise<CredentialsActionTypes>, CredentialsState, void>>;
}

interface Props extends DispatchActions {
  credentials: CredentialsState;
  credentialOptions: CredentialOptionsState;
  credentialRequests: CredentialRequestsState;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}


class RequestToMatchView extends Component<Props> {

  static navigationOptions = ({ navigation }: NavigationScreenDetails<void> ) => ({
    header: (props: any) => <Header {...props}></Header>
  });

  constructor(props: Props) {
    super(props);

    this.load = this.load.bind(this);
    this.renderBody = this.renderBody.bind(this);
    this.onReqSelect = this.onReqSelect.bind(this);
    this.onCredSelect = this.onCredSelect.bind(this);
    this.onRawCredSelect = this.onRawCredSelect.bind(this);
  }

  async componentDidMount() {
    this.load();
  }

  private async load() {
    await Promise.all([
      this.props.fetchCredentials(),
      this.props.fetchCredentialRequests(),
      this.props.fetchCredentialOptions(),
    ]);
  }

  private renderCredentials() {
    const { credentialsWithState } = this.props.credentials;
    if (credentialsWithState.isEmpty()) {
      return <Text style={styles.credential}>You haven't added any credentials</Text>;
    }
    return credentialsWithState.map(credentialWithState => {
      const {
        id,
        name,
        state,
      } = credentialWithState;
      switch (state) {
        case 'normal':
          const onPress = async () => {
            try {
              await this.props.removeCredential(id);
            } catch(e) {
              await this.props.errorToast(e.message);
            }
          };
          return (
            <Card key={id} style={styles.credentialCard}>
              <Text style={styles.credential}>{name}</Text>
              <TouchableOpacity onPress={onPress}>
                <MaterialIcons name="delete" size={24} />
              </TouchableOpacity>
            </Card>
          );
        case 'deleting':
          return (
            <Card key={id} style={styles.deletingCard}>
              <ActivityIndicator />
            </Card>
          );
        default:
          // Ensure exhaustiveness of select
          const _: never = state;
      }
    });
  }

  private renderCredentialRequests() {
    const { credentialRequestsWithState } = this.props.credentialRequests;
    if (credentialRequestsWithState.isEmpty()) {
      return <Text style={styles.credentialRequest}>You don't have any requests</Text>;
    }
    return credentialRequestsWithState.map(credentialWithState => {
      const {
        id,
        name,
        state,
      } = credentialWithState;
      switch (state) {
        case 'normal':
          const onPress = async () => {
            try {
              await this.props.removeCredentialRequest(id);
            } catch(e) {
              await this.props.errorToast(e.message);
            }
          };
          return (
            <Card key={id} style={styles.credentialRequestCard}>
              <Text style={styles.credentialRequest}>{name}</Text>
              <TouchableOpacity onPress={onPress}>
                <MaterialIcons name="delete" size={24} />
              </TouchableOpacity>
            </Card>
          );
        case 'deleting':
          return (
            <Card key={id} style={styles.deletingCard}>
              <ActivityIndicator />
            </Card>
          );
        default:
          // Ensure exhaustiveness of select
          const _: never = state;
      }
    });
  }

  private async onReqSelect(elem: FilterableElement): Promise<void> {
    try {
      await this.props.addCredentialRequest({ id: elem.id, name: elem.value });
    } catch (e) {
      await this.props.errorToast(e.message);
    }
  }

  private async onCredSelect(elem: { id: number, value: string }): Promise<void> {
    try {
      await this.props.addCredential(elem.value);
    } catch (e) {
      await this.props.errorToast(e.message);
    }
  }

  private async onRawCredSelect(value: string) {
    try {
      await this.props.addCredential(value);
    } catch (e) {
      await this.props.errorToast(e.message);
    }
  }

  private renderBody() {
    const { credentialOptions } = this.props.credentialOptions;

    const addCredentialButton = (onPress: () => void) => {
      return (
        <TouchableOpacity onPress={onPress} style={styles.addButton}>
          <MaterialIcons name="add-circle" size={32} color={Colors.HIVE_PRIMARY} />
        </TouchableOpacity>
      );
    }

    return (
      <View style={styles.container}>
        <ScrollView keyboardShouldPersistTaps={'always'}>
          <TextHeader>Active Requests</TextHeader>
          <View style={styles.credentialRequestContainer}>
            {this.renderCredentialRequests()}
          </View>
          <View style={styles.credentialHeaderContainer}>
            <TextHeader>Your Credentials</TextHeader>
            <FilterListModal
              data={credentialOptions.map(cred => { return { id: cred.id, value: cred.name }}).toList()}
              onSelect={this.onCredSelect}
              onRawSelect={this.onRawCredSelect}
              placeholder="I am a..."
              buttonComponent={addCredentialButton}
            />
          </View>
          <View style={styles.credentialRequestContainer}>
            {this.renderCredentials()}
          </View>
        </ScrollView>
      </View>
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
  ({ credentialRequests, credentialOptions, credentials }: RootState) => {
    return { credentialOptions, credentialRequests, credentials };
  }, {
    addCredential,
    addCredentialRequest,
    errorToast,
    fetchCredentials,
    fetchCredentialRequests,
    fetchCredentialOptions,
    removeCredential,
    removeCredentialRequest,
  })(RequestToMatchView);

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
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
  credentialHeaderContainer: {
    justifyContent: 'space-between',
    flexDirection: 'row',
    width: SCREEN_WIDTH,
  },
  credentialCard: {
    flexDirection: 'row',
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
  credential: {
    fontWeight: 'bold',
    fontSize: 18,
  },
  addButton: {
    margin: 12,
  },
  topContainer: {
    width: SCREEN_WIDTH,
    backgroundColor: '#FFC107',
  },
})
