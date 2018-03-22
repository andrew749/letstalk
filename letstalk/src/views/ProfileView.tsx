import React, { Component } from 'react';
import {
  ActivityIndicator,
  ScrollView,
  AppRegistry,
  Text,
  TextInput,
  View,
  FlatList,
  StyleSheet,
  Image,
} from 'react-native';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import { bindActionCreators } from 'redux'
import { NavigationScreenProp, NavigationStackAction, NavigationActions } from 'react-navigation';

import auth from '../services/auth';
import { ActionButton, Card, Header } from '../components';
import { genderIdToString } from '../models';
import { RootState } from '../redux';
import { State as BootstrapState, fetchBootstrap } from '../redux/bootstrap/reducer';
import { ActionTypes } from '../redux/bootstrap/actions';

interface DispatchActions {
  fetchBootstrap: ActionCreator<ThunkAction<Promise<ActionTypes>, BootstrapState, void>>;
}

interface Props extends BootstrapState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class ProfileView extends Component<Props> {

  constructor(props: Props) {
    super(props);

    this.onLogoutPress = this.onLogoutPress.bind(this);
    this.load = this.load.bind(this);
  }

  private async onLogoutPress() {
    try {
      await auth.logout();
    } catch(error) {}
    this.props.navigation.dispatch(NavigationActions.reset({
      index: 0,
      key: null,
      actions: [NavigationActions.navigate({ routeName: 'Login' })]
    }));
  }

  private async load() {
    await this.props.fetchBootstrap();
  }

  renderBody() {
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
        // TODO: right now the code assumes that cohort always exists. remove this assumption.
        const {
          programId,
          gradYear,
          sequence,
        } = this.props.bootstrap.cohort;

        const {
          gender,
          email,
          birthdate,
        } = this.props.bootstrap.me;

        const capitalize = (s: string) => s.charAt(0).toUpperCase() + s.slice(1);

        const genderStr = capitalize(genderIdToString(gender));
        const options = { year: 'numeric', month: 'long', day: 'numeric' };
        const birthdayStr = birthdate.toLocaleDateString('en-US', options);

        return (
          <View style={styles.contentContainer} >
            <Image style={styles.image} source={require('../img/profile.jpg')} />
            <Card>
              <Text style={styles.sectionHeader}>Profile</Text>
              <View style={styles.listItem}>
                <Text style={styles.label}>Gender:</Text>
                <Text style={styles.value}>{genderStr}</Text>
              </View>
              <View style={styles.listItem}>
                <Text style={styles.label}>Email:</Text>
                <Text style={styles.value}>{email}</Text>
              </View>
              <View style={styles.listItem}>
                <Text style={styles.label}>Birthday:</Text>
                <Text style={styles.value}>{birthdayStr}</Text>
              </View>
            </Card>
            <Card>
              <Text style={styles.sectionHeader}>Cohort</Text>
              <View style={styles.listItem}>
                <Text style={styles.label}>Program:</Text>
                <Text style={styles.value}>{programId}</Text>
              </View>
              <View style={styles.listItem}>
                <Text style={styles.label}>Sequence:</Text>
                <Text style={styles.value}>{sequence}</Text>
              </View>
              <View style={styles.listItem}>
                <Text style={styles.label}>Grad year:</Text>
                <Text style={styles.value}>{gradYear}</Text>
              </View>
            </Card>
          </View>
        );
      default:
        // Ensure exhaustiveness of select
        const _: never = state;
    }
  }

  render() {
    const body = this.renderBody();
    const headerText = this.props.bootstrap ?
      this.props.bootstrap.me.firstName + ' ' + this.props.bootstrap.me.lastName : 'Profile';
    return (
      <ScrollView contentContainerStyle={styles.container}>
        <Header>{headerText}</Header>
        {body}
        <ActionButton onPress={this.onLogoutPress} title='LOGOUT'/>
      </ScrollView>
    );
  }
}

export default connect(({bootstrap}: RootState) => bootstrap, { fetchBootstrap })(ProfileView);

const styles = StyleSheet.create({
  centeredContainer: {
    flex: 1,
    justifyContent: 'center'
  },
  container: {
    paddingBottom: 10,
  },
  contentContainer: {
    alignItems: 'center',
    marginHorizontal: 25
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
  image: {
    width: 150,
    height: 150,
    borderRadius: 75
  },
  listItem: {
    flex: 1,
    flexDirection: 'row',
  },
  sectionHeader: {
    fontWeight: 'bold',
    fontSize: 18,
    marginBottom: 5,
  },
  label: {
    fontWeight: 'bold',
    fontSize: 12,
  },
  value: {
    fontSize: 12,
    marginLeft: 10,
  },
});
