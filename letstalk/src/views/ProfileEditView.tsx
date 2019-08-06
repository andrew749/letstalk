import React, { Component, SFC } from 'react';
import {
  Dimensions,
  EmitterSubscription,
  Picker,
  Keyboard,
  ScrollView,
  StyleSheet,
  Text,
  View,
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
import { GenderId } from '../models/user';

import {
  ActionButton,
  FloatingButton,
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
import { ActionTypes as ProfileActionTypes } from '../redux/profile/actions';
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
import { headerStyle, headerTitleStyle, headerTintColor } from './TopHeader';
import { AnalyticsHelper } from '../services';
import { required, phoneNumber } from '../validators';

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

interface EditFormProps extends FormProps<EditFormData>, EditFormData {
  cohorts: Immutable.List<Cohort>;
}

class EditFormRefs {
  firstNameFieldRef: Field<FormInputProps>;
  lastNameFieldRef: Field<FormInputProps>;
  phoneNumberFieldRef: Field<FormInputProps>;
}

interface State {
  readonly isKeyboardShown: boolean;
}

type EditFormComponentProps = FormProps<EditFormData> & EditFormProps;

class EditForm extends Component<EditFormComponentProps, State> {
  private keyboardDidHideListener: EmitterSubscription;
  private keyboardDidShowListener: EmitterSubscription;

  constructor(props: EditFormComponentProps) {
    super(props)

    this.state = { isKeyboardShown: false };

    this.keyboardDidShow = this.keyboardDidShow.bind(this)
    this.keyboardDidHide = this.keyboardDidHide.bind(this)
  }

  componentDidMount () {
    this.keyboardDidShowListener = Keyboard.addListener('keyboardDidShow', this.keyboardDidShow);
    this.keyboardDidHideListener = Keyboard.addListener('keyboardDidHide', this.keyboardDidHide);
  }

  componentWillUnmount () {
    this.keyboardDidShowListener.remove();
    this.keyboardDidHideListener.remove();
  }

  private keyboardDidShow () {
    this.setState({ isKeyboardShown: true });
  }

  private keyboardDidHide () {
    this.setState({ isKeyboardShown: false });
  }

  render() {
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
    } = this.props;
    const { isKeyboardShown } = this.state;
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
      <View>
        <KeyboardAwareScrollView
          keyboardShouldPersistTaps="always"
        >
          <Header>Personal Info</Header>
          <View style={styles.profilePicContainer}>
            <Field
              name="profilePic"
              component={ProfileAvatarEditableFormElement as "input" & typeof ProfileAvatarEditableFormElement}
              containerStyle={styles.profilePicContainerStyle}
            />
          </View>
          <Field
            label="First name"
            name="firstName"
            component={LabeledFormInput as "input" & typeof LabeledFormInput}
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
            component={LabeledFormInput as "input" & typeof LabeledFormInput}
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
            component={LabeledFormInput as "input" & typeof LabeledFormInput}
            ref={(ref: Field<FormInputProps>) => fieldRefs.phoneNumberFieldRef = ref}
            withRef={true}
            keyboardType={'phone-pad' as 'phone-pad'}
            validate={phoneNumber}
          />
          <Field
            label="Gender (optional)"
            name="gender"
            component={ButtonPicker as "input" & typeof ButtonPicker}
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
            androidMode={'spinner' as 'spinner'}
            component={ModalDatePicker as "input" & typeof ModalDatePicker}
            validate={required}
          />
          <Header>Your Cohort</Header>
          <Field
            label="Program"
            name="programId"
            component={ModalPicker as "input" & typeof ModalPicker}
            validate={required}
          >
            {programItems}
          </Field>
          <Field
            label="Sequence"
            name="sequenceId"
            component={ModalPicker as "input" & typeof ModalPicker}
            validate={required}
          >
            {sequenceItems}
          </Field>
          <Field
            label="Grad Year"
            name="gradYear"
            component={ModalPicker as "input" & typeof ModalPicker}
            validate={required}
          >
            {gradYearItems}
          </Field>
          <Field
            label="Your Preferred Role"
            name="mentorshipPreference"
            component={ModalPicker as "input" & typeof ModalPicker}
            validate={required}
          >
            <Picker.Item key="mentor" label="Mentor" value={MENTORSHIP_PREFERENCE_MENTEE} />
            <Picker.Item key="mentee" label="Mentee" value={MENTORSHIP_PREFERENCE_MENTOR}/>
            <Picker.Item key="none" label="I don't know yet" value={MENTORSHIP_PREFERENCE_NONE}/>
          </Field>
          <Header>Additional Info</Header>
          <Text style={styles.hint}>Optional</Text>
          <Field
            label="Hometown"
            name="hometown"
            component={LabeledFormInput as "input" & typeof LabeledFormInput}
            autoCorrect={false}
            placeholder="e.g. Waterloo, ON"
          />
          <Field
            label="Bio"
            name="bio"
            component={LabeledFormInput as "input" & typeof LabeledFormInput}
            multiline={true}
            numberOfLines={10}
            inputStyle={{width: "100%"}}
            containerStyle={!isKeyboardShown && { marginBottom: 40 }}
            placeholder="e.g. I enjoy Inuit throat singing. (Tell us what you're passionate about, your hobbies, or whatever describes you as a person!)"
          />
          {isKeyboardShown && <ActionButton
            backgroundColor={Colors.HIVE_PRIMARY}
            disabled={!valid}
            loading={submitting}
            title={submitting ? null : "Save"}
            onPress={handleSubmit(onSubmit)}
          />}
          {error && <FormValidationMessage>{error}</FormValidationMessage>}
        </KeyboardAwareScrollView>
        <FloatingButton
          backgroundColor={Colors.HIVE_PRIMARY}
          disabled={!valid}
          loading={submitting}
          title={submitting ? null : "Save"}
          onPress={handleSubmit(onSubmit)}
        />
      </View>
    );
  }
}

const cohortSelector = formValueSelector('profile-edit');

const EditFormWithReduxBuilder = (initialValues: EditFormData) => {
  return reduxForm<EditFormData, FormP<EditFormData>>({
    form: 'profile-edit',
    initialValues,
  })(connect<{}, {}, EditFormComponentProps>((state: RootState) => ({
    programId: cohortSelector(state, 'programId'),
    sequenceId: cohortSelector(state, 'sequenceId'),
    gradYear: cohortSelector(state, 'gradYear'),
    cohorts: state.cohorts.cohorts,
  }))(EditForm as any));
}

interface DispatchActions {
  fetchProfile: ActionCreator<ThunkAction<Promise<ProfileActionTypes>, ProfileState, void>>;
  fetchCohorts: ActionCreator<ThunkAction<Promise<CohortsActionTypes>, CohortsState, void>>;
}

interface Props extends ProfileState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  cohorts: CohortsState,
}

class ProfileEditView extends Component<Props> {
  EDIT_PROFILE_VIEW_IDENTIFIER = "ProfileEditView";

  static navigationOptions = () => ({
    headerTitle: 'Edit Personal Info',
    headerStyle,
    headerTitleStyle,
    headerTintColor
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

  private async onSubmit(values: EditFormData) {
    try {
      const {
        firstName,
        lastName,
        phoneNumber,
        gender = GenderId.Unspecified,
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
      if (profilePic && profilePic.uri) {
        await photoService.uploadProfilePhoto(profilePic.uri);
      }
      await this.props.fetchProfile();
      await this.props.navigation.goBack();
    } catch(e) {
      throw new SubmissionError({_error: e.errorMsg});
    }
  }

  private renderBody() {
    const {
      firstName,
      lastName,
      gender = GenderId.Unspecified,
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
      profilePic: {uri: profilePic},
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
}, { fetchProfile, fetchCohorts })(ProfileEditView as any);

const styles = StyleSheet.create({
  profilePicContainer: {
    justifyContent: 'center',
    alignItems: 'center',
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
