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
import { combineFetchStates } from '../redux/actions';
import {
  State as CredentialsState,
  addCredential,
  fetchCredentials,
} from '../redux/credentials/reducer';
import {
  State as CredentialOptionsState,
  fetchCredentialOptions,
} from '../redux/credential-options/reducer';
import { ActionTypes as CredentialsActionTypes } from '../redux/credentials/actions';
import { ActionTypes as CredentialOptionsActionTypes } from '../redux/credential-options/actions';
import { ActionButton, Card, FilterListModal, Header, Loading } from '../components';
import { CredentialPair, CredentialFilterableElement } from '../models/credential';

interface DispatchActions {
  addCredential: ActionCreator<
    ThunkAction<Promise<CredentialsActionTypes>, CredentialsState, void>>;
  fetchCredentials: ActionCreator<
    ThunkAction<Promise<CredentialsActionTypes>, CredentialsState, void>>;
  fetchCredentialOptions: ActionCreator<
    ThunkAction<Promise<CredentialOptionsActionTypes>, CredentialOptionsState, void>>;
}

interface Props extends DispatchActions {
  credentials: CredentialsState;
  credentialOptions: CredentialOptionsState;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class CredentialEditView extends Component<Props> {
  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: 'Edit Credentials',
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
    await Promise.all([this.props.fetchCredentials(), this.props.fetchCredentialOptions()]);
  }

  private renderCredentials() {
    const { credentialsWithState } = this.props.credentials;
    if (credentialsWithState.isEmpty()) {
      return <Text style={styles.credential}>You haven't added any credentials</Text>;
    }
    return credentialsWithState.map(credentialWithState => {
      const {
        credentialId,
        organizationName,
        positionName,
      } = credentialWithState;
      const name = `${positionName} at ${organizationName}`;
      return (
        <Card key={credentialId} style={styles.credentialCard}>
          <Text style={styles.credential}>{name}</Text>
          <TouchableOpacity style={styles.delete}>
            <MaterialIcons name="delete" size={24} />
          </TouchableOpacity>
        </Card>
      );
    });
  }

  private async onSelect(elem: CredentialFilterableElement): Promise<void> {
    this.props.addCredential(elem);
  }

  private renderBody() {
    const { credentialElements } = this.props.credentialOptions;
    return (
      <ScrollView keyboardShouldPersistTaps={'always'}>
        <Header>Your Credentials</Header>
        <View style={styles.credentialContainer}>
          {this.renderCredentials()}
        </View>
        <Text style={styles.addCredentialText}>Add a credential:</Text>
        <FilterListModal
          data={credentialElements}
          onSelect={this.onSelect}
          placeholder="I am a..."
        />
      </ScrollView>
    );
  }

  render() {
    const {
      state,
      errorMsg,
    } = combineFetchStates(
      this.props.credentials.fetchState,
      this.props.credentialOptions.fetchState,
    );
    return (
      <Loading
        state={state}
        load={this.load}
        errorMsg={errorMsg}
        renderBody={this.renderBody}
      />
    );
  }
}

export default connect(
  ({ credentials, credentialOptions }: RootState) => {
    return { credentials, credentialOptions };
  },
  { addCredential, fetchCredentials, fetchCredentialOptions }
)(CredentialEditView);

const styles = StyleSheet.create({
  credentialContainer: {
    flex: 1,
    alignItems: 'center'
  },
  addCredentialText: {
    fontWeight: 'bold',
    fontSize: 18,
    padding: 10,
  },
  credentialCard: {
    flexDirection: 'row',
  },
  credential: {
    fontWeight: 'bold',
    fontSize: 18,
  },
  // TODO: Make this float right
  delete: {
  },
})
