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
  ModalPicker,
} from '../../components';
import { DataItem, Select } from '../../components/AutocompleteInput';
import { headerStyle, headerTitleStyle, headerTintColor } from '../TopHeader';
import { AnalyticsHelper } from '../../services/analytics';
import Colors from '../../services/colors';
import autocompleteService from '../../services/autocomplete-service';
import { ActionTypes as ProfileActionTypes } from '../../redux/profile/actions';
import {
  State as ProfileState,
  addGroup,
} from '../../redux/profile/reducer';
import { groupNamesByGroupId } from '../../models/group';

interface AddGroupFormData {
  groupId: string;
}

const AddGroupForm: SFC<FormProps<{}>> = props => {
  const {
    error,
    handleSubmit,
    onSubmit,
    reset,
    submitting,
    valid,
  } = props;

  const groupItems = groupNamesByGroupId.mapEntries(([value, label]) => {
    return [value, <Picker.Item key={value} label={label} value={value}/>];
  }).toList().toJS();

  return (
    <KeyboardAwareScrollView
      keyboardShouldPersistTaps="always"
    >
      <Field
        label="Group"
        name="groupId"
        containerStyle={styles.pickerContainer}
        component={ModalPicker}
        validate={required}
      >
        {groupItems}
      </Field>
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
        backgroundColor={Colors.HIVE_PRIMARY}
        disabled={!valid}
        loading={submitting}
        title={submitting ? null : "Join"}
        onPress={handleSubmit(onSubmit)}
      />
    </KeyboardAwareScrollView>
  );
}

const selector = formValueSelector('add-position');

const AddGroupFormWithRedux =
  reduxForm<AddGroupFormData, FormP<AddGroupFormData>>({
    form: 'add-position',
  })(AddGroupForm);

interface DispatchActions {
  addGroup: ActionCreator<ThunkAction<Promise<ProfileActionTypes>, ProfileState, void>>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

interface Props extends DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class AddGroupView extends Component<Props> {
  ADD_GROUP_VIEW_IDENTIFIER = 'AddGroupView';

  static navigationOptions = {
    headerTitle: 'Add Group',
    headerStyle,
    headerTitleStyle,
    headerTintColor
  }

  constructor(props: Props) {
    super(props);

    this.onSubmit = this.onSubmit.bind(this);
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.ADD_GROUP_VIEW_IDENTIFIER);
    });
  }

  private async onSubmit(values: AddGroupFormData) {
    try {
      const { groupId } = values;
      const groupName = groupNamesByGroupId.get(groupId);
      await this.props.addGroup({ groupName, groupId });
      await this.props.infoToast(`Successfully joined the ${groupName} group`);
      await this.props.navigation.goBack();
    } catch (e) {
      throw new SubmissionError({_error: e.errorMsg});
    }
  }

  render() {
    return <AddGroupFormWithRedux onSubmit={this.onSubmit} />;
  }
}

const styles = StyleSheet.create({
  selectionContainer: {
    padding: 10,
  },
  selectionText: {
    fontSize: 16,

  },
  textInputContainerStyle: {
    borderRadius: 5
  },
  roleStyle: {
    zIndex: 2,
  },
  pickerContainer: {
    marginHorizontal: 15
  },
});

export default connect(null, { addGroup, infoToast })(AddGroupView);
