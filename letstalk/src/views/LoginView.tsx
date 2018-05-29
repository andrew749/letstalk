import React, { Component } from 'react';
import { Button as ReactNativeButton, View, StyleSheet, Platform } from 'react-native';
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
import Colors from '../services/colors';
import {AnalyticsHelper} from '../services/analytics';


interface LoginFormData {
  email: string;
  password: string;
}

// TODO: move elsewhere
const required = (value: any) => (value ? undefined : 'Required')
const email = (value: string) =>
  value && !/^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$/i.test(value)
    ? 'Invalid email address'
    : undefined

const LoginForm: React.SFC<FormProps<LoginFormData>> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting, valid } = props;
  const onSubmitWithReset = async (values: LoginFormData): Promise<void> => {
    await onSubmit(values);
    reset();
  };
  return (
    <View>
      <Field
        label="Email"
        name="email"
        keyboardType={'email-address' as 'email-address'}
        component={LabeledFormInput}
        autoCorrect={false}
        autoCapitalize={'none' as 'none'}
        validate={[required, email]}
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
        buttonStyle={{backgroundColor: Colors.HIVE_MAIN_BG}}
        textStyle={{color: Colors.HIVE_MAIN_FONT}}
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

interface FBLoginFormData {}

const FBLoginForm: React.SFC<FormProps<FBLoginFormData>> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting, valid } = props;
  const onSubmitWithReset = async (values: FBLoginFormData): Promise<void> => {
    await onSubmit(values);
    reset();
  };
  return (
    <View>
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
        backgroundColor={Colors.FB_BG}
        loading={submitting}
        title={submitting ? null : "Log in with Facebook"}
        onPress={handleSubmit(onSubmitWithReset)}
      />
    </View>
  );
}

const FBLoginFormWithRedux = reduxForm<LoginFormData, FormP<LoginFormData>>({
  form: 'fblogin',
})(FBLoginForm);

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

export default class LoginView extends Component<Props> {
  LOGIN_VIEW_IDENTIFIER = "LoginView";

  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: 'Log in',
    headerRight: Platform.OS === 'ios' ? <ReactNativeButton title="Sign up"
      onPress={() => navigation.navigate('Signup')} /> : null,
    headerStyle: {
      backgroundColor: Colors.HIVE_MAIN_BG,
    },
  })

  constructor(props: Props) {
    super(props);
    this.onSubmit = this.onSubmit.bind(this);
    this.onSubmitFb = this.onSubmitFb.bind(this);
  }

  async componentDidMount() {
    AnalyticsHelper.getInstance().recordPage(this.LOGIN_VIEW_IDENTIFIER);
  }

  async registerForPushNotificationsAsync(): Promise<string> {
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
    console.log("Registered with expo notification service: " + token);
    return token;
  }

  async onSubmitFb() {
    try {
      let token: string = null;
      // don't fail if expo is down
      try {
        token = await this.registerForPushNotificationsAsync();
      } catch(e){
        console.log("Failed to register for notification")
      }
      if (await auth.loginWithFb(token)) {
        this.props.navigation.dispatch(NavigationActions.reset({
          index: 0,
          actions: [NavigationActions.navigate({ routeName: 'Tabbed' })]
        }));
      }
    } catch(e) {
      throw new SubmissionError({_error: e.message});
    }
  }

  async onSubmit(values: LoginFormData) {
    AnalyticsHelper.getInstance().recordAction(this.LOGIN_VIEW_IDENTIFIER, "login", "", 1);
    const {
      email,
      password,
    } = values;
    try {
      let token: string = null;
      // don't fail if expo is down
      try {
        token = await this.registerForPushNotificationsAsync();
      } catch(e){
        console.log("Failed to register for notification")
      }
      await auth.login(email, password, token);
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        actions: [NavigationActions.navigate({ routeName: 'Tabbed' })]
      }));
    } catch(e) {
      throw new SubmissionError({_error: e.message});
    }
  }

  render() {
    const signupButton = Platform.OS === 'android' ?
      <ActionButton
        title="Sign up"
        buttonStyle={{backgroundColor: Colors.HIVE_MAIN_BG}}
        textStyle={{color: Colors.HIVE_MAIN_FONT}}
        onPress={() => this.props.navigation.dispatch(NavigationActions.navigate({routeName: 'Signup'}))} />
      : null;

    return (
      <View>
        <LoginFormWithRedux onSubmit={this.onSubmit} />
        {signupButton}
        <FBLoginFormWithRedux onSubmit={this.onSubmitFb} />
      </View>
    );
  }
}
