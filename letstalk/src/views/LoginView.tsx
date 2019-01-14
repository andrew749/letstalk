import React, { Component } from 'react';
import { Button as ReactNativeButton, View, StyleSheet, Platform, TouchableOpacity } from 'react-native';
import { FormValidationMessage, Text } from 'react-native-elements';
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
import { headerStyle, headerTitleStyle, headerTintColor } from './TopHeader';
import { required, email } from '../validators';

const styles = StyleSheet.create({
  forgotPasswordButton: {
    marginLeft: 20,
    marginTop: 2,
  },
  forgotPasswordText: {
    color: Colors.HIVE_PRIMARY,
    fontSize: 14,
  },
});

interface LoginFormData {
  email: string;
  password: string;
}

interface NavProps {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

const LoginForm: React.SFC<FormProps<LoginFormData, NavProps>> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting, valid, pristine } = props;
  const onSubmitWithReset = async (values: LoginFormData): Promise<void> => {
    await onSubmit(values);
    reset();
  };

  const forgotPasswordButton = (
    <TouchableOpacity
      style={styles.forgotPasswordButton}
      onPress={() => {
        props.navigation.dispatch(NavigationActions.navigate({routeName: 'ForgotPassword'}))
      }}
    >
      <Text style={styles.forgotPasswordText}>Forgot my password</Text>
    </TouchableOpacity>
  );

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
        autoCapitalize={'none' as 'none'}
      />
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      {forgotPasswordButton}
      <ActionButton
        buttonStyle={{backgroundColor: Colors.HIVE_PRIMARY}}
        textStyle={{color: Colors.HIVE_MAIN_FONT}}
        disabled={pristine || !valid}
        loading={submitting}
        title={submitting ? null : "Log in"}
        onPress={handleSubmit(onSubmitWithReset)}
      />
    </View>
  );
}

const LoginFormWithRedux = reduxForm<LoginFormData, FormP<LoginFormData, NavProps>>({
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
      headerStyle,
      headerTitleStyle,
      headerTintColor
  })

  constructor(props: Props) {
    super(props);
    this.onSubmit = this.onSubmit.bind(this);
    this.onSubmitFb = this.onSubmitFb.bind(this);
  }

  async componentDidMount() {
    AnalyticsHelper.getInstance().recordPage(this.LOGIN_VIEW_IDENTIFIER);
  }

  async onSubmitFb() {
    try {
      let token: string = null;
      // don't fail if expo is down
      try {
        token = await auth.registerForPushNotificationsAsync();
      } catch(e){
        console.log("Failed to register for notification " + e);
      }
      console.log("Logging in with Facebook")
      if (await auth.loginWithFb(token)) {
        this.props.navigation.dispatch(NavigationActions.reset({
          index: 0,
          actions: [NavigationActions.navigate({ routeName: 'SplashScreen' })]
        }));
      }
    } catch(e) {
      console.log("Unable to login: " + e)
      throw new SubmissionError({_error: e.errorMsg});
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
        token = await auth.registerForPushNotificationsAsync();
      } catch(e){
        console.log("Failed to register for notification " + e);
      }
      await auth.login(email, password, token);
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        actions: [NavigationActions.navigate({ routeName: 'SplashScreen' })]
      }));
    } catch(e) {
      console.log("Unable to login: " + e)
      throw new SubmissionError({_error: e.errorMsg});
    }
  }

  render() {
    const signupButton = Platform.OS === 'android' ?
      <ActionButton
        title="Sign up"
        buttonStyle={{backgroundColor: Colors.HIVE_PRIMARY}}
        textStyle={{color: Colors.HIVE_MAIN_FONT}}
        onPress={() => this.props.navigation.dispatch(NavigationActions.navigate({routeName: 'Signup'}))} />
      : null;

    return (
      <View>
        <LoginFormWithRedux onSubmit={this.onSubmit} navigation={this.props.navigation} />
        {signupButton}
        <FBLoginFormWithRedux onSubmit={this.onSubmitFb} />
      </View>
    );
  }
}
