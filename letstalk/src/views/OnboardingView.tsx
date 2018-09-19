import React, { Component, ReactNode } from 'react';
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
import { KeyboardAwareScrollView } from 'react-native-keyboard-aware-scroll-view'

import { RootState } from '../redux';
import { State } from '../redux/onboarding/reducer';
import {
  ONBOARDING_COHORT,
  ONBOARDING_VECTOR_ME,
  ONBOARDING_VECTOR_YOU,
  ONBOARDING_DONE,
} from '../models/onboarding';
import {
  Cohort,
  OnboardingState,
  OnboardingStatus,
} from '../models';
import {
  setOnboardingStatusAction,
  SetOnboardingStatusAction,
} from '../redux/onboarding/actions';
import Loading from './Loading';
import {
  ActionButton,
  Emoji,
  FormP,
  FormProps,
  Header,
  InfoText,
  LabeledFormInput,
  ModalPicker,
  Rating,
} from '../components';
import profileService, {
  PersonalityVector,
  UserVectorPreferenceType,
} from '../services/profile-service';
import { State as BootstrapState, fetchBootstrap } from '../redux/bootstrap/reducer';
import { State as CohortsState, fetchCohorts } from '../redux/cohorts/reducer';
import { ActionTypes as BootstrapActionTypes} from '../redux/bootstrap/actions';
import { ActionTypes as CohortsActionTypes } from '../redux/cohorts/actions';
import {
  getCohortId,
  programOptions,
  sequenceOptions,
  gradYearOptions,
  ValueLabel,
} from '../models/cohort';
import {
  MENTORSHIP_PREFERENCE_MENTOR,
  MENTORSHIP_PREFERENCE_MENTEE,
  MENTORSHIP_PREFERENCE_NONE,
} from '../models/user';
import { headerStyle, headerTitleStyle, headerTintColor } from './TopHeader';
import { AnalyticsHelper } from '../services';
import Colors from '../services/colors';
import { required } from '../validators';

interface CohortFormData {
  programId: string,
  sequenceId: string;
  gradYear: number;
  mentorshipPreference: number;
  bio: string | null,
  hometown: string | null,
}

interface CohortFormProps extends FormProps<CohortFormData>, CohortFormData {
  cohorts: Immutable.List<Cohort>;
}

const CohortForm: React.SFC<FormProps<CohortFormData> & CohortFormProps>
  = props => {
  const {
    cohorts,
    error,
    handleSubmit,
    onSubmit,
    programId,
    reset,
    sequenceId,
    submitting,
    valid,
  } = props;
  const onSubmitWithReset = async (values: CohortFormData): Promise<void> => {
    await onSubmit(values);
    reset();
  };
  const buildItems = (rows: Immutable.List<ValueLabel>) => {
    return rows.map(({ value, label }) => {
      return <Picker.Item key={value} label={label} value={value}/>;
    });
  };
  const programItems = buildItems(programOptions(cohorts)).toJS();
  const sequenceItems = buildItems(sequenceOptions(cohorts, programId)).toJS();
  const gradYearItems = buildItems(gradYearOptions(cohorts, programId, sequenceId)).toJS();
  return (
    <View style={styles.cohortForm}>
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
      <Field
        label="Mentorship Preference"
        name="mentorshipPreference"
        component={ModalPicker}
        validate={required}
      >
        <Picker.Item key="mentor" label="Mentor" value={MENTORSHIP_PREFERENCE_MENTEE} />
        <Picker.Item key="mentee" label="Mentee" value={MENTORSHIP_PREFERENCE_MENTOR}/>
        <Picker.Item key="none" label="Neither" value={MENTORSHIP_PREFERENCE_NONE}/>
      </Field>
      <Header>Additional Info</Header>
      <Field
        label="Hometown"
        name="hometown"
        component={LabeledFormInput}
        autoCorrect={false}
        placeholder="e.g. Waterloo, ON"
      />
      <Field
        label="Bio"
        name="bio"
        component={LabeledFormInput}
        autoCorrect={false}
        placeholder="e.g. I enjoy Inuit throat singing. (Tell us what you're passionate about, your hobbies, or whatever describes you as a person!)"
        multiline={true}
        numberOfLines={10}
        inputStyle={{width: "100%"}}
      />
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
        backgroundColor={Colors.HIVE_PRIMARY}
        disabled={!valid}
        loading={submitting}
        title={submitting ? null : "Submit"}
        onPress={handleSubmit(onSubmitWithReset)}
      />
    </View>
  );
}

const cohortSelector = formValueSelector('onboarding-cohort');

const CohortFormWithRedux = reduxForm<CohortFormData, FormP<CohortFormData>>({
  form: 'onboarding-cohort',
})(connect((state: RootState) => ({
  programId: cohortSelector(state, 'programId'),
  sequenceId: cohortSelector(state, 'sequenceId'),
  gradYear: cohortSelector(state, 'gradYear'),
  cohorts: state.cohorts.cohorts,
}))(CohortForm));

type PersonalityFormData = PersonalityVector;

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
})(connect((state: RootState) => ({
  values: personalitySelector(state,
    'sociable', 'hardworking', 'ambitious', 'energetic', 'carefree', 'confident'),
}))(PersonalityForm));

interface DispatchActions {
  fetchBootstrap: ActionCreator<ThunkAction<Promise<BootstrapActionTypes>, BootstrapState, void>>;
  fetchCohorts: ActionCreator<ThunkAction<Promise<CohortsActionTypes>, CohortsState, void>>;
  setOnboardingStatusAction(onboardingStatus: OnboardingStatus): SetOnboardingStatusAction;
}

interface Props extends State, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  cohorts: CohortsState;
}

class OnboardingView extends Component<Props> {
  ONBOARDING_VIEW_IDENTIFIER = "OnboardingView";

  static navigationOptions = {
    headerTitle: 'Onboarding',
    headerStyle, 
    headerTitleStyle, 
    headerTintColor 
  }

  constructor(props: Props) {
    super(props);

    this.onSubmitCohort = this.onSubmitCohort.bind(this);
    this.load = this.load.bind(this);
    this.renderBody = this.renderBody.bind(this);
  }

  async componentDidMount() {
    AnalyticsHelper.getInstance().recordPage(this.ONBOARDING_VIEW_IDENTIFIER);
    await this.load();
  }

  private async load() {
    await this.props.fetchCohorts();
  }

  async onSubmitCohort(values: CohortFormData) {
    const { programId, sequenceId, gradYear, mentorshipPreference, bio, hometown } = values;
    const cohortId = getCohortId(this.props.cohorts.cohorts, programId, sequenceId, gradYear);
    try {
      const onboardingStatus = await profileService.updateCohort({
        cohortId,
        mentorshipPreference,
        bio,
        hometown,
      });
      this.props.setOnboardingStatusAction(onboardingStatus);
      await this.props.fetchBootstrap();
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        actions: [NavigationActions.navigate({ routeName: 'Tabbed' })]
      }));
    } catch(e) {
      throw new SubmissionError({_error: e.errorMsg});
    }
  }

  async onSubmitPersonality(preferenceType: UserVectorPreferenceType, values: PersonalityFormData) {
    try {
      const onboardingStatus = await profileService.updateVector(preferenceType, values);
      this.props.setOnboardingStatusAction(onboardingStatus);
    } catch(e) {
      throw new SubmissionError({_error: e.errorMsg});
    }
  }

  renderBody() {
    const { state } = this.props.onboardingStatus;
    switch (state) {
      case ONBOARDING_COHORT:
        // TODO: Update copy here
        return (
          <KeyboardAwareScrollView>
            <Header>Your Cohort</Header>
            <InfoText>
              Based on your information, we'll be better able to match you with a mentor/mentee!
            </InfoText>
            <CohortFormWithRedux onSubmit={this.onSubmitCohort} />
          </KeyboardAwareScrollView>
        );
      case ONBOARDING_VECTOR_ME:
        // NOTE: This will not show up now
        const onSubmitMine = async (values: PersonalityFormData) => {
          await this.onSubmitPersonality(UserVectorPreferenceType.PREFERENCE_TYPE_ME, values);
        };
        return (
          <ScrollView>
            <Header>Your Personality</Header>
            <InfoText>
              Looks like you're going to be a big <Emoji name="man"/>! We'd like to get
              to know you a little better, so that we can find you some lit <Emoji name="fire"/>
              noobies to mentor.
            </InfoText>
            <PersonalityFormWithRedux onSubmit={onSubmitMine} />
          </ScrollView>
        );
      case ONBOARDING_VECTOR_YOU:
        // NOTE: This will not show up now
        const onSubmitYour = async (values: PersonalityFormData) => {
          await this.onSubmitPersonality(UserVectorPreferenceType.PREFERENCE_TYPE_YOU, values);
          // Reload bootstrap data after updating
          await this.props.fetchBootstrap();
          this.props.navigation.dispatch(NavigationActions.reset({
            index: 0,
            actions: [NavigationActions.navigate({ routeName: 'Tabbed' })]
          }));
        };
        return (
          <ScrollView>
            <Header>Their Personality</Header>
            <InfoText>
              We'd also like to get a sense of what kind of noobies <Emoji name="baby"/> you would
              like to mentor.
            </InfoText>
            <PersonalityFormWithRedux onSubmit={onSubmitYour} />
          </ScrollView>
        );
      case ONBOARDING_DONE:
        // TODO: What to do in this case
        return (
          <ScrollView>
            <Header>You're done</Header>
          </ScrollView>
        );
      default:
        // Ensure exhaustiveness of select
        const _: never = state;
    }
  }

  render() {
    const {
      state,
      errorMsg,
      errorType,
    } = this.props.cohorts.fetchState;
    return (
      <Loading
        state={state}
        errorMsg={errorMsg}
        errorType={errorType}
        load={this.load}
        renderBody={this.renderBody}
        navigation={this.props.navigation}
      />
    );
  }
}

const styles = StyleSheet.create({
  actionButton: {
    marginBottom: 10,
  },
  cohortForm: {
    paddingBottom: 100,
  },
  hint: {
    color: 'gray',
    fontSize: 14,
    marginTop: -10,
    marginLeft: 10,
  },
  longForm: {
    height: 40
  }
});

export default connect(({ onboarding, cohorts }: RootState) => {
  return { ...onboarding, cohorts }
}, { fetchBootstrap, fetchCohorts, setOnboardingStatusAction })(OnboardingView);
