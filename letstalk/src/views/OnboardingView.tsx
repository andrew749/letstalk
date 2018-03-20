import React, { Component } from 'react';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import { Picker, ScrollView } from 'react-native';
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
import { FormValidationMessage } from 'react-native-elements';
import Immutable from 'immutable';

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
        title={submitting ? null : "Finish Onboarding"}
        onPress={handleSubmit(onSubmitWithReset)}
      />
    </ScrollView>
  );
}

const selector = formValueSelector('onboarding');

const OnboardingFormWithRedux = reduxForm<OnboardingFormData, FormP<OnboardingFormData>>({
  form: 'onboarding',
})(connect(state => ({
  programId: selector(state, 'programId'),
  sequenceId: selector(state, 'sequenceId'),
}))(OnboardingForm));

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
    return <OnboardingFormWithRedux onSubmit={this.onSubmit} />;
  }
}

export default connect(null, { fetchBootstrap })(OnboardingView);
