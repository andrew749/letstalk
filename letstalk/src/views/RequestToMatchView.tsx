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
import {
  ActionButton,
  Card,
  FilterableElement,
  FilterListModal,
  Header,
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
}

interface Props extends DispatchActions {
  credentialOptions: CredentialOptionsState;
  credentialRequests: CredentialRequestsState;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

const Head: SFC<{}> = (props: {}) => {
  return <View style={styles.header}/>;
}

class RequestToMatchView extends Component<Props> {
  static navigationOptions = () => ({
    header: <Head/>,
    headerStyle: {
      backgroundColor: Colors.HIVE_MAIN_BG,
    },
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
              <TouchableOpacity onPress={onPress} style={styles.delete}>
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

  private async onSelect(elem: FilterableElement): Promise<void> {
    try {
      await this.props.addCredentialRequest({ id: elem.id, name: elem.value });
    } catch (e) {
      await this.props.errorToast(e.message);
    }
  }

  private renderBody() {
    const { credentials } = this.props.credentialOptions;
    return (
      <View>
        <View style={styles.topContainer}>
          <FilterListModal
            data={credentials.map(cred => { return { id: cred.id, value: cred.name }}).toList()}
            onSelect={this.onSelect}
            placeholder="Find someone who is a..."
          />
        </View>
        <ScrollView keyboardShouldPersistTaps={'always'}>
          <Header>Active Requests</Header>
          <View style={styles.credentialRequestContainer}>
            {this.renderCredentialRequests()}
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
  searchBar: {

  },
  header: {
    height: Platform.OS == "ios" ? 20 : 0,
    backgroundColor: Colors.HIVE_MAIN_BG,
  },
  // TODO: Make this float right
  delete: {
  },
  topContainer: {
    width: SCREEN_WIDTH,
    backgroundColor: '#FFC107',
  },
})
