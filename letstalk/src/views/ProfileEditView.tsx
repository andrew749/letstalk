import React, { Component, SFC } from 'react';
import {
  ScrollView,
  StyleSheet,
  Picker,
  Text,
} from 'react-native';
import { reduxForm, Field, InjectedFormProps, SubmissionError } from 'redux-form';
import { connect, ActionCreator } from 'react-redux';
import { FormValidationMessage } from 'react-native-elements';
import { ThunkAction } from 'redux-thunk';
import { bindActionCreators } from 'redux'
import { NavigationScreenProp, NavigationStackAction, NavigationActions } from 'react-navigation';

import auth from '../services/auth';
import {
  ActionButton,
  ButtonPicker,
  Card,
  FormProps,
  Header,
  LabeledFormInput,
  Loading,
  ModalDatePicker,
  ProfileAvatar,
} from '../components';
import { genderIdToString } from '../models/user';
import { RootState } from '../redux';
import { State as BootstrapState, fetchBootstrap } from '../redux/bootstrap/reducer';
import { ActionTypes } from '../redux/bootstrap/actions';
import photoService, {PhotoResult} from '../services/photo_service';
import Colors from '../services/colors';

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

// TODO: move elsewhere
const required = (value: any) => (value ? undefined : 'Required')
const email = (value: string) =>
  value && !/^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$/i.test(value)
    ? 'Invalid email address'
    : undefined
const phoneNumber = (value: string) =>
  value && !/^(0|[1-9][0-9]{9})$/i.test(value)
    ? 'Invalid phone number, must be 10 digits'
    : undefined

const EditForm: SFC<FormProps<EditFormData>> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting, valid } = props;
  return (
    <ScrollView>
      <Field
        name="profilePic"
        component={ProfileAvatar}
        containerStyle={styles.profilePicContainerStyle}
        editable={true}
      />
      <Field
        label="First name"
        name="firstName"
        component={LabeledFormInput}
        autoCorrect={false}
        validate={required}
      />
      <Field
        label="Last name"
        name="lastName"
        component={LabeledFormInput}
        autoCorrect={false}
        validate={required}
      />
      <Field
        label="Email"
        name="email"
        component={LabeledFormInput}
        keyboardType={'email-address' as 'email-address'}
        autoCorrect={false}
        autoCapitalize={'none' as 'none'}
        validate={[required, email]}
      />
      <Field
        label="Phone number"
        name="phoneNumber"
        component={LabeledFormInput}
        keyboardType={'phone-pad' as 'phone-pad'}
        validate={[required, phoneNumber]}
      />
      <Field
        label="Password"
        name="password"
        component={LabeledFormInput}
        secureTextEntry={true}
        validate={required} // Add some rules for password
      />
      <Field
        label="Gender"
        name="gender"
        component={ButtonPicker}
        validate={required}
      >
        <Picker.Item
          label="Male"
          value="male"
        />
        <Picker.Item
          label="Female"
          value="female"
        />
      </Field>
      <Field
        label="Birthday"
        name="birthday"
        mode={'date' as 'date'}
        defaultDate={new Date('1996-11-07T00:00:00.000Z')}
        component={ModalDatePicker}
        validate={required}
      />
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
        backgroundColor={Colors.HIVE_MAIN_BG}
        style={styles.submitButton}
        disabled={!valid}
        loading={submitting}
        title={submitting ? null : "Sign up"}
        onPress={handleSubmit(onSubmit)}
      />
    </ScrollView>
  );
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
    return <Text>Yo</Text>;
  }
}

export default connect(({bootstrap}: RootState) => bootstrap)(ProfileEditView);

const styles = StyleSheet.create({
  submitButton: {
    marginBottom: 100,
  },
  profilePicContainerStyle: {
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: 20,
    marginLeft: 20,
  }
});
