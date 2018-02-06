import React, { Component } from 'react';
import { Button as ReactNativeButton, Dimensions, View } from 'react-native';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import { FormInput, FormLabel, FormValidationMessage } from 'react-native-elements';
import { ThunkAction } from 'redux-thunk';
import { NavigationScreenProp, NavigationStackAction, NavigationActions } from 'react-navigation';

import { ActionButton } from '../components';
import { RootState } from '../redux';
import { login, State as LoginState } from '../redux/login/reducer';
import {
  resetAction,
  ResetAction,
  setPasswordAction,
  setUsernameAction,
  SetPasswordAction,
  SetUsernameAction,
} from '../redux/login/actions';
import { InvalidCredentialsError } from '../services/sessionService';

interface DispatchActions {
  login: ActionCreator<ThunkAction<Promise<void>, LoginState, void>>;
  setUsername: ActionCreator<SetUsernameAction>;
  setPassword: ActionCreator<SetPasswordAction>;
  reset: ActionCreator<ResetAction>;
};

interface Props extends LoginState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class LoginView extends Component<Props> {
  usernameInput: FormInput
  passwordInput: FormInput

  static navigationOptions = {
    headerRight: <ReactNativeButton title="Sign up" onPress={() => null} />,
    headerTitle: 'Log in',
  }

  constructor(props: Props) {
    super(props);
    this.usernameInput = null;
    this.passwordInput = null;

    this.onLoginPress = this.onLoginPress.bind(this);
    this.onUsernameChange = this.onUsernameChange.bind(this);
    this.onPasswordChange = this.onPasswordChange.bind(this);
  }

  onUsernameChange(username: string) {
    this.props.setUsername(username);
  }

  onPasswordChange(password: string) {
    this.props.setPassword(password);
  }

  async onLoginPress() {
    const {
      username,
      password,
    } = this.props;
    try {
      await this.props.login(username, password);
      this.props.reset();
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        key: null,
        actions: [NavigationActions.navigate({ routeName: 'Main' })]
      }));
    } catch(e) {
      console.log(e.message);
      if (e.message !== InvalidCredentialsError.tpe) throw e;
      this.usernameInput.shake();
      this.passwordInput.shake();
    }
  }

  render() {
    // TODO: Add header
    let validationMessage = this.props.errorMsg || ' ';
    return (
      <View>
        <FormLabel>Username</FormLabel>
        <FormInput
          ref={input => this.usernameInput = input}
          autoCorrect={false}
          autoCapitalize={'none'}
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
        <ActionButton
          loading={this.props.isFetching}
          title={this.props.isFetching ? null : "Log in"}
          onPress={this.onLoginPress}
        />
      </View>
    );
  }
}

export default connect(({ login }: RootState) => login, {
  login,
  setUsername: setUsernameAction,
  setPassword: setPasswordAction,
  reset: resetAction,
})(LoginView);
