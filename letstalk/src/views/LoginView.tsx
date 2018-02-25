import React, { Component } from 'react';
import { Button as ReactNativeButton, View } from 'react-native';
import { FormValidationMessage } from 'react-native-elements';
import { reduxForm, Field, SubmissionError } from 'redux-form';
import {
  NavigationScreenProp,
  NavigationStackAction,
  NavigationActions,
  NavigationScreenDetails,
} from 'react-navigation';
import { Permissions, Notifications } from 'expo';

import {
  ActionButton,
  FormP,
  FormProps,
  LabeledFormInput
} from '../components';
import { InvalidCredentialsError } from '../services/session-service';
import auth from '../services/auth';

interface LoginFormData {
  username: string;
  password: string;
}

// TODO: move elsewhere
const required = (value: any) => (value ? undefined : 'Required')

const LoginForm: React.SFC<FormProps<LoginFormData>> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting, valid } = props;
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
        validate={required}
      />
      <Field
        label="Password"
        name="password"
        component={LabeledFormInput}
        secureTextEntry={true}
        validate={required}
      />
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
        disabled={!valid}
        loading={submitting}
        title={submitting ? null : "Log in"}
        onPress={handleSubmit(onSubmitWithReset)}
      />
    </View>
  );
}

const LoginFormWithRedux = reduxForm<LoginFormData, FormP<LoginFormData>>({
  // TODO: Enum with these to make sure there are no conflicts
  form: 'login',
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

  async registerForPushNotificationsAsync() {
    const { status: existingStatus } = await Permissions.getAsync(Permissions.NOTIFICATIONS as any);
    let finalStatus = existingStatus;

    // only ask if permissions have not already been determined, because
    // iOS won't necessarily prompt the user a second time.
    if (existingStatus !== 'granted') {
      // Android remote notification permissions are granted during the app
      // install, so this will only ask on iOS
      const { status } = await Permissions.askAsync(Permissions.NOTIFICATIONS as any);
      finalStatus = status;
    }

    // Stop here if the user did not grant permissions
    if (finalStatus !== 'granted') {
      return;
    }

    // Get the token that uniquely identifies this device
    let token = await Notifications.getExpoPushTokenAsync();

    // TODO: POST the token to your backend server from where you can retrieve it to send push
    // notifications.
    console.log(token);
  }

  async onSubmit(values: LoginFormData) {
    const {
      username,
      password,
    } = values;
    try {
      await auth.login(username, password);
      await this.registerForPushNotificationsAsync();
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        actions: [NavigationActions.navigate({ routeName: 'Home' })]
      }));
    } catch(e) {
      throw new SubmissionError({_error: e.message});
    }
  }

  render() {
    return <LoginFormWithRedux onSubmit={this.onSubmit} />;
  }
}
