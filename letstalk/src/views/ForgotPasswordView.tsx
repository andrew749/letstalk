import React, {Component} from 'react';
import { NavigationScreenProp, NavigationStackAction, NavigationActions } from 'react-navigation';
import { View, StyleSheet } from 'react-native';
import { reduxForm, Field, InjectedFormProps, SubmissionError } from 'redux-form';
import { FormValidationMessage, Text } from 'react-native-elements';
import { ActionButton, FormP, FormProps, LabeledFormInput } from '../components';
import { KeyboardAvoidingView } from 'react-native';
import { AnalyticsHelper } from '../services';
import { infoToast, errorToast } from '../redux/toast';
import { headerStyle } from './TopHeader';
import Colors from '../services/colors';
import { ScrollView } from 'react-native';
import auth from '../services/auth';
import { connect, Dispatch } from 'react-redux';
import { RootState } from '../redux';

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

interface ForgotPasswordFormData {
  email: string,
}

const required = (value: any) => (value ? undefined : 'Required')
const email = (value: string) =>
  value && !/^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$/i.test(value)
    ? 'Invalid email address'
    : undefined;
const ForgotPasswordForm: React.SFC<FormProps<ForgotPasswordFormData>> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting, valid, pristine } = props;
  const onSubmitWithReset = async (values: ForgotPasswordFormData): Promise<void> => {
    await onSubmit(values);
  };
  return (
    <ScrollView>
      <Field
        label="Email"
        name="email"
        keyboardType={'email-address' as 'email-address'}
        component={LabeledFormInput}
        autoCorrect={false}
        autoCapitalize={'none' as 'none'}
        validate={[required, email]}
      />
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
        buttonStyle={{backgroundColor: Colors.HIVE_PRIMARY}}
        textStyle={{color: Colors.HIVE_MAIN_FONT}}
        disabled={pristine || !valid}
        loading={submitting}
        title={submitting ? null : "Reset Password"}
        onPress={handleSubmit(onSubmitWithReset)}
      />
    </ScrollView>
  );
}

const ForgotPasswordFormWithRedux = reduxForm<ForgotPasswordFormData, FormP<ForgotPasswordFormData>>({
  form: 'forgotPassword',
})(ForgotPasswordForm);

export class ForgotPasswordView extends Component<Props> {
  FORGOT_PASSWORD_VIEW_IDENTIFIER = "ForgotPasswordView";

  static navigationOptions = {
    headerTitle: 'Forgot Password',
    headerStyle,
  }

  constructor(props: Props) {
    super(props);
    this.onSubmit = this.onSubmit.bind(this);
  }

  async componentDidMount() {
    AnalyticsHelper.getInstance().recordPage(this.FORGOT_PASSWORD_VIEW_IDENTIFIER);
  }

  async onSubmit(values: ForgotPasswordFormData) {
    // TODO submit password reset request
    try {
      await auth.forgotPassword(values.email);
      await this.props.infoToast("Sent an email with reset instructions.");
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        key: null,
        actions: [NavigationActions.navigate({ routeName: 'Login' })]
      }));
    } catch(e) {
      await this.props.errorToast(e.errorMsg);
      throw new SubmissionError({_error: e.errorMsg});
    }
}

  render() {
    return (
      <KeyboardAvoidingView behavior="padding">
        <ForgotPasswordFormWithRedux onSubmit={this.onSubmit} />
      </KeyboardAvoidingView>
    );
  }
}

export default connect(null, {infoToast, errorToast})(ForgotPasswordView);

const styles = StyleSheet.create({
  container: {
    flex: 1,
  }
});
