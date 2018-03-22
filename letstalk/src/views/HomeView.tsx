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
import { State as BootstrapState, fetchBootstrap } from '../redux/bootstrap/reducer';
import { ActionTypes } from '../redux/bootstrap/actions';
import { ActionButton, Card, Header } from '../components';

interface DispatchActions {
  fetchBootstrap: ActionCreator<ThunkAction<Promise<ActionTypes>, BootstrapState, void>>;
}

interface Props extends BootstrapState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class HomeView extends Component<Props> {
  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: 'Home',
    headerRight: <ReactNativeButton title="Profile"
      onPress={() => navigation.navigate('Profile')} />,
  })

  constructor(props: Props) {
    super(props);

    this.load = this.load.bind(this);
  }

  async componentDidMount() {
    this.load();
  }

  async componentWillReceiveProps(nextProps: Props) {
    if (nextProps.bootstrap && nextProps.bootstrap.state === 'account_created') {
      // Account not yet setup, so take to onboarding page
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        actions: [NavigationActions.navigate({ routeName: 'Onboarding' })]
      }));
    }
  }

  private async load() {
    await this.props.fetchBootstrap();
  }

  private renderMatches() {
    const { relationships } = this.props.bootstrap;
    const matchCards = relationships.map(relationship => {
      const {
        userId,
        userType,
        firstName,
        lastName,
        email,
      } = relationship;
      // TODO: Handle errors for links
      const link = 'mailto:' + email;
      return (
        <Card key={userId}>
          <Text style={styles.name}>{firstName + ' ' + lastName}</Text>
          <Text style={styles.relationshipType}>{userType}</Text>
          <TouchableOpacity style={styles.emailContainer} onPress={() => Linking.openURL(link)}>
            <MaterialIcons name="email" size={24} />
            <Text style={styles.email}>{email}</Text>
          </TouchableOpacity>
        </Card>
      );
    });
    return (
      <ScrollView>
        <Header>Matches</Header>
        {matchCards}
      </ScrollView>
    );
  }

  private renderHome() {
    const { state } = this.props.bootstrap;
    switch (state) {
      case 'account_created':
        // Should in reality not be shown, since we never show home page for account_created
        return (
          <View style={styles.centeredContainer}>
            <Text style={styles.headline}>Waiting for you to finish onboarding</Text>
            <ActionButton onPress={() => this.load()} title="Check again" />
          </View>
        );
      case 'account_setup':
        return (
          <View style={styles.centeredContainer}>
            <Text style={styles.headline}>Waiting for your match</Text>
            <ActionButton onPress={() => this.load()} title="Check again" />
          </View>
        );
      case 'account_matched':
        return this.renderMatches();
      default:
        // Ensure exhaustiveness of select
        const _: never = state;
    }
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
        return this.renderHome();
      default:
        // Ensure exhaustiveness of select
        const _: never = state;
    }
  }
}

export default connect(({ bootstrap }: RootState) => bootstrap, { fetchBootstrap })(HomeView);

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
