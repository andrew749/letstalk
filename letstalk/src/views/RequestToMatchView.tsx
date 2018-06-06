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
  fetchCredentials,
  removeCredential,
} from '../redux/credentials/reducer';
import {
  State as CredentialRequestsState,
  fetchCredentialRequests,
  removeCredentialRequest,
} from '../redux/credential-requests/reducer';
import {
  State as CredentialOptionsState,
  fetchCredentialOptions,
} from '../redux/credential-options/reducer';
import {
  ActionTypes as SearchBarActionTypes,
  SEARCH_LIST_TYPE_CREDENTIALS,
  SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS,
} from '../redux/search-bar/actions';
import {
  State as SearchBarState,
  updateFocus,
  updateListType
} from '../redux/search-bar/reducer';
import { ActionTypes as CredentialsActionTypes } from '../redux/credentials/actions';
import { ActionTypes as CredentialRequestsActionTypes } from '../redux/credential-requests/actions';
import { ActionTypes as CredentialOptionsActionTypes } from '../redux/credential-options/actions';
import {
  ActionButton,
  Button,
  Card,
  Header,
} from '../components';
import Loading from './Loading';
import { Credential } from '../models/credential';
import Colors from '../services/colors';
import TopHeader, { headerStyle } from './TopHeader';
import AllFilterableModals from './AllFilterableModals';

const SCREEN_WIDTH = Dimensions.get('window').width;

interface DispatchActions {
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  removeCredentialRequest: ActionCreator<
    ThunkAction<Promise<CredentialRequestsActionTypes>, CredentialRequestsState, void>>;
  fetchCredentialRequests: ActionCreator<
    ThunkAction<Promise<CredentialRequestsActionTypes>, CredentialRequestsState, void>>;
  fetchCredentialOptions: ActionCreator<
    ThunkAction<Promise<CredentialOptionsActionTypes>, CredentialOptionsState, void>>;
  removeCredential: ActionCreator<
    ThunkAction<Promise<CredentialsActionTypes>, CredentialsState, void>>;
  fetchCredentials: ActionCreator<
    ThunkAction<Promise<CredentialsActionTypes>, CredentialsState, void>>;
  updateFocus: ActionCreator<
    ThunkAction<Promise<SearchBarActionTypes>, SearchBarState, void>>;
  updateListType: ActionCreator<
    ThunkAction<Promise<SearchBarActionTypes>, SearchBarState, void>>;
}

interface Props extends DispatchActions {
  credentials: CredentialsState;
  credentialOptions: CredentialOptionsState;
  credentialRequests: CredentialRequestsState;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class RequestToMatchView extends Component<Props> {
  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: <TopHeader navigation={navigation} />,
    headerStyle,
  })

  constructor(props: Props) {
    super(props);

    this.load = this.load.bind(this);
    this.renderBody = this.renderBody.bind(this);
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
      const onPress = () => {
        this.props.updateFocus(true);
        this.props.updateListType(SEARCH_LIST_TYPE_CREDENTIALS);
      };
      return (
        <View style={styles.noCredentialsContainer}>
          <Text style={styles.noCredentials}>You haven't added any credentials</Text>
          <Button buttonStyle={styles.noCredentialsButton} title="Add credential" onPress={onPress} />
        </View>
      );
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
              await this.props.errorToast(e.errorMsg);
            }
          };
          return (
            <Card key={id} style={styles.credentialCard}>
              <Text style={styles.credential}>{name}</Text>
              <TouchableOpacity style={styles.deleteCredential} onPress={onPress}>
                <MaterialIcons name="close" size={18} />
              </TouchableOpacity>
            </Card>
          );
        case 'deleting':
          return (
            <Card key={id} style={styles.credentialCard}>
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
      const onPress = () => {
        this.props.updateFocus(true);
        this.props.updateListType(SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS);
      };
      return (
        <View style={styles.noCredentialsContainer}>
          <Text style={styles.noCredentials}>You don't have any requests</Text>
          <Button buttonStyle={styles.noCredentialsButton} title="Add request" onPress={onPress} />
        </View>
      );
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
              await this.props.errorToast(e.errorMsg);
            }
          };
          return (
            <Card key={id} style={styles.credentialCard}>
              <Text style={styles.credential}>{name}</Text>
              <TouchableOpacity style={styles.deleteCredential} onPress={onPress}>
                <MaterialIcons name="close" size={18} />
              </TouchableOpacity>
            </Card>
          );
        case 'deleting':
          return (
            <Card key={id} style={styles.credentialCard}>
              <ActivityIndicator />
            </Card>
          );
        default:
          // Ensure exhaustiveness of select
          const _: never = state;
      }
    });
  }

  private renderBody() {
    const { credentials } = this.props.credentialOptions;

    const onAddCredentialPress = () => {
      this.props.updateFocus(true);
      this.props.updateListType(SEARCH_LIST_TYPE_CREDENTIALS);
    }

    return (
      <View style={styles.container}>
        <ScrollView>
          <Header>Active Requests</Header>
          <View style={styles.credentialContainer}>
            {this.renderCredentialRequests()}
          </View>
          <View style={styles.credentialHeaderContainer}>
            <Header>Your Credentials</Header>
            <TouchableOpacity onPress={onAddCredentialPress} style={styles.addButton}>
              <MaterialIcons name="add-circle" size={32} color={Colors.HIVE_PRIMARY} />
            </TouchableOpacity>
          </View>
          <View style={styles.credentialContainer}>
            {this.renderCredentials()}
          </View>
        </ScrollView>
        <AllFilterableModals />
      </View>
    );
  }

  render() {
    const {
      state,
      errorMsg,
      errorType,
    } = combineFetchStates(
      this.props.credentialRequests.fetchState,
      this.props.credentialOptions.fetchState,
    );
    return (
      <Loading
        state={state}
        errorMsg={errorMsg}
        errorType={errorType}
        load={this.load}
        renderBody={this.renderBody}
        navigation={this.props.navigation}
      />
    );
  }
}

export default connect(
  ({ credentialRequests, credentialOptions, credentials }: RootState) => {
    return { credentialOptions, credentialRequests, credentials };
  }, {
    errorToast,
    fetchCredentials,
    fetchCredentialRequests,
    fetchCredentialOptions,
    removeCredential,
    removeCredentialRequest,
    updateFocus,
    updateListType,
  })(RequestToMatchView);

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  credentialContainer: {
    flex: 1,
    alignItems: 'stretch'
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
    padding: 10,
    marginVertical: 5,
    marginHorizontal: 10,
  },
  noCredentialsContainer: {
    alignSelf: 'center',
    alignItems: 'center'
  },
  noCredentialsButton: {
    width: 200,
    marginTop: 10,
  },
  noCredentials: {
    fontSize: 14,
  },
  deleteCredential: {
    position: 'absolute',
    top: 10,
    right: 10,
  },
  credential: {
    fontSize: 14,
    paddingRight: 30,
  },
  addButton: {
    margin: 12,
  },
})
