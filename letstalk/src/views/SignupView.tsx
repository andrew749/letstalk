import React, { Component } from 'react';
import { Picker, ScrollView } from 'react-native';
import {
  NavigationScreenProp,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import { reduxForm, Field, InjectedFormProps, SubmissionError } from 'redux-form';
import { FormValidationMessage } from 'react-native-elements';

import {
  ActionButton,
  FormP,
  FormProps,
  LabeledFormInput,
  ModalPicker,
  ModalDatePicker,
} from '../components';
import profileService from '../services/profile-service';

interface SignupFormData {
  firstName: string;
  lastName: string;
  email: string;
  phoneNumber: string;
  password: string;
  gender: string;
  birthday: Date;
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
    console.log(values);
    await onSubmit(values);
    reset();
  };
  return (
    <ScrollView>
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
        component={ModalPicker}
        validate={required} // Add some rules for password
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
        validate={required} // Add some rules for password
      />
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
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
  static navigationOptions = {
    headerTitle: 'Signup',
  }

  constructor(props: Props) {
    super(props);

    this.onSubmit = this.onSubmit.bind(this);
  }

  async onSubmit(values: SignupFormData) {
    try {
      await profileService.signup({
        ...values,
        birthday: Math.round(values.birthday.getTime() / 1000),
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
    return <SignupFormWithRedux onSubmit={this.onSubmit} />;
  }
}
