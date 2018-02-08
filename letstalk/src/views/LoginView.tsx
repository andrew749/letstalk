import React, { Component } from 'react';
import { Button as ReactNativeButton, View } from 'react-native';
import { FormValidationMessage } from 'react-native-elements';
import { reduxForm, Field, InjectedFormProps, SubmissionError } from 'redux-form';
import {
  NavigationScreenProp,
  NavigationStackAction,
  NavigationActions,
  NavigationScreenDetails,
} from 'react-navigation';

import { ActionButton, LabeledFormInput } from '../components';
import { InvalidCredentialsError } from '../services/sessionService';
import auth from '../services/auth';

interface LoginFormData {
  username: string;
  password: string;
}

interface LoginFormProps {
  onSubmit(values: LoginFormData): void;
}

const LoginForm: React.SFC<
    LoginFormProps & InjectedFormProps<LoginFormData, LoginFormProps>> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting } = props;
  const onSubmitWithReset = async (values: LoginFormData): Promise<void> => {
    await onSubmit(values);
    reset();
  };
  return (
    <View>
      <Field
        label="Username"
        name="username"
        component={LabeledFormInput}
        autoCorrect={false}
        autoCapitalize={'none' as 'none'}
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
        title={submitting ? null : "Log in"}
        onPress={handleSubmit(onSubmitWithReset)}
      />
    </View>
  );
}

const LoginFormWithRedux = reduxForm<LoginFormData, LoginFormProps>({
  form: 'signup',
})(LoginForm);

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

export default class LoginView extends Component<Props> {
  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: 'Log in',
    headerRight: <ReactNativeButton title="Sign up"
      onPress={() => navigation.navigate('Signup')} />,
  })

  constructor(props: Props) {
    super(props);

    this.onSubmit = this.onSubmit.bind(this);
  }

  async onSubmit(values: LoginFormData) {
    const {
      username,
      password,
    } = values;
    console.log(values);
    try {
      await auth.login(username, password);
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        key: null,
        actions: [NavigationActions.navigate({ routeName: 'Main' })]
      }));
    } catch(e) {
      throw new SubmissionError({_error: e.message});
    }
  }

  render() {
    return <LoginFormWithRedux onSubmit={this.onSubmit} />;
  }
}
