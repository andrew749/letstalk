import React, { Component } from 'react';
import {
  Dimensions,
  ScrollView,
  AppRegistry,
  Text,
  TextInput,
  View,
  FlatList,
  StyleSheet,
  Image,
} from 'react-native';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux'
import { NavigationScreenProp, NavigationStackAction, NavigationActions } from 'react-navigation';

import auth from '../services/auth';
import { ActionButton, Card, Header } from '../components';
import { RootState } from '../redux';
import { State as BootstrapState } from '../redux/bootstrap/reducer';

const window = Dimensions.get('window');

interface Props extends BootstrapState {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class ProfileView extends Component<Props> {

  constructor(props: Props) {
    super(props);

    this.onLogoutPress = this.onLogoutPress.bind(this);
  }

  async onLogoutPress() {
    await auth.logout();
    this.props.navigation.dispatch(NavigationActions.reset({
      index: 0,
      key: null,
      actions: [NavigationActions.navigate({ routeName: 'Login' })]
    }));
  }

  render() {
    const {
      programId,
      gradYear,
      sequence,
    } = this.props.bootstrap.cohort;
    return (
      <ScrollView contentContainerStyle={styles.container}>
        <Header title="Profile" />
        <View style={styles.contentContainer} >
          <Image style={styles.image} source={require('../img/profile.jpg')} />
          <Card>
            <Text style={styles.cohort}>Cohort</Text>
            <Text style={styles.cohortText}>{programId}</Text>
            <Text style={styles.cohortText}>{gradYear}</Text>
            <Text style={styles.cohortText}>{sequence}</Text>
          </Card>
        </View>
        <ActionButton onPress={this.onLogoutPress} title='LOGOUT'/>
      </ScrollView>
    );
  }
}

export default connect(({bootstrap}: RootState) => bootstrap)(ProfileView);

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
  cohort: {
    fontWeight: 'bold',
    fontSize: 18,
    marginBottom: 5,
  },
  cohortText: {
    fontSize: 12,
  },
});
