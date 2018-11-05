import React, { Component, SFC } from 'react';
import {
  Dimensions,
  EmitterSubscription,
  Picker,
  Platform,
  Keyboard,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from 'react-native';
import {
  reduxForm,
  Field,
  SubmissionError,
} from 'redux-form';
import { connect, ActionCreator } from 'react-redux';
import { FormValidationMessage, FormInputProps } from 'react-native-elements';
import { ThunkAction } from 'redux-thunk';
import { NavigationScreenProp, NavigationStackAction, NavigationActions } from 'react-navigation';
import Immutable from 'immutable';
import Moment from 'moment';
import { KeyboardAwareScrollView } from 'react-native-keyboard-aware-scroll-view'

import profileService from '../services/profile-service';
import { State as CohortsState, fetchCohorts } from '../redux/cohorts/reducer';
import { ActionTypes as CohortsActionTypes } from '../redux/cohorts/actions';
import { combineFetchStates } from '../redux/actions';

const SCREEN_WIDTH = Dimensions.get('window').width;

import {
  ActionButton,
  ButtonPicker,
  FormP,
  FormProps,
  Header,
  LabeledFormInput,
  ModalDatePicker,
  ModalPicker,
} from '../components';
import Loading from './Loading';
import { RootState } from '../redux';
import { State as ProfileState, fetchProfile } from '../redux/profile/reducer';
import { ActionTypes as ProfileActionTypes } from '../redux/profile/actions';
import {
  Cohort,
  getCohortId
} from '../models/cohort';
import {
  MENTORSHIP_PREFERENCE_MENTOR,
  MENTORSHIP_PREFERENCE_MENTEE,
  MENTORSHIP_PREFERENCE_NONE,
  GenderId,
} from '../models/user';
import Colors from '../services/colors';
import { headerStyle, headerTitleStyle, headerTintColor } from './TopHeader';
import { AnalyticsHelper } from '../services';
import { required, phoneNumber} from '../validators';

interface EditFormData {
  firstName: string;
  lastName: string;
  phoneNumber: string;
  gender: number;
  birthdate: string;
  mentorshipPreference: number;
  bio: string | null;
  hometown: string | null;
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
      error,
      handleSubmit,
      onSubmit,
      reset,
      submitting,
      valid
    } = this.props;
    const { isKeyboardShown } = this.state;
    const fieldRefs = new EditFormRefs();
    return (
      <View>
        <KeyboardAwareScrollView
          keyboardShouldPersistTaps="always"
        >
          <Header>Personal Info</Header>
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
            label="Gender (optional)"
            name="gender"
            component={ButtonPicker}
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
            label="Birthday (optional)"
            name="birthdate"
            mode={'date' as 'date'}
            androidMode={'spinner' as 'spinner'}
            component={ModalDatePicker}
          />
          {(Platform.OS === 'android') && <Text style={[styles.hint, styles.mentorshipLabel]}>
            Mentorship Preference
          </Text>}
          <Field
            label="Your Preferred Role"
            name="mentorshipPreference"
            component={ModalPicker}
            validate={required}
            containerStyle={styles.adjustMargin}
          >
            <Picker.Item key="mentor" label="Mentor" value={MENTORSHIP_PREFERENCE_MENTEE} />
            <Picker.Item key="mentee" label="Mentee" value={MENTORSHIP_PREFERENCE_MENTOR}/>
            <Picker.Item key="none" label="I don't know yet" value={MENTORSHIP_PREFERENCE_NONE}/>
          </Field>
          <Header>Additional Info</Header>
          <Text style={[styles.hint, styles.adjustMargin]}>Optional info to let others know you better</Text>
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
            multiline={true}
            numberOfLines={10}
            inputStyle={{width: "100%"}}
            placeholder="e.g. I enjoy Inuit throat singing. (Tell us what you're passionate about, your hobbies, or whatever describes you as a person!)"
          />
          <ActionButton
            backgroundColor={Colors.HIVE_PRIMARY}
            disabled={!valid}
            loading={submitting}
            title={submitting ? null : "Save"}
            onPress={handleSubmit(onSubmit)}
           />
          {error && <FormValidationMessage>{error}</FormValidationMessage>}
        </KeyboardAwareScrollView>
      </View>
    );
  }
}

const EditFormWithReduxBuilder = (initialValues: EditFormData) => {
  return reduxForm<EditFormData, FormP<EditFormData>>({
    form: 'profile-edit',
    initialValues,
  })(connect()(EditForm));
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
    if (!this.props.profile) {
      // Don't await on purpose
      this.props.fetchProfile();
    }
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
        mentorshipPreference,
        bio,
        hometown
      } = values;
      const {
        programId,
        sequenceId,
        gradYear
      } = this.props.profile;
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
      phoneNumber,
      gender = GenderId.Unspecified,
      birthdate,
      mentorshipPreference,
      bio,
      hometown,
    } = this.props.profile;
    const EditFormWithRedux = EditFormWithReduxBuilder({
      firstName,
      lastName,
      phoneNumber,
      gender,
      birthdate: birthdate && Moment.utc(birthdate).format("YYYY-MM-DD"),
      mentorshipPreference,
      bio,
      hometown,
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
    } = combineFetchStates(this.props.fetchState, this.props.cohorts.fetchState);
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

export default connect(({ profile, cohorts }: RootState) => {
  return { ...profile, cohorts }
}, { fetchProfile, fetchCohorts })(ProfileEditView);

const styles = StyleSheet.create({
  adjustMargin: {
    marginHorizontal: 10
  },
  hint: {
    color: Colors.HIVE_SUBDUED,
    fontSize: 14,
  },
  mentorshipLabel: {
    marginTop: 10,
    marginHorizontal: 15
  }
});
