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
import Immutable from 'immutable';

import { RootState } from '../redux';
import {
  State as CredentialOptionsState,
  fetchCredentialOptions,
} from '../redux/credential-options/reducer';
import { ActionTypes } from '../redux/credential-options/actions';
import { ActionButton, Card, FilterListModal, Header } from '../components';
import { CredentialElement } from '../models/credential';

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
  }

  async componentDidMount() {
    this.load();
  }

  private async load() {
    await this.props.fetchCredentialOptions();
  }

  private onSelect(elem: CredentialElement) {
    console.log(elem);
  }

  private renderSearch(data: Immutable.List<CredentialElement>) {
    return (
      <FilterListModal
        data={data}
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
    switch (state) {
      case 'prefetch':
      case 'fetching':
        // TODO: Separate component for loading pages
        return (
          <View style={styles.centeredContainer}>
            <Text style={styles.headline}>Soon...</Text>
            <ActivityIndicator size="large" />
          </View>
        );
      case 'error':
        // TODO: Separate component for error pages
        return (
          <View style={styles.centeredContainer}>
            <Text style={styles.headline}>Something went wrong :(</Text>
            <Text style={styles.error}>{errorMsg}</Text>
            <ActionButton onPress={() => this.load()} title="Retry" />
          </View>
        );
      case 'success':
        return this.renderSearch(this.props.credentialElements);
      default:
        // Ensure exhaustiveness of select
        const _: never = state;
    }
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
  error: {
    fontWeight: 'bold',
    fontSize: 18,
    color: 'red',
    textAlign: 'center',
  },
  name: {
    fontWeight: 'bold',
    fontSize: 18,
  },
  relationshipType: {
    fontSize: 12,
    color: 'gray',
    marginBottom: 5,
  },
  emailContainer: {
    flex: 1,
    flexDirection: 'row',
  },
  email: {
    paddingTop: 2,
    marginLeft: 5,
    fontSize: 16,
  },
})
