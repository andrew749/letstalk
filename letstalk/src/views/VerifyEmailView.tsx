import React, {Component} from 'react';
import { NavigationScreenProp, NavigationStackAction, NavigationActions } from 'react-navigation';
import { View, StyleSheet } from 'react-native';
import { reduxForm, Field, SubmissionError } from 'redux-form';
import { FormValidationMessage, Text } from 'react-native-elements';
import { ActionButton, FormP, FormProps, LabeledFormInput } from '../components';
import { AnalyticsHelper } from '../services';
import { infoToast, errorToast } from '../redux/toast';
import { headerStyle, headerTitleStyle, headerTintColor } from './TopHeader';
import Colors from '../services/colors';
import auth from '../services/auth';
import {ActionCreator, connect, Dispatch} from 'react-redux';
import { RootState } from '../redux';

import { required, uwEmail } from '../validators';
import {ThunkAction} from "redux-thunk";
import { ActionTypes as BootstrapActionTypes } from '../redux/bootstrap/actions';
import {
  State as BootstrapState,
  fetchBootstrap,
} from '../redux/bootstrap/reducer';

interface Props extends BootstrapState {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  fetchBootstrap: ActionCreator<ThunkAction<Promise<BootstrapActionTypes>, BootstrapState, void>>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

interface State {
  refreshing: boolean;
}

interface VerifyEmailFormData {
  email: string,
}

const VerifyEmailForm: React.SFC<FormProps<VerifyEmailFormData>> = props => {
  const { error, handleSubmit, onSubmit, submitting, valid, pristine } = props;
  const onSubmitWithReset = async (values: VerifyEmailFormData): Promise<void> => {
    await onSubmit(values);
  };
  return (
    <View>
      <Field
        label="UW Email"
        name="email"
        keyboardType={'email-address' as 'email-address'}
        component={LabeledFormInput}
        autoCorrect={false}
        autoCapitalize={'none' as 'none'}
        validate={[required, uwEmail]}
      />
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
        buttonStyle={{backgroundColor: Colors.HIVE_PRIMARY}}
        textStyle={{color: Colors.HIVE_MAIN_FONT}}
        disabled={pristine || !valid}
        loading={submitting}
        title={submitting ? null : "Send verification email"}
        onPress={handleSubmit(onSubmitWithReset)}
      />
    </View>
  );
}

const VerifyEmailFormWithRedux = reduxForm<VerifyEmailFormData, FormP<VerifyEmailFormData>>({
  form: 'sendVerificationEmail',
})(VerifyEmailForm);

export class VerifyEmailView extends Component<Props, State> {
  VERIFY_EMAIL_VIEW_IDENTIFIER = "VerifyEmailView";

  static navigationOptions = {
    headerTitle: 'Verify Email',
    headerStyle, 
    headerTitleStyle, 
    headerTintColor
  }

  constructor(props: Props) {
    super(props);
    this.onSubmit = this.onSubmit.bind(this);
    this.state = { refreshing: false };
  }

  async componentDidMount() {
    AnalyticsHelper.getInstance().recordPage(this.VERIFY_EMAIL_VIEW_IDENTIFIER);

    this.load();
  }

  private async load() {
    await Promise.all([
      this.props.fetchBootstrap(),
    ]);
  }

  async componentWillReceiveProps(nextProps: Props) {
    if (nextProps.bootstrap && nextProps.bootstrap.state === 'account_email_verified') {
      // Email has been verified, move on to onboarding screen.
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        actions: [NavigationActions.navigate({ routeName: 'Onboarding' })]
      }));
    }
  }

  async onSubmit(values: VerifyEmailFormData) {
    try {
      await auth.sendVerificationEmail(values.email);
      await this.props.infoToast("Verification email sent");
    } catch(e) {
      await this.props.errorToast(e.errorMsg);
      throw new SubmissionError({_error: e.errorMsg});
    }
}

  render() {
    return (
      <View style={styles.container}>
        <Text style={styles.headline}>Please provide your UW email address to verify you are a Waterloo student.</Text>
        <VerifyEmailFormWithRedux onSubmit={this.onSubmit} />
        <ActionButton
          buttonStyle={{backgroundColor: Colors.HIVE_PRIMARY}}
          textStyle={{color: Colors.HIVE_MAIN_FONT}}
          title="Refresh"
          onPress={() => this.load()}
        />
      </View>
    );
  }
}

export default connect(({ bootstrap }: RootState) => bootstrap,
  { infoToast, errorToast, fetchBootstrap })(VerifyEmailView);

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  centeredContainer: {
    flex: 1,
    justifyContent: 'center'
  },
  headline: {
    fontWeight: 'bold',
    fontSize: 14,
    textAlign: 'left',
    margin: 20,
  },
});
