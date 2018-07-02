import React, { Component, SFC } from 'react';
import {
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
import { FormValidationMessage, FormInputProps } from 'react-native-elements';
import { ThunkAction } from 'redux-thunk';
import { bindActionCreators } from 'redux'
import { NavigationScreenProp, NavigationStackAction, NavigationActions } from 'react-navigation';
import Immutable from 'immutable';
import Moment from 'moment';
import { KeyboardAwareScrollView } from 'react-native-keyboard-aware-scroll-view'

import auth from '../services/auth';
import profileService from '../services/profile-service';
import { State as CohortsState, fetchCohorts } from '../redux/cohorts/reducer';
import { ActionTypes as CohortsActionTypes } from '../redux/cohorts/actions';

import {
  ActionButton,
  ButtonPicker,
  Card,
  FormP,
  FormProps,
  Header,
  LabeledFormInput,
  ModalDatePicker,
  ModalPicker,
  ProfileAvatarEditableFormElement,
} from '../components';
import Loading from './Loading';
import { genderIdToString } from '../models/user';
import { RootState } from '../redux';
import { State as ProfileState, fetchProfile } from '../redux/profile/reducer';
import { ActionTypes } from '../redux/profile/actions';
import photoService, {PhotoResult} from '../services/photo_service';
import {
  Cohort,
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
import { headerStyle } from './TopHeader';
import { AnalyticsHelper } from '../services';

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
  profilePic: PhotoResult;
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

class EditFormRefs {
  firstNameFieldRef: Field<FormInputProps>;
  lastNameFieldRef: Field<FormInputProps>;
  phoneNumberFieldRef: Field<FormInputProps>;
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
  const fieldRefs = new EditFormRefs();
  return (
    <KeyboardAwareScrollView
      keyboardShouldPersistTaps={true}
      >
      <Header>Personal Info</Header>
      <Field
        name="profilePic"
        component={ProfileAvatarEditableFormElement}
        containerStyle={styles.profilePicContainerStyle}
      />
      <Field
        label="First name"
        name="firstName"
        component={LabeledFormInput}
        ref={(ref: Field<FormInputProps>) => fieldRefs.firstNameFieldRef = ref}
        onSubmitEditing={() => {
          // @ts-ignore
          fieldRefs.lastNameFieldRef.getRenderedComponent().focus();
        }}
        withRef={true}
        autoCorrect={false}
        validate={required}
      />
      <Field
        label="Last name"
        name="lastName"
        component={LabeledFormInput}
        ref={(ref: Field<FormInputProps>) => fieldRefs.lastNameFieldRef = ref}
        onSubmitEditing={() => {
          // @ts-ignore
          fieldRefs.phoneNumberFieldRef.getRenderedComponent().focus();
        }}
        withRef={true}
        autoCorrect={false}
        validate={required}
      />
      <Field
        label="Phone number"
        name="phoneNumber"
        component={LabeledFormInput}
        ref={(ref: Field<FormInputProps>) => fieldRefs.phoneNumberFieldRef = ref}
        withRef={true}
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
        component={ModalDatePicker}
        validate={required}
      />
      <Header>Your Cohort</Header>
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
      <Text style={styles.hint}>Optional</Text>
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
    </KeyboardAwareScrollView>
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
    cohorts: state.cohorts.cohorts,
  }))(EditForm));
}

interface DispatchActions {
  fetchProfile: ActionCreator<ThunkAction<Promise<ActionTypes>, ProfileState, void>>;
  fetchCohorts: ActionCreator<ThunkAction<Promise<CohortsActionTypes>, CohortsState, void>>;
}

interface Props extends ProfileState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  cohorts: CohortsState,
}

class ProfileEditView extends Component<Props> {
  EDIT_PROFILE_VIEW_IDENTIFIER = "ProfileEditView";

  static navigationOptions = () => ({
    headerTitle: 'Edit Profile',
    headerStyle,
  })

  constructor(props: Props) {
    super(props);

    this.onSubmit = this.onSubmit.bind(this);
    this.load = this.load.bind(this);
    this.renderBody = this.renderBody.bind(this);
  }

  async componentDidMount() {
    AnalyticsHelper.getInstance().recordPage(this.EDIT_PROFILE_VIEW_IDENTIFIER);
    await this.load();
  }

  private async load() {
    await this.props.fetchCohorts();
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
        profilePic,
      } = values;
      const cohortId = getCohortId(this.props.cohorts.cohorts, programId, sequenceId, gradYear);
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
      // update profile pic
      await photoService.uploadProfilePhoto(profilePic.uri);
      await this.props.fetchProfile();
      this.props.navigation.goBack();
    } catch(e) {
      throw new SubmissionError({_error: e.errorMsg});
    }
  }

  renderBody() {
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
      profilePic,
    } = this.props.profile;
    const EditFormWithRedux = EditFormWithReduxBuilder({
      firstName,
      lastName,
      gender,
      birthdate: Moment.utc(birthdate).format("YYYY-MM-DD"),
      phoneNumber,
      programId,
      sequenceId,
      gradYear,
      mentorshipPreference,
      bio,
      hometown,
      profilePic: {uri: profilePic, data: null},
    });
    return (
      <EditFormWithRedux onSubmit={this.onSubmit} />
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

export default connect(({profile, cohorts}: RootState) => {
  return { ...profile, cohorts }
}, { fetchProfile, fetchCohorts })(ProfileEditView);

const styles = StyleSheet.create({
  submitButton: {
    marginBottom: 100,
  },
  profilePicContainerStyle: {
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: 20,
    marginLeft: 20,
  },
  hint: {
    color: 'gray',
    fontSize: 14,
    marginTop: -10,
    marginLeft: 10,
  },
});
