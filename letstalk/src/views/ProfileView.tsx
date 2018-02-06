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
import { Button, FormLabel, FormInput, FormValidationMessage } from 'react-native-elements'
import { NavigationScreenProp, NavigationStackAction, NavigationActions } from 'react-navigation';

import auth from '../services/auth';

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
      <ScrollView contentContainerStyle={ styles.container }>
        <Text style= { styles.textInput }>{ title } </Text>
        <Image style= {styles.image} source={require('../img/profile.jpg')} />
        <View>
          <FormLabel>Name</FormLabel>
          <FormInput containerStyle = {styles.formInput} placeholder={placeholderText}/>
        </View>
        <View>
          <FormLabel>Program</FormLabel>
          <FormInput containerStyle = {styles.formInput} placeholder={placeholderText}/>
        </View>
        <View style = {styles.row} >
          <View style = {styles.unit} >
            <FormLabel>Term</FormLabel>
            <FormInput placeholder={placeholderText}/>
          </View>
          <View style = {styles.unit}>
            <FormLabel>Stream</FormLabel>
            <FormInput placeholder={placeholderText}/>
          </View>
        </View>
        <View>
          <FormLabel>Email</FormLabel>
          <FormInput containerStyle = {styles.formInput} placeholder={placeholderText}/>
        </View>
        <View>
          <FormLabel>Phone Number</FormLabel>
          <FormInput containerStyle = {styles.formInput} placeholder={placeholderText}/>
        </View>
        <Button onPress={this.onLogoutPress} title='LOGOUT' backgroundColor='#EB5757'/>
      </ScrollView>
    );
  }
}

export default ProfileView;

const styles = StyleSheet.create({
  container: {
    alignItems: 'center',
    marginHorizontal: 25
  },
  textInput: {
    padding: 20
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
