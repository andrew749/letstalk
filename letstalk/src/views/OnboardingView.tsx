// TODO: maybe rename to OnboardingView
import React, { Component } from 'react';
import { Picker, ScrollView } from 'react-native';
import { reduxForm, Field, InjectedFormProps, SubmissionError } from 'redux-form';
import { FormValidationMessage } from 'react-native-elements';

import {
  ActionButton,
  FormP,
  FormProps,
  ModalPicker,
} from '../components';
import profileService from '../services/profile-service';

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

interface Props { }

export default class OnboardingView extends Component<Props> {
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
    } catch(e) {
      throw new SubmissionError({_error: e.message});
    }
  }

  render() {
    return <OnboardingFormWithRedux onSubmit={this.onSubmit} />;
  }
}
