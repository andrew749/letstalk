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
import { ActionTypes as CredentialsActionTypes } from '../redux/credentials/actions';
import { ActionTypes as CredentialRequestsActionTypes } from '../redux/credential-requests/actions';
import { ActionTypes as CredentialOptionsActionTypes } from '../redux/credential-options/actions';
import {
  ActionButton,
  Card,
  Header,
  Loading,
} from '../components';
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
}

interface Props extends DispatchActions {
  credentials: CredentialsState;
  credentialOptions: CredentialOptionsState;
  credentialRequests: CredentialRequestsState;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class RequestToMatchView extends Component<Props> {
  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: <TopHeader placeholder="Find someone who is a..." />,
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

  private renderBody() {
    const { credentials } = this.props.credentialOptions;

    const addCredentialButton = (onPress: () => void) => {
      return (
        <TouchableOpacity onPress={onPress} style={styles.addButton}>
          <MaterialIcons name="add-circle" size={32} color={Colors.HIVE_PRIMARY} />
        </TouchableOpacity>
      );
    }

    return (
      <View style={styles.container}>
        <ScrollView>
          <Header>Active Requests</Header>
          <View style={styles.credentialRequestContainer}>
            {this.renderCredentialRequests()}
          </View>
          <View style={styles.credentialHeaderContainer}>
            <Header>Your Credentials</Header>
          </View>
          <View style={styles.credentialRequestContainer}>
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
})
