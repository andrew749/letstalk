import React, { Component, SFC, ReactNode } from 'react';
import { Picker, StyleSheet, Text, View } from 'react-native';
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
  FormErrors,
} from 'redux-form';
import Sentry from 'sentry-expo';
import { KeyboardAwareScrollView } from 'react-native-keyboard-aware-scroll-view'
import { ThunkAction } from 'redux-thunk';
import _ from 'underscore';
import Moment from 'moment';

import { infoToast } from '../../redux/toast';
import { RootState } from '../../redux';
import { required } from '../../validators';
import Loading from '../Loading';
import {
  ActionButton,
  AutocompleteInput,
  FormProps,
  FormP,
  Header,
  ModalDatePicker,
} from '../../components';
import { DataItem, Select } from '../../components/AutocompleteInput';
import { headerStyle } from '../TopHeader';
import { AnalyticsHelper } from '../../services/analytics';
import Colors from '../../services/colors';
import autocompleteService from '../../services/autocomplete-service';
import { ActionTypes } from '../../redux/profile/actions';
import {
  State as ProfileState,
  addPosition,
} from '../../redux/profile/reducer';

interface AddPositionFormData {
  role: Select;
  organization: Select;
  startDate: string;
  endDate: string;
}

const THROTTLE_TIME = 250; // ms

const onRoleQueryChange = async (query: string, setData: (data: Array<DataItem>) => void) => {
  let res: Array<DataItem> = [];
  if (query !== '') res = await autocompleteService.autocompleteRole(query, 10);
  setData(res);
}
const onRoleQueryChangeThrottled = _.throttle(onRoleQueryChange, THROTTLE_TIME);

const onOrganizationQueryChange = async (query: string, setData: (data: Array<DataItem>) => void) => {
  let res: Array<DataItem> = [];
  if (query !== '') res = await autocompleteService.autocompleteOrganization(query, 10);
  setData(res);
}
const onOrganizationQueryChangeThrottled = _.throttle(onOrganizationQueryChange, THROTTLE_TIME);

const AddPositionForm: SFC<FormProps<AddPositionFormData> & AddPositionFormData> = props => {

  const {
    error,
    handleSubmit,
    onSubmit,
    reset,
    submitting,
    valid,
    role,
    organization,
    startDate,
    endDate,
  } = props;
  let selection: ReactNode = null;


  if (!!role && !!organization) {
    const dateFmt = "MMM Do, YYYY";
    const untilSelection = !endDate ? 'present' : Moment(endDate).format(dateFmt);
    const dateSelection = !startDate ? null : (
      <Text>
        <Text>{' from '}</Text>
        <Text style={{fontWeight: "900"}}>{Moment(startDate).format(dateFmt)}</Text>
        <Text>{' to '}</Text>
        <Text style={{fontWeight: "900"}}>{untilSelection}</Text>
      </Text>
    );
    selection = (
      <View style={styles.selectionContainer}>
        <Text style={styles.selectionText}>
          <Text>{'You want to add the position: '}</Text>
          <Text style={{fontWeight: "900"}}>{role.name}</Text>
          <Text>{' @ '}</Text>
          <Text style={{fontWeight: "900"}}>{organization.name}</Text>
          { dateSelection }
          <Text>{'.'}</Text>
        </Text>
      </View>
    );
  }

  let roleSelection: ReactNode = null;
  let organizationSelection: ReactNode = null;
  if (!!role && role.type ==='CUSTOM_ITEM') {
    roleSelection = (
      <View style={styles.selectionContainer}>
        <Text style={styles.selectionText}>
          <Text style={{fontWeight: "900"}}>{role.name}</Text>
          {' is a new role. Before adding, please make sure that the spelling and' +
           ' capitalization are correct.'}
        </Text>
      </View>
    );
  }
  if (!!organization && organization.type === 'CUSTOM_ITEM') {
    organizationSelection = (
      <View style={styles.selectionContainer}>
        <Text style={styles.selectionText}>
          <Text style={{fontWeight: "900"}}>{organization.name}</Text>
          {' is a new organization. Before adding, please make sure that the spelling and' +
           ' capitalization are correct.'}
        </Text>
      </View>
    );
  }

  return (
    <KeyboardAwareScrollView
      keyboardShouldPersistTaps="always"
    >
      <Field
        label="Role"
        name="role"
        placeholder="Search for roles (e.g. Accountant, Software Engineer)"
        component={AutocompleteInput}
        onQueryChange={onRoleQueryChangeThrottled}
        validate={required}
        containerStyle={styles.roleContainerStyle}
      />
      <Field
        label="Organization"
        name="organization"
        placeholder="Search for companies, clubs or sports teams"
        component={AutocompleteInput}
        onQueryChange={onOrganizationQueryChangeThrottled}
        validate={required}
      />
      <Field
        label="Start Date"
        name="startDate"
        androidMode={'spinner' as 'spinner'}
        mode={'date' as 'date'}
        component={ModalDatePicker}
        validate={required}
      />
      <Field
        label="End Date (Optional)"
        name="endDate"
        androidMode={'spinner' as 'spinner'}
        mode={'date' as 'date'}
        component={ModalDatePicker}
        maxDate={null}
      />
      { selection }
      { roleSelection }
      { organizationSelection }
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
        backgroundColor={Colors.HIVE_PRIMARY}
        disabled={!valid}
        loading={submitting}
        title={submitting ? null : "Add"}
        onPress={handleSubmit(onSubmit)}
      />
    </KeyboardAwareScrollView>
  );
}

const selector = formValueSelector('add-position');

const endDateAfterStartDate = (values: AddPositionFormData) => {
  const errors: FormErrors<AddPositionFormData> = {};
  if (!!values.endDate && !!values.startDate && values.endDate <= values.startDate) {
    errors.endDate = 'End date must be after start date'
  }
  return errors;
}

const AddPositionFormWithRedux =
  reduxForm<AddPositionFormData, FormP<AddPositionFormData>>({
    form: 'add-position',
    validate: endDateAfterStartDate,
  })(connect((state: RootState) => ({
    role: selector(state, 'role'),
    organization: selector(state, 'organization'),
    startDate: selector(state, 'startDate'),
    endDate: selector(state, 'endDate'),
  }))(AddPositionForm));

interface DispatchActions {
  addPosition: ActionCreator<ThunkAction<Promise<ActionTypes>, ProfileState, void>>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

interface Props extends DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class AddPositionView extends Component<Props> {
  ADD_POSITION_VIEW_IDENTIFIER = 'AddPositionView';

  static navigationOptions = {
    headerTitle: 'Add Position',
    headerStyle,
  }

  constructor(props: Props) {
    super(props);

    this.onSubmit = this.onSubmit.bind(this);
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.ADD_POSITION_VIEW_IDENTIFIER);
    });
  }

  private async onSubmit(values: AddPositionFormData) {
    try {
      const { role, organization, startDate, endDate } = values;
      if (role === null || organization === null) {
        const err = 'You have not selected a role or organization';
        Sentry.captureMessage(err, { level: 'error' });
        throw new Error(err);
      } else {
        const req = {
          roleId: role.type === 'ITEM' ? role.id as number : null,
          roleName: role.type === 'ITEM' ? null : role.name.trim(),
          organizationId: organization.type === 'ITEM' ? organization.id as number : null,
          organizationName: organization.type === 'ITEM' ? null : organization.name.trim(),
          startDate,
          endDate,
        };
        await this.props.addPosition(req);
        await this.props.infoToast(`Successfully added position "${role.name} @ ${organization.name}"`);
        await this.props.navigation.goBack();
      }
    } catch (e) {
      throw new SubmissionError({_error: e.errorMsg});
    }
  }

  render() {
    return <AddPositionFormWithRedux onSubmit={this.onSubmit} />;
  }
}

const styles = StyleSheet.create({
  selectionContainer: {
    padding: 10,
  },
  selectionText: {
    fontSize: 16,
  },
  roleContainerStyle: {
    zIndex: 2,
  },
});

export default connect(null, { addPosition, infoToast })(AddPositionView);
