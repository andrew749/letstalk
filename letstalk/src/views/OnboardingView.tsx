import React, { Component } from 'react';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import { Picker, ScrollView, StyleSheet, Text, View } from 'react-native';
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
  Rating,
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

interface CohortFormData {
  programId: string,
  sequenceId: string;
  gradYear: string;
}

// TODO: move elsewhere
const required = (value: any) => (value ? undefined : 'Required')

const CohortForm: React.SFC<FormProps<CohortFormData> & CohortFormData> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting, valid, programId, sequenceId } = props;
  const onSubmitWithReset = async (values: CohortFormData): Promise<void> => {
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
    <View>
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
    </View>
  );
}

const cohortSelector = formValueSelector('onboarding-cohort');

const CohortFormWithRedux = reduxForm<CohortFormData, FormP<CohortFormData>>({
  form: 'onboarding-cohort',
})(connect(state => ({
  programId: cohortSelector(state, 'programId'),
  sequenceId: cohortSelector(state, 'sequenceId'),
  gradYear: cohortSelector(state, 'gradYear'),
}))(CohortForm));

interface PersonalityFormData {
  sociable: number;
  hardworking: number;
  ambitious: number;
  energetic: number;
  carefree: number;
  confident: number;
}

const vectorsWithLabels = Immutable.Map({
  sociable: 'Sociable',
  hardworking: 'Hardworking',
  ambitious: 'Ambitious',
  energetic: 'Energetic',
  carefree: 'Carefree',
  confident: 'Confident',
});

interface PersonalityProps {
  values: PersonalityFormData;
}

const TOTAL_POINTS = 15;

const PersonalityForm: React.SFC<FormProps<PersonalityFormData> & PersonalityProps> = props => {
  const { error, handleSubmit, onSubmit, reset, submitting, valid } = props;
  if (!props.initialized) {
    props.initialize({
      sociable: 4,
      hardworking: 0,
      ambitious: 0,
      energetic: 0,
      carefree: 0,
      confident: 0,
    });
  }
  const values: Immutable.Map<string, number> = Immutable.Map(props.values);
  const total = values.reduce((acc, value) => acc + value, 0);
  const onSubmitWithReset = async (values: PersonalityFormData): Promise<void> => {
    await onSubmit(values);
    reset();
  };
  let buttonLabel = "Submit Personality";
  if (total < TOTAL_POINTS) buttonLabel = "You have " + (TOTAL_POINTS - total) + " points left";
  if (total > TOTAL_POINTS) buttonLabel = (total - TOTAL_POINTS) + " points too many";
  const fields = vectorsWithLabels.map((label, name) => (
    <Field
      key={name}
      label={label}
      name={name}
      numElems={5}
      component={Rating}
    />
  )).toArray();
  return (
    <View>
      <InfoText>
        You have <Text style={{fontWeight: 'bold'}}>{total + '/' + TOTAL_POINTS}</Text> points.
      </InfoText>
      {fields}
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
        style={styles.actionButton}
        disabled={total !== TOTAL_POINTS}
        loading={submitting}
        title={submitting ? null : buttonLabel}
        onPress={handleSubmit(onSubmitWithReset)}
      />
    </View>
  );
}

const personalitySelector = formValueSelector('onboarding-personality');

const PersonalityFormWithRedux = reduxForm<PersonalityFormData, FormP<PersonalityFormData>>({
  form: 'onboarding-personality',
})(connect(state => ({
  values: personalitySelector(state,
    'sociable', 'hardworking', 'ambitious', 'energetic', 'carefree', 'confident'),
}))(PersonalityForm));

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

  async onSubmitCohort(values: CohortFormData) {
    const { programId, sequenceId, gradYear } = values;
    const cohortId = getCohortId(programId, sequenceId, gradYear);
    try {
      // await profileService.updateCohort({ cohortId });
    } catch(e) {
      throw new SubmissionError({_error: e.message});
    }
    this.props.setStep(Step.MY_VECTOR);
  }

  async onSubmit(values: CohortFormData) {
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
            <CohortFormWithRedux onSubmit={this.onSubmitCohort} />
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
            <PersonalityFormWithRedux onSubmit={(values: PersonalityFormData) => null as void} />
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

const styles = StyleSheet.create({
  actionButton: {
    marginBottom: 10,
  },
});

export default connect(({ onboarding }: RootState) => onboarding,
  { fetchBootstrap, setStep }
)(OnboardingView);
