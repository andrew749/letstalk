import React, { Component } from 'react';
import { Button as ReactNativeButton, Dimensions, View } from 'react-native';
import { FormInput, FormLabel, FormValidationMessage } from 'react-native-elements';

import { ActionButton } from '../components';
import auth from '../services/auth';
import { InvalidCredentialsError } from '../services/sessionService';

interface Props {
}

interface State {
  username: string
  password: string
  error: string | null
}

export default class LoginPage extends Component<Props, State> {
  usernameInput: FormInput
  passwordInput: FormInput

  static navigationOptions = {
    headerRight: <ReactNativeButton title="Sign up" onPress={() => null} />,
    headerTitle: 'Log in',
  }

  constructor(props: Props) {
    super(props);
    this.state = {
      username: '',
      password: '',
      error: null,
    };
    this.usernameInput = null;
    this.passwordInput = null;

    this.onLoginPress = this.onLoginPress.bind(this);
    this.onUsernameChange = this.onUsernameChange.bind(this);
    this.onPasswordChange = this.onPasswordChange.bind(this);
  }

  onUsernameChange(username: string) {
    this.setState({ username });
  }

  onPasswordChange(password: string) {
    this.setState({ password });
  }

  async onLoginPress() {
    const {
      username,
      password,
    } = this.state;
    try {
      await auth.login(username, password)
    } catch(e) {
      if (e.message !== InvalidCredentialsError.tpe) throw e;
      this.setState({ error: 'Invalid username or password' });
      this.usernameInput.shake();
      this.passwordInput.shake();
      this.passwordInput.focus();
    }
  }

  render() {
    // TODO: Add header
    let validationMessage = this.state.error || ' ';
    return (
      <View>
        <FormLabel>Username</FormLabel>
        <FormInput
          ref={input => this.usernameInput = input}
          autoCorrect={false}
          onChangeText={this.onUsernameChange}
        />
        <FormLabel>Password</FormLabel>
        <FormInput
          ref={input => this.passwordInput = input}
          secureTextEntry={true}
          onChangeText={this.onPasswordChange}
        />
        <FormValidationMessage>
          {validationMessage}
        </FormValidationMessage>
        <ActionButton title="Log in" onPress={this.onLoginPress} />
      </View>
    );
  }
}
