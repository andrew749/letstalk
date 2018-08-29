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
} from 'redux-form';
import Sentry from 'sentry-expo';
import { KeyboardAwareScrollView } from 'react-native-keyboard-aware-scroll-view'
import { ThunkAction } from 'redux-thunk';
import _ from 'underscore';

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
} from '../../components';
import { DataItem, Select } from '../../components/AutocompleteInput';
import { headerStyle } from '../TopHeader';
import { AnalyticsHelper } from '../../services/analytics';
import Colors from '../../services/colors';
import autocompleteService from '../../services/autocomplete-service';
import requestToMatchService from '../../services/request-to-match-service';

interface AddSimpleTraitFormData {
  simpleTrait: Select;
}

const THROTTLE_TIME = 250; // ms

const onQueryChange = async (query: string, setData: (data: Array<DataItem>) => void) => {
  let res: Array<DataItem> = [];
  if (query !== '') res = await autocompleteService.autocompleteSimpleTrait(query, 10);
  setData(res);
}

const onQueryChangeThrottled = _.throttle(onQueryChange, THROTTLE_TIME);

const AddSimpleTraitForm: SFC<FormProps<AddSimpleTraitFormData> & AddSimpleTraitFormData> = props => {

  const {
    error,
    handleSubmit,
    onSubmit,
    reset,
    submitting,
    valid,
    simpleTrait,
  } = props;
  let selection: ReactNode = null;

  if (!simpleTrait) {
    selection = (
      <Text>Choose a trait from the dropdown above to add it to your list</Text>
    );
  } else if (simpleTrait.type === 'CUSTOM_ITEM') {
    selection = (
      <Text>
        {'Looks like you want to add a new trait '}
        <Text style={{fontWeight: "900"}}>{simpleTrait.name}</Text>
        {'. Make sure that the spelling and capitalization is correct before adding. '}
        {'Click the button below to add.'}
      </Text>
    );
  } else if (simpleTrait.type === 'ITEM') {
    selection = (
      <Text>
        {'Looks like you want to add '}
        <Text style={{fontWeight: "900"}}>{simpleTrait.name}</Text>
        {' to your list of traits. Click the button below to add.'}
      </Text>
    );
  }

  return (
    <KeyboardAwareScrollView
      keyboardShouldPersistTaps={true}
    >
      <Field
        label="Trait"
        name="simpleTrait"
        placeholder="Search for traits (e.g. cycling, climbing)"
        component={AutocompleteInput}
        onQueryChange={onQueryChangeThrottled}
        validate={required}
      >
      </Field>
      <View style={styles.selectionContainer}>
        <Text style={styles.selectionText}>
          { selection }
        </Text>
      </View>
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

const selector = formValueSelector('add-simple-trait');

const AddSimpleTraitFormWithRedux =
  reduxForm<AddSimpleTraitFormData, FormP<AddSimpleTraitFormData>>({
    form: 'add-simple-trait',
  })(connect((state: RootState) => ({
    simpleTrait: selector(state, 'simpleTrait'),
  }))(AddSimpleTraitForm));

interface DispatchActions {
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

interface Props extends DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class AddSimpleTraitView extends Component<Props> {
  ADD_SIMPLE_TRAIT_VIEW_IDENTIFIER = 'AddSimpleTraitView';

  static navigationOptions = {
    headerTitle: 'Add Trait',
    headerStyle,
  }

  constructor(props: Props) {
    super(props);

    this.onSubmit = this.onSubmit.bind(this);
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.ADD_SIMPLE_TRAIT_VIEW_IDENTIFIER);
    });
  }

  private async onSubmit(values: AddSimpleTraitFormData) {
    try {
      const { simpleTrait } = values;
      let err = '';
      if (simpleTrait === null) {
        err = 'You have not selected a trait';
        Sentry.captureMessage(err, { level: 'error' });
        throw new Error(err);
      } else {
        if (simpleTrait.type === 'CUSTOM_ITEM') {
          const newTraitName = simpleTrait.name.trim();
          await requestToMatchService.addUserSimpleTraitByName(newTraitName);
          await this.props.infoToast(`Successfully added new trait "${newTraitName}"`);
        } else if (simpleTrait.type === 'ITEM') {
          await requestToMatchService.addUserSimpleTraitById(simpleTrait.id as number);
          await this.props.infoToast(`Successfully added trait "${simpleTrait.name}"`);
        } else {
          err = 'Invalid item type selected';
          Sentry.captureMessage(err, { level: 'error' });
          throw new Error(err);
        }
        await this.props.navigation.goBack();
      }
    } catch (e) {
      throw new SubmissionError({_error: e.errorMsg});
    }
  }

  render() {
    return <AddSimpleTraitFormWithRedux onSubmit={this.onSubmit} />;
  }
}

const styles = StyleSheet.create({
  selectionContainer: {
    padding: 10,
  },
  selectionText: {
    fontSize: 16,
  }
});

export default connect(null, { infoToast })(AddSimpleTraitView);
