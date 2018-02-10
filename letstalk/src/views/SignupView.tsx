import React, { Component } from 'react';
import { ScrollView } from 'react-native';
import {
  NavigationScreenProp,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import { reduxForm, Field, InjectedFormProps, SubmissionError } from 'redux-form';
import { FormValidationMessage } from 'react-native-elements';
import { FormInput } from 'react-native-elements';

import {
  ActionButton,
  FormP,
  FormProps,
  LabeledFormInput,
} from '../components';
import profileService from '../services/profile-service';

interface SignupFormData {
  firstName: string;
  lastName: string;
  email: string;
  phoneNumber: string;
  password: string;
}

const SignupForm: React.SFC<FormProps<SignupFormData>> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting } = props;
  const onSubmitWithReset = async (values: SignupFormData): Promise<void> => {
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
      />
      <Field
        label="Last name"
        name="lastName"
        component={LabeledFormInput}
        autoCorrect={false}
      />
      <Field
        label="Email"
        name="email"
        component={LabeledFormInput}
        keyboardType={'email-address' as 'email-address'}
        autoCorrect={false}
        autoCapitalize={'none' as 'none'}
      />
      <Field
        label="Phone number"
        name="phoneNumber"
        component={LabeledFormInput}
        keyboardType={'phone-pad' as 'phone-pad'}
      />
      <Field
        label="Password"
        name="password"
        component={LabeledFormInput}
        secureTextEntry={true}
      />
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
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
    const userId = await profileService.signup({
      ...values,
      gender: 'male',
      birthday: 847324800,
    });
  }

  render() {
    return <SignupFormWithRedux onSubmit={this.onSubmit} />;
  }
}
