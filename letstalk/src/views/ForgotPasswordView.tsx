import React, {Component} from 'react';
import { NavigationScreenProp, NavigationStackAction, StackActions, NavigationActions } from 'react-navigation';
import { StyleSheet } from 'react-native';
import { reduxForm, Field, SubmissionError } from 'redux-form';
import { FormValidationMessage } from 'react-native-elements';
import { ActionButton, FormP, FormProps, LabeledFormInput } from '../components';
import { AnalyticsHelper } from '../services';
import { infoToast, errorToast } from '../redux/toast';
import { headerStyle, headerTitleStyle, headerTintColor } from './TopHeader';
import Colors from '../services/colors';
import auth from '../services/auth';
import { connect, Dispatch } from 'react-redux';
import { RootState } from '../redux';
import { KeyboardAwareScrollView } from 'react-native-keyboard-aware-scroll-view'

import { required, email } from '../validators';

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

interface ForgotPasswordFormData {
  email: string,
}

const ForgotPasswordForm: React.SFC<FormProps<ForgotPasswordFormData>> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting, valid, pristine } = props;
  const onSubmitWithReset = async (values: ForgotPasswordFormData): Promise<void> => {
    await onSubmit(values);
  };
  return (
    <KeyboardAwareScrollView>
      <Field
        label="Email"
        name="email"
        keyboardType={'email-address' as 'email-address'}
        component={LabeledFormInput as "input" & typeof LabeledFormInput}
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
    </KeyboardAwareScrollView>
  );
}

const ForgotPasswordFormWithRedux = reduxForm<ForgotPasswordFormData, FormP<ForgotPasswordFormData>>({
  form: 'forgotPassword',
})(ForgotPasswordForm as any);

export class ForgotPasswordView extends Component<Props> {
  FORGOT_PASSWORD_VIEW_IDENTIFIER = "ForgotPasswordView";

  static navigationOptions = {
    headerTitle: 'Forgot Password',
    headerStyle, 
    headerTitleStyle, 
    headerTintColor
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
      await this.props.infoToast("Check your email for reset instructions!");
      this.props.navigation.dispatch(StackActions.reset({
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
      <ForgotPasswordFormWithRedux onSubmit={this.onSubmit} />
    );
  }
}

export default connect(null, {infoToast, errorToast})(ForgotPasswordView as any);

const styles = StyleSheet.create({
  container: {
    flex: 1,
  }
});
