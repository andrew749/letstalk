import React, { Component } from 'react';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import { Picker, ScrollView } from 'react-native';
import {
  NavigationScreenProp,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import { reduxForm, Field, InjectedFormProps, SubmissionError } from 'redux-form';
import { FormValidationMessage } from 'react-native-elements';

import { RootState } from '../redux';
import {
  ActionButton,
  FormP,
  FormProps,
  ModalPicker,
} from '../components';
import profileService from '../services/profile-service';
import { State as BootstrapState, fetchBootstrap } from '../redux/bootstrap/reducer';
import { ActionTypes } from '../redux/bootstrap/actions';

interface OnboardingFormData {
  cohortId: number,
}

// TODO: move elsewhere
const required = (value: any) => (value ? undefined : 'Required')

const OnboardingForm: React.SFC<FormProps<OnboardingFormData>> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting, valid } = props;
  const onSubmitWithReset = async (values: OnboardingFormData): Promise<void> => {
    await onSubmit(values);
    reset();
  };
  // TODO: Maybe pull these from the server and split the form up into parts (program, year)
  return (
    <ScrollView>
      <Field
        label="Cohort"
        name="cohortId"
        component={ModalPicker}
        validate={required}
      >
        <Picker.Item
          label="Software Engineering 2019"
          value={1}
        />
        <Picker.Item
          label="Computer Engineering Stream 8 2019"
          value={2}
        />
        <Picker.Item
          label="Computer Engineering Stream 4 2019"
          value={3}
        />
      </Field>
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
        disabled={!valid}
        loading={submitting}
        title={submitting ? null : "Finish Onboarding"}
        onPress={handleSubmit(onSubmitWithReset)}
      />
    </ScrollView>
  );
}

const OnboardingFormWithRedux = reduxForm<OnboardingFormData, FormP<OnboardingFormData>>({
  form: 'onboarding',
})(OnboardingForm);

interface DispatchActions {
  fetchBootstrap: ActionCreator<ThunkAction<Promise<ActionTypes>, BootstrapState, void>>;
}

interface Props extends DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class OnboardingView extends Component<Props> {
  static navigationOptions = {
    headerTitle: 'Onboarding',
  }

  constructor(props: Props) {
    super(props);

    this.onSubmit = this.onSubmit.bind(this);
  }

  async onSubmit(values: OnboardingFormData) {
    try {
      await profileService.updateCohort(values);
      // Reload bootstrap data after updating
      await this.props.fetchBootstrap();
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        actions: [NavigationActions.navigate({ routeName: 'Home' })]
      }));
    } catch(e) {
      throw new SubmissionError({_error: e.message});
    }
  }

  render() {
    return <OnboardingFormWithRedux onSubmit={this.onSubmit} />;
  }
}

export default connect(null, { fetchBootstrap })(OnboardingView);
