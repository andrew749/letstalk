import React, { Component } from 'react';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import { Picker, ScrollView, Text } from 'react-native';
import {
  NavigationScreenProp,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import {
  reduxForm,
  formValueSelector,
  Field,
  InjectedFormProps,
  SubmissionError,
} from 'redux-form';
import { FormValidationMessage, Rating } from 'react-native-elements';
import Immutable from 'immutable';

import { RootState } from '../redux';
import { State as OnboardingState } from '../redux/onboarding/reducer';
import { setStep, SetStepAction, Step } from '../redux/onboarding/actions';
import {
  ActionButton,
  Emoji,
  FormP,
  FormProps,
  Header,
  InfoText,
  ModalPicker,
} from '../components';
import profileService from '../services/profile-service';
import { State as BootstrapState, fetchBootstrap } from '../redux/bootstrap/reducer';
import { ActionTypes } from '../redux/bootstrap/actions';
import {
  getCohortId,
  programOptions,
  sequenceOptions,
  gradYearOptions,
  ValueLabel,
} from '../models/cohort';

interface OnboardingFormData {
  programId: string,
  sequenceId: string;
  gradYear: string;
}

// TODO: move elsewhere
const required = (value: any) => (value ? undefined : 'Required')

const OnboardingForm: React.SFC<FormProps<OnboardingFormData> & OnboardingFormData> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting, valid, programId, sequenceId } = props;
  const onSubmitWithReset = async (values: OnboardingFormData): Promise<void> => {
    await onSubmit(values);
    reset();
  };
  const buildItems = (rows: Immutable.List<ValueLabel>) => {
    return rows.map(({ value, label }) => {
      return <Picker.Item key={value} label={label} value={value}/>;
    });
  };
  const programItems = buildItems(programOptions()).toJS();
  const sequenceItems = buildItems(sequenceOptions(programId)).toJS();
  const gradYearItems = buildItems(gradYearOptions(programId, sequenceId)).toJS();
  return (
    <ScrollView>
      <Field
        label="Program"
        name="programId"
        component={ModalPicker}
        validate={required}
      >
        {programItems}
      </Field>
      <Field
        label="Sequence"
        name="sequenceId"
        component={ModalPicker}
        validate={required}
      >
        {sequenceItems}
      </Field>
      <Field
        label="Grad Year"
        name="gradYear"
        component={ModalPicker}
        validate={required}
      >
        {gradYearItems}
      </Field>
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
        disabled={!valid}
        loading={submitting}
        title={submitting ? null : "Choose cohort"}
        onPress={handleSubmit(onSubmitWithReset)}
      />
    </ScrollView>
  );
}

const cohortSelector = formValueSelector('onboarding-cohort');

const OnboardingFormWithRedux = reduxForm<OnboardingFormData, FormP<OnboardingFormData>>({
  form: 'onboarding-cohort',
})(connect(state => ({
  programId: cohortSelector(state, 'programId'),
  sequenceId: cohortSelector(state, 'sequenceId'),
}))(OnboardingForm));

interface DispatchActions {
  fetchBootstrap: ActionCreator<ThunkAction<Promise<ActionTypes>, BootstrapState, void>>;
  setStep(step: Step): SetStepAction;
}

interface Props extends OnboardingState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class OnboardingView extends Component<Props> {
  static navigationOptions = {
    headerTitle: 'Onboarding',
  }

  constructor(props: Props) {
    super(props);

    this.onSubmit = this.onSubmit.bind(this);
    this.onSubmitCohort = this.onSubmitCohort.bind(this);
  }

  async onSubmitCohort(values: OnboardingFormData) {
    const { programId, sequenceId, gradYear } = values;
    const cohortId = getCohortId(programId, sequenceId, gradYear);
    try {
      // await profileService.updateCohort({ cohortId });
    } catch(e) {
      throw new SubmissionError({_error: e.message});
    }
    this.props.setStep(Step.MY_VECTOR);
  }

  async onSubmit(values: OnboardingFormData) {
    const { programId, sequenceId, gradYear } = values;
    const cohortId = getCohortId(programId, sequenceId, gradYear);
    try {
      await profileService.updateCohort({ cohortId });

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
    const { step } = this.props;
    switch (step) {
      case Step.COHORT:
        // TODO: Gender specific emoji
        return (
          <ScrollView>
            <Header>Your Cohort</Header>
            <InfoText>
              Based on your cohort, you will either be a big <Emoji name="man"/>, mentoring other
              students, or a small <Emoji name="baby"/>, being mentored by an upper year student.
            </InfoText>
            <OnboardingFormWithRedux onSubmit={this.onSubmitCohort} />
          </ScrollView>
        );
      case Step.MY_VECTOR:
        return (
          <ScrollView>
            <Header>Your Personality</Header>
            <InfoText>
              Looks like you're going to be a big <Emoji name="man"/>! We'd like to get
              to know you a little better, so that we can find you some lit <Emoji name="fire"/>
              noobies to mentor.
            </InfoText>
          </ScrollView>
        );
      case Step.YOUR_VECTOR:
        return <Text>your vector</Text>;
      default:
        // Ensure exhaustiveness of select
        const _: never = step;
    }
  }
}

export default connect(({ onboarding }: RootState) => onboarding,
  { fetchBootstrap, setStep }
)(OnboardingView);
