import {FileSystem} from 'expo';
import React, { Component } from 'react';
import { Picker, ScrollView, StyleSheet, KeyboardAvoidingView } from 'react-native';
import {
  NavigationScreenProp,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import { reduxForm, Field, InjectedFormProps, SubmissionError } from 'redux-form';
import { FormValidationMessage } from 'react-native-elements';

import {
  ActionButton,
  ButtonPicker,
  FormP,
  FormProps,
  LabeledFormInput,
  ModalPicker,
  ModalDatePicker,
  ProfileAvatarEditableFormElement,
} from '../components';
import profileService, {SignupRequest} from '../services/profile-service';
import photoService, {PhotoResult} from '../services/photo_service';
import Colors from '../services/colors';
import {AnalyticsHelper} from '../services/analytics';

interface SignupFormData {
  firstName: string;
  lastName: string;
  email: string;
  phoneNumber: string;
  password: string;
  gender: number;
  birthdate: string;
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

const SignupForm: React.SFC<FormProps<SignupFormData>> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting, valid } = props;
  const onSubmitWithReset = async (values: SignupFormData): Promise<void> => {
    await onSubmit(values);
    reset();
  };
  return (
    <ScrollView>
      <Field
        name="profilePic"
        component={ProfileAvatarEditableFormElement}
        containerStyle={styles.profilePicContainerStyle}
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
          value={2}
        />
        <Picker.Item
          label="Female"
          value={1}
        />
      </Field>
      <Field
        label="Birthday"
        name="birthdate"
        mode={'date' as 'date'}
        defaultDate={new Date('1996-11-07T00:00:00.000Z')}
        component={ModalDatePicker}
        validate={required}
      />
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
        backgroundColor={Colors.HIVE_PRIMARY}
        style={styles.submitButton}
        disabled={!valid}
        loading={submitting}
        title={submitting ? null : "Sign up"}
        onPress={handleSubmit(onSubmitWithReset)}
      />
    </ScrollView>
  );
}

const SignupFormWithRedux = reduxForm<SignupFormData, FormP<SignupFormData>>({
  form: 'signup',
})(SignupForm);

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

export default class SignupView extends Component<Props> {
  SIGNUP_VIEW_IDENTIFIER = "SignupView";

  static navigationOptions = {
    headerTitle: 'Signup',
    headerStyle: {
      backgroundColor: Colors.HIVE_PRIMARY,
    }
  }

  constructor(props: Props) {
    super(props);

    this.onSubmit = this.onSubmit.bind(this);
  }

  async componentDidMount() {
    AnalyticsHelper.getInstance().recordPage(this.SIGNUP_VIEW_IDENTIFIER);
  }

  async onSubmit(values: SignupFormData) {
    try {
      await profileService.signup({
        ...{
          "firstName": values.firstName,
          "lastName": values.lastName,
          "email": values.email,
          "phoneNumber": values.phoneNumber,
          "gender": values.gender,
          "password": values.password,
          "profilePic": values.profilePic ? values.profilePic.data: undefined,
          "birthdate": values.birthdate,
        }
      });
      // TODO: have a prompt saying successfully signed up
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        key: null,
        actions: [NavigationActions.navigate({ routeName: 'Login' })]
      }));
    } catch(e) {
      throw new SubmissionError({_error: e.message});
    }
  }

  render() {
    return (
      <KeyboardAvoidingView behavior="padding">
        <SignupFormWithRedux onSubmit={this.onSubmit} />
      </KeyboardAvoidingView>
    );
  }
}

// Hack to make submit button available when keyboard is open.
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
