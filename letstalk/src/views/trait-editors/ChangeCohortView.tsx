import React, { Component, SFC } from 'react';
import { Picker, Text } from 'react-native';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import {
  NavigationScreenProp,
  NavigationStackAction,
} from 'react-navigation';
import Immutable from 'immutable';
import { FormValidationMessage } from 'react-native-elements';
import {
  formValueSelector,
  reduxForm,
  Field,
  InjectedFormProps,
  SubmissionError,
} from 'redux-form';
import { KeyboardAwareScrollView } from 'react-native-keyboard-aware-scroll-view'
import { ThunkAction } from 'redux-thunk';

import { infoToast } from '../../redux/toast';
import { RootState } from '../../redux';
import {
  Cohort,
  getCohortId,
  programOptions,
  sequenceOptions,
  gradYearOptions,
  ValueLabel,
} from '../../models/cohort';
import { State as CohortsState, fetchCohorts } from '../../redux/cohorts/reducer';
import { ActionTypes as CohortsActionTypes } from '../../redux/cohorts/actions';
import { State as ProfileState, fetchProfile } from '../../redux/profile/reducer';
import { ActionTypes as ProfileActionTypes } from '../../redux/profile/actions';
import { required } from '../../validators';
import Loading from '../Loading';
import {
  ActionButton,
  FormProps,
  FormP,
  Header,
  ModalPicker,
} from '../../components';
import { headerStyle } from '../TopHeader';
import { AnalyticsHelper } from '../../services/analytics';
import Colors from '../../services/colors';
import profileService from '../../services/profile-service';

interface ChangeCohortFormData {
  programId: string,
  sequenceId: string;
  gradYear: number;
}

interface ChangeCohortFormProps extends FormProps<ChangeCohortFormData>, ChangeCohortFormData {
  cohorts: Immutable.List<Cohort>;
}

const ChangeCohortForm: SFC<FormProps<ChangeCohortFormData> & ChangeCohortFormProps> = props => {
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
    gradYear,
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
    <KeyboardAwareScrollView
      keyboardShouldPersistTaps={true}
    >
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
        backgroundColor={Colors.HIVE_PRIMARY}
        disabled={!valid}
        loading={submitting}
        title={submitting ? null : "Save"}
        onPress={handleSubmit(onSubmit)}
      />
    </KeyboardAwareScrollView>
  );
}

const cohortSelector = formValueSelector('profile-edit');

const ChangeCohortFormWithReduxBuilder = (initialValues: ChangeCohortFormData) => {
  return reduxForm<ChangeCohortFormData, FormP<ChangeCohortFormData>>({
    form: 'profile-edit',
    initialValues,
  })(connect((state: RootState) => ({
    programId: cohortSelector(state, 'programId'),
    sequenceId: cohortSelector(state, 'sequenceId'),
    gradYear: cohortSelector(state, 'gradYear'),
    cohorts: state.cohorts.cohorts,
  }))(ChangeCohortForm));
}

interface DispatchActions {
  fetchProfile: ActionCreator<ThunkAction<Promise<ProfileActionTypes>, ProfileState, void>>;
  fetchCohorts: ActionCreator<ThunkAction<Promise<CohortsActionTypes>, CohortsState, void>>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

// TODO: Maybe take current cohort info as a prop instead of pulling from redux. Makes this more
// reusable.
interface Props extends DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction & ChangeCohortFormData>;
  cohorts: CohortsState;
}

class ChangeCohortView extends Component<Props> {
  CHANGE_COHORT_VIEW_IDENTIFIER = 'ChangeCohortView';

  static navigationOptions = {
    headerTitle: 'Change Cohort',
    headerStyle,
  }

  constructor(props: Props) {
    super(props);

    this.load = this.load.bind(this);
    this.onSubmit = this.onSubmit.bind(this);
    this.renderBody = this.renderBody.bind(this);
  }

  private async load() {
    await this.props.fetchCohorts();
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.CHANGE_COHORT_VIEW_IDENTIFIER);
    });
    await this.load();
  }

  private async onSubmit(values: ChangeCohortFormData) {
    try {
      const {
        programId,
        sequenceId,
        gradYear,
      } = values;
      const cohortId = getCohortId(this.props.cohorts.cohorts, programId, sequenceId, gradYear);
      await profileService.updateCohort({ cohortId });
      await this.props.infoToast('Successfully changed cohort');
      await this.props.fetchProfile();
      await this.props.navigation.goBack();
    } catch (e) {
      throw new SubmissionError({_error: e.errorMsg});
    }
  }

  private renderBody() {
    const { navigation } = this.props;
    const programId = navigation.getParam('programId', null);
    const sequenceId = navigation.getParam('sequenceId', null);
    const gradYear = navigation.getParam('gradYear', null);
    const ChangeCohortFormWithRedux = ChangeCohortFormWithReduxBuilder({
      programId,
      sequenceId,
      gradYear,
    });
    return (
      <ChangeCohortFormWithRedux onSubmit={this.onSubmit} />
    )
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

export default connect(({cohorts}: RootState) => {
  return { cohorts }
}, { fetchProfile, fetchCohorts, infoToast })(ChangeCohortView);
