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
import { ActionButton, Header } from '../components';

const window = Dimensions.get('window');

interface Props {
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
    const title = `My Profile`;
    const placeholderText = `Lorem Ipsum`;

    return(
      <ScrollView contentContainerStyle={styles.container}>
        <Header title="Profile" />
        <View style={styles.contentContainer} >
          <Image style={styles.image} source={require('../img/profile.jpg')} />
        </View>
        <ActionButton onPress={this.onLogoutPress} title='LOGOUT'/>
      </ScrollView>
    );
  }
}

export default ProfileView;

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
  formInput: {
    width: window.width * .8
  },
  row : {
    flexDirection: 'row'
  },
  unit: {
    flex: 1,
  }
});
