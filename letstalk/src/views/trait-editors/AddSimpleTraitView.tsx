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
import { required } from '../../validators';
import Loading from '../Loading';
import {
  ActionButton,
  AutocompleteInput,
  FormProps,
  FormP,
  Header,
} from '../../components';
import { Select } from '../../components/AutocompleteInput';
import { headerStyle } from '../TopHeader';
import { AnalyticsHelper } from '../../services/analytics';
import Colors from '../../services/colors';
import autocompleteService from '../../services/autocomplete-service';

interface AddSimpleTraitFormData {
  simpleTrait: Select;
}

const AddSimpleTraitForm: SFC<FormProps<AddSimpleTraitFormData>> = props => {
  const onQueryChange = async (query: string) => {
    if (query === '') return [];
    const res = await autocompleteService.autocompleteSimpleTrait(query, 10);
    return res;
  }

  const {
    error,
    handleSubmit,
    onSubmit,
    reset,
    submitting,
    valid,
  } = props;
  return (
    <KeyboardAwareScrollView
      keyboardShouldPersistTaps={true}
    >
      <Field
        label="Trait"
        name="simpleTrait"
        component={AutocompleteInput}
        allowCustom={true}
        onQueryChange={onQueryChange}
        validate={required}
      >
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

const AddSimpleTraitFormWithRedux =
  reduxForm<AddSimpleTraitFormData, FormP<AddSimpleTraitFormData>>({
    form: 'add-simple-trait',
  })(AddSimpleTraitForm);

interface DispatchActions {
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

// TODO: Maybe take current cohort info as a prop instead of pulling from redux. Makes this more
// reusable.
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
    this.renderBody = this.renderBody.bind(this);
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.ADD_SIMPLE_TRAIT_VIEW_IDENTIFIER);
    });
  }

  private async onSubmit(values: AddSimpleTraitFormData) {
    try {
      const { simpleTrait } = values;
      console.log(simpleTrait);
      // TODO: Fill in
    } catch (e) {
      throw new SubmissionError({_error: e.errorMsg});
    }
  }

  private renderBody() {
    return <AddSimpleTraitFormWithRedux onSubmit={this.onSubmit} />;
  }

  render() {
    return this.renderBody();
  }
}

export default connect(null, { infoToast })(AddSimpleTraitView);
