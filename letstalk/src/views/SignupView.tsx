import {FileSystem} from 'expo';
import React, { Component } from 'react';
import { Dimensions, Picker, ScrollView, StyleSheet, View } from 'react-native';
import {
  NavigationScreenProp,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import { reduxForm, Field, InjectedFormProps, SubmissionError } from 'redux-form';
import {FormValidationMessage, FormInputProps, FormInput, Text} from 'react-native-elements';
import { KeyboardAwareScrollView } from 'react-native-keyboard-aware-scroll-view'
import { GenderId } from '../models/user';

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
import { headerStyle, headerTitleStyle, headerTintColor } from './TopHeader';
import auth from "../services/auth";
import { required, email, phoneNumber, password } from '../validators';

const SCREEN_WIDTH = Dimensions.get('window').width;

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

class SignupFormRefs {
  firstNameFieldRef: Field<FormInputProps>;
  lastNameFieldRef: Field<FormInputProps>;
  emailFieldRef: Field<FormInputProps>;
  phoneNumberFieldRef: Field<FormInputProps>;
  passwordFieldRef: Field<FormInputProps>;
}

const SignupForm: React.SFC<FormProps<SignupFormData>> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting, valid } = props;
  const onSubmitWithReset = async (values: SignupFormData): Promise<void> => {
    await onSubmit(values);
    reset();
  };
  const fieldRefs = new SignupFormRefs();
  return (
    <KeyboardAwareScrollView
      keyboardShouldPersistTaps="always"
    >
      <View style={styles.profilePicContainer}>
        <Field
          name="profilePic"
          component={ProfileAvatarEditableFormElement}
        />
      </View>
      <Field
        label="First name"
        name="firstName"
        component={LabeledFormInput}
        ref={(ref: Field<FormInputProps>) => fieldRefs.firstNameFieldRef = ref}
        onSubmitEditing={() => {
          // @ts-ignore
          fieldRefs.lastNameFieldRef.getRenderedComponent().focus();
        }}
        withRef={true}
        autoCorrect={false}
        validate={required}
      />
      <Field
        label="Last name"
        name="lastName"
        component={LabeledFormInput}
        ref={(ref: Field<FormInputProps>) => fieldRefs.lastNameFieldRef = ref}
        onSubmitEditing={() => {
          // @ts-ignore
          fieldRefs.emailFieldRef.getRenderedComponent().focus();
        }}
        withRef={true}
        autoCorrect={false}
        validate={required}
      />
      <Field
        label="Email"
        name="email"
        component={LabeledFormInput}
        ref={(ref: Field<FormInputProps>) => fieldRefs.emailFieldRef = ref}
        onSubmitEditing={() => {
          // @ts-ignore
          fieldRefs.phoneNumberFieldRef.getRenderedComponent().focus();
        }}
        withRef={true}
        keyboardType={'email-address' as 'email-address'}
        autoCorrect={false}
        autoCapitalize={'none' as 'none'}
        validate={[required, email]}
      />
      <Field
        label="Phone number"
        ref={(ref: Field<FormInputProps>) => fieldRefs.phoneNumberFieldRef = ref}
        name="phoneNumber"
        onSubmitEditing={() => {
          // @ts-ignore
          fieldRefs.passwordFieldRef.getRenderedComponent().focus();
        }}
        withRef={true}
        component={LabeledFormInput}
        keyboardType={'phone-pad' as 'phone-pad'}
        validate={[required, phoneNumber]}
      />
      <Text style={styles.formTipText}>Your primary email and phone number will be displayed on your profile to your matches.</Text>
      <Field
        label="Password"
        name="password"
        component={LabeledFormInput}
        ref={(ref: Field<FormInputProps>) => fieldRefs.passwordFieldRef = ref}
        withRef={true}
        secureTextEntry={true}
        validate={[required, password]} // Add some rules for password
        autoCapitalize={'none' as 'none'}
        containerStyle={{marginBottom: 10}}
      />
      <Field
        label="Gender (optional)"
        name="gender"
        component={ButtonPicker}
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
        label="Birthday (optional)"
        name="birthdate"
        androidMode={'spinner' as 'spinner'}
        mode={'date' as 'date'}
        component={ModalDatePicker}
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
    </KeyboardAwareScrollView>
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
    headerStyle,
    headerTitleStyle,
    headerTintColor
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
          "gender": values.gender ? values.gender : GenderId.Unspecified,
          "password": values.password,
          "profilePic": values.profilePic ? values.profilePic.data: undefined,
          "birthdate": values.birthdate,
        }
      });
    } catch(e) {
      throw new SubmissionError({_error: e.errorMsg});
    }
    // TODO: have a prompt saying successfully signed up
    // Immediately log in.
    try {
      let token: string = null;
      // Don't fail if expo is down
      try {
        token = await auth.registerForPushNotificationsAsync();
      } catch(e){
        console.log("Failed to register for notifications");
      }
      await auth.login(values.email, values.password, token);
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        actions: [NavigationActions.navigate({ routeName: 'SplashScreenView' })]
      }));
    } catch(e) {
      throw new SubmissionError({_error: e.errorMsg});
    }
  }

  render() {
    return (
      <SignupFormWithRedux onSubmit={this.onSubmit} />
    );
  }
}

// Hack to make submit button available when keyboard is open.
const styles = StyleSheet.create({
  profilePicContainer: {
    justifyContent: 'center',
    alignItems: 'center',
  },
  submitButton: {
    marginBottom: 100,
  },
  formTipText: {
    color: Colors.HIVE_SUBDUED,
    fontSize: 14,
    paddingLeft: 20,
    paddingRight: 20,
  },
});
