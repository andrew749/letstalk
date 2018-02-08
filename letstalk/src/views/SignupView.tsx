import React, { Component } from 'react';
import { ScrollView } from 'react-native';
import {
  NavigationScreenProp,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import { FormInput } from 'react-native-elements';

import { LabeledFormInput } from '../components';

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

export default class SignupView extends Component<Props> {
  firstNameInput: FormInput
  lastNameInput: FormInput
  emailInput: FormInput
  phoneNumberInput: FormInput
  passwordInput: FormInput

  static navigationOptions = {
    headerTitle: 'Signup',
  }

  constructor(props: Props) {
    super(props);
    this.firstNameInput = null;
    this.lastNameInput = null;
    this.emailInput = null;
    this.phoneNumberInput = null;
    this.passwordInput = null;
  }

  render() {
    /*return (
      <ScrollView>
        <LabeledFormInput
          label="First name"
          formInputRef={input => this.firstNameInput = input}
          autoCorrect={false}
        />
        <LabeledFormInput
          label="Last name"
          formInputRef={input => this.lastNameInput = input}
          autoCorrect={false}
        />
        <LabeledFormInput
          label="Email"
          formInputRef={input => this.emailInput = input}
          keyboardType="email-address"
          autoCorrect={false}
          autoCapitalize="none"
        />
        <LabeledFormInput
          label="Phone number"
          formInputRef={input => this.phoneNumberInput = input}
          keyboardType="phone-pad"
        />
        <LabeledFormInput
          label="Password"
          formInputRef={input => this.passwordInput = input}
          secureTextEntry={true}
        />
      </ScrollView>
    );
  }*/
    return <ScrollView></ScrollView>;
  }
}
