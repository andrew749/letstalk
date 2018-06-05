import React, { Component, SFC } from 'react';
import {
  KeyboardAvoidingView,
  Picker,
  ScrollView,
  StyleSheet,
  Text,
} from 'react-native';
import {
  formValueSelector,
  reduxForm,
  Field,
  InjectedFormProps,
  SubmissionError,
} from 'redux-form';
import { connect, ActionCreator } from 'react-redux';
import { FormValidationMessage } from 'react-native-elements';
import { ThunkAction } from 'redux-thunk';
import { bindActionCreators } from 'redux'
import { NavigationScreenProp, NavigationStackAction, NavigationActions } from 'react-navigation';
import Immutable from 'immutable';

import auth from '../services/auth';
import profileService from '../services/profile-service';
import {
  ActionButton,
  ButtonPicker,
  Card,
  FormP,
  FormProps,
  Header,
  LabeledFormInput,
  Loading,
  ModalDatePicker,
  ModalPicker,
  ProfileAvatarEditableFormElement,
} from '../components';
import { genderIdToString } from '../models/user';
import { RootState } from '../redux';
import { State as ProfileState, fetchProfile } from '../redux/profile/reducer';
import { ActionTypes } from '../redux/profile/actions';
import photoService, {PhotoResult} from '../services/photo_service';
import {
  Cohort,
  COHORTS,
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
import Colors from '../services/colors';

interface EditFormData {
  firstName: string;
  lastName: string;
  phoneNumber: string;
  gender: number;
  birthdate: string;
  programId: string,
  sequenceId: string;
  gradYear: number;
  mentorshipPreference: number;
  bio: string | null;
  hometown: string | null;
}

// TODO: move elsewhere
const required = (value: any) => (value ? undefined : 'Required')
const phoneNumber = (value: string) =>
  value && !/^(0|[1-9][0-9]{9})$/i.test(value)
    ? 'Invalid phone number, must be 10 digits'
    : undefined

interface EditFormProps extends FormProps<EditFormData>, EditFormData {
  cohorts: Immutable.List<Cohort>;
}

const EditForm: SFC<FormProps<EditFormData> & EditFormProps> = props => {
  const {
    cohorts,
    error,
    handleSubmit,
    onSubmit,
    reset,
    submitting,
    valid,
    programId,
    sequenceId,
    gradYear
  } = props;
  const buildItems = (rows: Immutable.List<ValueLabel>) => {
    return rows.map(({ value, label }) => {
      return <Picker.Item key={value} label={label} value={value}/>;
    });
  };
  const programItems = buildItems(programOptions(cohorts)).toJS();
  const sequenceItems = buildItems(sequenceOptions(cohorts, programId)).toJS();
  const gradYearItems = buildItems(gradYearOptions(cohorts, programId, sequenceId)).toJS();
  return (
    <ScrollView>
      <Field
        label="First name"
        name="firstName"
        component={LabeledFormInput}
        autoCorrect={false}
        validate={required}
      />
      <Field
        label="Last name"
        name="lastName"
        component={LabeledFormInput}
        autoCorrect={false}
        validate={required}
      />
      <Field
        label="Phone number"
        name="phoneNumber"
        component={LabeledFormInput}
        keyboardType={'phone-pad' as 'phone-pad'}
        validate={phoneNumber}
      />
      <Field
        label="Gender"
        name="gender"
        component={ButtonPicker}
        validate={required}
      >
        <Picker.Item
          label="Male"
          value={2}
        />
        <Picker.Item
          label="Female"
          value={1}
        />
      </Field>
      <Field
        label="Birthday"
        name="birthdate"
        mode={'date' as 'date'}
        defaultDate={new Date('1996-11-07')}
        component={ModalDatePicker}
        validate={required}
      />
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
        placeholder="e.g. I enjoy long walks on the beach"
      />
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
        backgroundColor={Colors.HIVE_PRIMARY}
        style={styles.submitButton}
        disabled={!valid}
        loading={submitting}
        title={submitting ? null : "Save"}
        onPress={handleSubmit(onSubmit)}
      />
    </ScrollView>
  );
}

const cohortSelector = formValueSelector('profile-edit');

const EditFormWithReduxBuilder = (initialValues: EditFormData) => {
  return reduxForm<EditFormData, FormP<EditFormData>>({
    form: 'profile-edit',
    initialValues,
  })(connect((state: RootState) => ({
    programId: cohortSelector(state, 'programId'),
    sequenceId: cohortSelector(state, 'sequenceId'),
    gradYear: cohortSelector(state, 'gradYear'),
    cohorts: COHORTS,
  }))(EditForm));
}

interface DispatchActions {
  fetchProfile: ActionCreator<ThunkAction<Promise<ActionTypes>, ProfileState, void>>;
}

interface Props extends ProfileState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class ProfileEditView extends Component<Props> {
  // TODO: Remove the padding on this component
  static navigationOptions = () => ({
    headerTitle: 'Edit Profile',
    headerStyle: {
      backgroundColor: Colors.HIVE_PRIMARY,
    },
  })

  constructor(props: Props) {
    super(props);

    this.onSubmit = this.onSubmit.bind(this);
  }

  async onSubmit(values: EditFormData) {
    try {
      const {
        firstName,
        lastName,
        phoneNumber,
        gender,
        birthdate,
        programId,
        sequenceId,
        gradYear,
        mentorshipPreference,
        bio,
        hometown,
      } = values;
      const cohortId = getCohortId(COHORTS, programId, sequenceId, gradYear);
      await profileService.profileEdit({
        firstName,
        lastName,
        phoneNumber,
        gender,
        birthdate,
        cohortId,
        mentorshipPreference,
        bio,
        hometown,
      });
      await this.props.fetchProfile();
      this.props.navigation.goBack();
    } catch(e) {
      throw new SubmissionError({_error: e.message});
    }
  }

  render() {
    const {
      firstName,
      lastName,
      gender,
      birthdate,
      phoneNumber,
      programId,
      sequenceId,
      gradYear,
      mentorshipPreference,
      bio,
      hometown,
    } = this.props.profile;
    const EditFormWithRedux = EditFormWithReduxBuilder({
      firstName,
      lastName,
      gender,
      birthdate,
      phoneNumber,
      programId,
      sequenceId,
      gradYear,
      mentorshipPreference,
      bio,
      hometown,
    });
    return (
      <KeyboardAvoidingView behavior="padding">
        <EditFormWithRedux onSubmit={this.onSubmit} />
      </KeyboardAvoidingView>
    );
  }
}

export default connect(({profile}: RootState) => profile, { fetchProfile })(ProfileEditView);

const styles = StyleSheet.create({
  submitButton: {
    marginBottom: 100,
  },
  profilePicContainerStyle: {
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: 20,
    marginLeft: 20,
  }
});
