import React, { Component, SFC } from 'react';
import { StyleSheet, Text } from 'react-native';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import { bindActionCreators } from 'redux'
import { NavigationScreenProp, NavigationStackAction, NavigationActions } from 'react-navigation';

import auth from '../services/auth';
import {
  ActionButton,
  Card,
  FormProps,
  Header,
  Loading,
} from '../components';
import { genderIdToString } from '../models/user';
import { RootState } from '../redux';
import { State as BootstrapState, fetchBootstrap } from '../redux/bootstrap/reducer';
import { ActionTypes } from '../redux/bootstrap/actions';
import photoService, {PhotoResult} from '../services/photo_service';

interface EditFormData {
  firstName: string;
  lastName: string;
  email: string;
  phoneNumber: string;
  password: string;
  gender: string;
  birthday: Date;
  profilePic: PhotoResult;
}

const EditForm: SFC<FormProps<EditFormData>> = props => {
  return <Text>Yo</Text>;
}

interface DispatchActions {
}

interface Props extends BootstrapState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class ProfileEditView extends Component<Props> {
  static navigationOptions = () => ({
    headerTitle: 'Edit Profile',
  })

  render() {
    return <Text>Yo</Text>
  }
}

export default connect(({bootstrap}: RootState) => bootstrap)(ProfileEditView);

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
