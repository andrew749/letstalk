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
import { Cohort } from '../models';
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
import profileService from '../services/profile-service';
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
import { GROUP_GENERIC } from '../services/survey';

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
        label="Your Preferred Role"
        name="mentorshipPreference"
        component={ModalPicker}
        validate={required}
      >
        <Picker.Item key="mentor" label="Mentor" value={MENTORSHIP_PREFERENCE_MENTEE} />
        <Picker.Item key="mentee" label="Mentee" value={MENTORSHIP_PREFERENCE_MENTOR}/>
        <Picker.Item key="none" label="I don't know yet" value={MENTORSHIP_PREFERENCE_NONE}/>
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

interface DispatchActions {
  fetchBootstrap: ActionCreator<ThunkAction<Promise<BootstrapActionTypes>, BootstrapState, void>>;
  fetchCohorts: ActionCreator<ThunkAction<Promise<CohortsActionTypes>, CohortsState, void>>;
}

interface Props extends DispatchActions {
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
      await profileService.updateCohort({
        cohortId,
        mentorshipPreference,
        bio,
        hometown,
      });
      await this.props.fetchBootstrap();
      await this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        actions: [NavigationActions.navigate({ routeName: 'Tabbed' })]
      }));
    } catch(e) {
      throw new SubmissionError({_error: e.errorMsg});
    }
  }

  renderBody() {
    return (
      <KeyboardAwareScrollView>
        <Header>Your Cohort</Header>
        <InfoText>
          Based on your information, we'll be better able to match you with a mentor/mentee!
        </InfoText>
        <CohortFormWithRedux onSubmit={this.onSubmitCohort} />
      </KeyboardAwareScrollView>
    );
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

export default connect(({ cohorts }: RootState) => {
  return { cohorts }
}, { fetchBootstrap, fetchCohorts })(OnboardingView);
