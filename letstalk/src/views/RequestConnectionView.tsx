import React, {Component} from 'react';
import { NavigationScreenProp, NavigationStackAction, NavigationActions } from 'react-navigation';
import { View, StyleSheet } from 'react-native';
import { reduxForm, Field, SubmissionError } from 'redux-form';
import { FormValidationMessage, Text } from 'react-native-elements';
import { ActionButton, FormP, FormProps, LabeledFormInput } from '../components';
import { AnalyticsHelper } from '../services';
import { infoToast, errorToast } from '../redux/toast';
import { headerStyle, headerTitleStyle, headerTintColor } from './TopHeader';
import {ActionCreator, connect, Dispatch} from 'react-redux';
import { RootState } from '../redux';
import { KeyboardAwareScrollView } from 'react-native-keyboard-aware-scroll-view'
import { ThunkAction } from "redux-thunk";

import { State as MatchProfileState, fetchMatchProfile } from '../redux/match-profile/reducer';
import { ActionTypes } from '../redux/match-profile/actions';
import { MatchProfileData } from '../models/profile';
import { ConnectionIntent } from '../models/connection';
import Colors from '../services/colors';
import requestToMatchService from '../services/request-to-match-service';

interface NavigationParams {
  readonly profile: MatchProfileData,
  readonly connectionIntent: ConnectionIntent,
}

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction & NavigationParams>;
  fetchMatchProfile: ActionCreator<ThunkAction<Promise<ActionTypes>, MatchProfileState, void>>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

interface State {
  refreshing: boolean;
}

interface RequestConnectionFormData {
  message: string,
}

interface RequestConnectionFormProps extends FormProps<RequestConnectionFormData> {
  profile: MatchProfileData;
}

const RequestConnectionForm: React.SFC<RequestConnectionFormProps> = props => {
  const { profile } = props;
  const { error, handleSubmit, onSubmit, submitting, valid, pristine } = props;
  const onSubmitWithReset = async (values: RequestConnectionFormData): Promise<void> => {
    await onSubmit(values);
  };
  return (
    <View>
      <Field
        label="Message (optional)"
        name="message"
        placeholder={'Tell ' + profile.firstName + ' why you want to connect'}
        component={LabeledFormInput as "input" & typeof LabeledFormInput}
        multiline={true}
        numberOfLines={10}
        inputStyle={{width: "100%"}}
      />
      {error && <FormValidationMessage>{error}</FormValidationMessage>}
      <ActionButton
        buttonStyle={{backgroundColor: Colors.HIVE_PRIMARY}}
        textStyle={{color: Colors.WHITE}}
        loading={submitting}
        title={submitting ? null : "Send request"}
        onPress={handleSubmit(onSubmitWithReset)}
      />
    </View>
  );
}

const RequestConnectionFormWithRedux = reduxForm<
  RequestConnectionFormData, FormP<RequestConnectionFormData>
>({
  form: 'requestConnection',
})(RequestConnectionForm as any);

export class RequestConnectionView extends Component<Props, State> {
  REQUEST_CONNECTION_VIEW_IDENTIFIER = "RequestConnectionView";

  static navigationOptions = {
    headerTitle: 'Request Connection',
    headerStyle, 
    headerTitleStyle, 
    headerTintColor
  }

  constructor(props: Props) {
    super(props);
    this.onSubmit = this.onSubmit.bind(this);
    this.state = { refreshing: false };
  }

  async componentDidMount() {
    AnalyticsHelper.getInstance().recordPage(this.REQUEST_CONNECTION_VIEW_IDENTIFIER);
  }

  async onSubmit(values: RequestConnectionFormData) {
    const { userId } = this.props.navigation.getParam('profile');
    const connectionIntent = this.props.navigation.getParam('connectionIntent')
    const { message } = values;
    try {
      await requestToMatchService.requestConnection({...connectionIntent, userId, message});
      await this.props.infoToast('Sent request');
    } catch(e) {
      await this.props.errorToast(e.errorMsg);
      throw new SubmissionError({_error: e.errorMsg});
    }
    this.props.fetchMatchProfile(userId);
    this.props.navigation.goBack();
}

  render() {
    const profile = this.props.navigation.getParam('profile');
    const connectionInent = this.props.navigation.getParam('connectionIntent')
    const { firstName, lastName } = profile;
    return (
      <KeyboardAwareScrollView style={styles.container}>
        <Text style={styles.headline}>
          <Text>{'Send a request to connect with '}</Text>
          <Text style={styles.boldText}>{profile.firstName} {profile.lastName}</Text>
          <Text>{'.'}</Text>
        </Text>
        // @ts-ignore - Couldn't find a way to add this to the props
        <RequestConnectionFormWithRedux
          onSubmit={this.onSubmit}
          profile={profile}
        />
      </KeyboardAwareScrollView>
    );
  }
}

export default connect(null, { infoToast, errorToast, fetchMatchProfile })(RequestConnectionView as any);

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  centeredContainer: {
    flex: 1,
    justifyContent: 'center'
  },
  headline: {
    fontSize: 18,
    textAlign: 'left',
    margin: 20,
  },
  boldText: {
    fontWeight: '900',
  },
});
