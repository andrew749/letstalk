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
import { ActionButton, Card, Header, Loading } from '../components';
import { genderIdToString } from '../models/user';
import { RootState } from '../redux';
import { State as BootstrapState, fetchBootstrap } from '../redux/bootstrap/reducer';
import { ActionTypes } from '../redux/bootstrap/actions';
import { programById, sequenceById } from '../models/cohort';

interface DispatchActions {
  fetchBootstrap: ActionCreator<ThunkAction<Promise<ActionTypes>, BootstrapState, void>>;
}

interface Props extends BootstrapState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class ProfileView extends Component<Props> {
  static navigationOptions = () => ({
    headerTitle: 'Profile',
  })

  constructor(props: Props) {
    super(props);

    this.onLogoutPress = this.onLogoutPress.bind(this);
    this.load = this.load.bind(this);
    this.renderInner = this.renderInner.bind(this);
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

  private renderInner() {
    const {
      programId,
      gradYear,
      sequenceId,
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
    const sequence = sequenceById(sequenceId);
    const program = programById(programId);

    const buildItems = (name_values: Array<[string, string]>) => {
      return name_values.map(([label, value]) => {
        return (
          <View key={label} style={styles.listItem}>
            <Text style={styles.label}>{label}:</Text>
            <Text style={styles.value}>{value}</Text>
          </View>
        )
      });
    };

    const profileItems = buildItems([
      ['Gender', genderStr],
      ['Email', email],
      ['Birthday', birthdayStr],
    ]);
    const cohortItems = buildItems([
      ['Program', program],
      ['Sequence', sequence],
      ['Grad year', String(gradYear)],
    ]);

    return (
      <View style={styles.contentContainer} >
        <Image style={styles.image} source={require('../img/profile.jpg')} />
        <Card>
          <Text style={styles.sectionHeader}>Profile</Text>
          {profileItems}
        </Card>
        <Card>
          <Text style={styles.sectionHeader}>Cohort</Text>
          {cohortItems}
        </Card>
      </View>
    );
  }

  renderBody() {
    const {
      state,
      errorMsg,
    } = this.props.fetchState;
    return (
      <Loading
        state={state}
        errorMsg={errorMsg}
        load={this.load}
        renderBody={this.renderInner}
      />
    );
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
  container: {
    paddingBottom: 10,
  },
  contentContainer: {
    alignItems: 'center',
    marginHorizontal: 25
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
