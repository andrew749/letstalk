import React, { Component, ReactNode } from 'react';
import {
  ActivityIndicator,
  AppRegistry,
  Alert,
  Button as ReactNativeButton,
  Dimensions,
  FlatList,
  Image,
  Linking,
  Platform,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';
import {
  formValueSelector,
  reduxForm,
  Field,
  InjectedFormProps,
  SubmissionError,
} from 'redux-form';
import { Icon } from 'react-native-elements';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import { bindActionCreators } from 'redux'
import {
  NavigationScreenDetails,
  NavigationScreenProp,
  NavigationStackAction,
  NavigationActions,
} from 'react-navigation';
import { MaterialIcons } from '@expo/vector-icons';
import Immutable from 'immutable';
import Moment from 'moment';

import auth from '../services/auth';
import { infoToast, errorToast } from '../redux/toast';
import {fbLogin} from '../services/fb';
import { ActionButton, Button} from '../components';
import Loading from './Loading';
import { genderIdToString } from '../models/user';
import { RootState } from '../redux';
import {
  State as ProfileState,
  fetchProfile,
  removePosition,
  removeSimpleTrait,
} from '../redux/profile/reducer';
import { ActionTypes } from '../redux/profile/actions';
import { AnalyticsHelper } from '../services/analytics';
import { 
  FormP, FormProps, ProfileAvatarEditableFormElement } from '../components';
import photoService, {PhotoResult} from '../services/photo_service';
import Colors from '../services/colors';
import TopHeader, { headerStyle, headerTitleStyle, headerTintColor } from './TopHeader';
import AllFilterableModals from './AllFilterableModals';
import { UserPosition } from '../models/position';
import { UserSimpleTrait } from '../models/simple-trait';
import {
  CohortInfo,
  PersonalInfo,
  UserPositions,
  UserSimpleTraits,
  styles,
} from './profile-components/ProfileComponents';

const SCREEN_WIDTH = Dimensions.get('window').width;

type EditFormComponentProps = FormProps<PhotoResult> & PhotoResult;

class EditForm extends Component<EditFormComponentProps, State> {
  constructor(props: EditFormComponentProps) {
    super(props)
  }

  render() {
    const {
      error,
      handleSubmit,
      onSubmit,
      reset,
      submitting,
      valid,
      initialValues: { uri }
    } = this.props;

    return (
      <View>
        <Field
          name="profilePic"
          component={ProfileAvatarEditableFormElement}
          onChange={onSubmit}
          uri={uri}
        />
      </View>
    );
  }
}

const EditFormWithReduxBuilder = (initialValues: PhotoResult) => {
  return reduxForm<PhotoResult, FormP<PhotoResult>>({
    form: 'profile-pic-edit',
    initialValues,
  })(connect()(EditForm));
}

interface DispatchActions {
  fetchProfile: ActionCreator<ThunkAction<Promise<ActionTypes>, ProfileState, void>>;
  removePosition: ActionCreator<ThunkAction<Promise<ActionTypes>, ProfileState, void>>;
  removeSimpleTrait: ActionCreator<ThunkAction<Promise<ActionTypes>, ProfileState, void>>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

interface Props extends ProfileState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

interface State {
  readonly showAllPositions: boolean;
  readonly showAllSimpleTraits: boolean;
  readonly forgotPasswordRequest: boolean;
  readonly logoutRequest: boolean;
};

const initialState: State = {
  showAllPositions: false,
  showAllSimpleTraits: false,
  forgotPasswordRequest: false,
  logoutRequest: false,
};

class ProfileView extends Component<Props, State> {
  PROFILE_VIEW_IDENTIFIER = "ProfileView";

  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: <TopHeader navigation={navigation} />,
    headerStyle, 
    headerTitleStyle, 
    headerTintColor
  });

  constructor(props: Props) {
    super(props);

    this.state = initialState;

    this.onSubmit = this.onSubmit.bind(this);
    this.onLogoutPress = this.onLogoutPress.bind(this);
    this.onChangePasswordPress = this.onChangePasswordPress.bind(this);
    this.onEditTraitsButtonPress = this.onEditTraitsButtonPress.bind(this);
    this.load = this.load.bind(this);
    this.renderBody = this.renderBody.bind(this);
  }

  private async onSubmit(profilePic: PhotoResult) {
    try {
      if (profilePic && profilePic.uri) {
        let res = await photoService.uploadProfilePhoto(profilePic.uri);
      }
    } catch(e) {
      throw new SubmissionError({_error: e.errorMsg});
    }
  }

  private async onLogoutPress() {
    const onLogout = async () => {
      try {
        await auth.logout();
      } catch (error) { }
      await this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        key: null,
        actions: [NavigationActions.navigate({ routeName: 'Login' })]
      }));
    }

    Alert.alert(
      'Logout',
      `Are you sure you want to logout?`,
      [
        {text: 'Cancel', onPress: () => null, style: 'cancel'},
        {text: 'Logout', onPress: onLogout, style: 'default'},
      ],
    );

  }

  private async onChangePasswordPress() {
    try {
      this.setState({...this.state, forgotPasswordRequest: true});
      let resp = await auth.forgotPassword(this.props.profile.email);
      let resp2 = await this.props.infoToast("Check your email for reset instructions!");
      if (resp && resp2) {
        this.setState({...this.state, forgotPasswordRequest: false});
      }
    } catch(e) {
      this.setState({...this.state, forgotPasswordRequest: false});
      await this.props.errorToast(e.errorMsg);
    }
  }

  private async onEditTraitsButtonPress() {
    const { profile } = this.props;
    await this.props.navigation.navigate('EditTraitsSelector', { profile });
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.PROFILE_VIEW_IDENTIFIER);
    });
    this.load();
  }

  private async load() {
    await this.props.fetchProfile();
  }

  private renderContactInfo() {
    const {
      email,
      phoneNumber,
      fbId,
      fbLink,
    } = this.props.profile;
    const buildItem = (label: string, value: string) => {
      return (
        <View key={label} style={styles.listItem}>
          <Text style={styles.label}>{label}: </Text>
          <Text style={styles.value}>{value}</Text>
        </View>
      );
    };

    const buildItems = (name_values: Array<[string, string]>) => {
      const items = name_values.filter(x => x != null).map(([label, value]) => {
        return buildItem(label, value);
      });
      return items;
    };

    const contactItems = buildItems([
      phoneNumber ? ['Phone', phoneNumber] : undefined,
      ['Email', email]
    ]);


    if (fbLink !== null) {
      contactItems.push(
        <TouchableOpacity
          key={'facebook'}
          style={styles.listItem}
          onPress={() => Linking.openURL(fbLink)}
        >
          <Text style={styles.value}>Visit Facebook Profile</Text>
        </TouchableOpacity>
      );
    } else {
      // link fb profile
      contactItems.push(
        <TouchableOpacity
          key={'facebook'}
          style={styles.listItem}
          onPress={async () => {
            await auth.linkFB();
            await this.props.fetchProfile();
          }}
        >
          <Text style={styles.value}>Link your Facebook account</Text>
        </TouchableOpacity>
      );
    }

    return (
      <View style={styles.sectionContainer}>
        <Text style={styles.sectionHeader}>Contact Info</Text>
        {contactItems}
      </View>
    )
  }

  private renderBody() {
    let profilePic;
    if (this.props.profile) {
      profilePic = this.props.profile.profilePic;
    }
    const EditFormWithRedux = EditFormWithReduxBuilder({
      uri: profilePic, 
      data: null,
    });

    return (
      <View>
        <ScrollView contentContainerStyle={[styles.container, { paddingBottom: 65 }]}>
          <View style={styles.contentContainer} >
            <EditFormWithRedux onSubmit={this.onSubmit} />
            <PersonalInfo
              {...this.props.profile}
              navigation={this.props.navigation}
              allowQrCode={true}
              allowEditing={true}
            />
            <CohortInfo
              programId={this.props.profile.programId}
              sequenceId={this.props.profile.sequenceId}
              gradYear={this.props.profile.gradYear}
              navigation={this.props.navigation}
              allowEditing={true}
            />
            {this.renderContactInfo()}
            <UserPositions
              userPositions={this.props.profile.userPositions}
              navigation={this.props.navigation}
              allowEditing={true}
              removePosition={this.props.removePosition}
              errorToast={this.props.errorToast}
            />
            <UserSimpleTraits
              userSimpleTraits={this.props.profile.userSimpleTraits}
              navigation={this.props.navigation}
              allowEditing={true}
              removeSimpleTrait={this.props.removeSimpleTrait}
              errorToast={this.props.errorToast}
            />
            <View style={styles.sectionContainer}>
              <Text style={styles.sectionHeader}>Account Actions</Text>
              <View style={{ alignItems: 'center' }}>
                <ActionButton
                  backgroundColor={Colors.WHITE}
                  buttonStyle={[styles.changePassButton, styles.profileActionButton,]}
                  textStyle={[styles.changePassButtonText, styles.buttonText]}
                  loading={this.state.forgotPasswordRequest}
                  title={this.state.forgotPasswordRequest ? null : "Change Password"}
                  onPress={this.onChangePasswordPress}
                />
                <ActionButton
                  backgroundColor={Colors.HIVE_SUBDUED}
                  buttonStyle={[styles.logoutButton, styles.profileActionButton]}
                  textStyle={[styles.logoutButtonText, styles.buttonText]}
                  loading={this.state.logoutRequest}
                  title={this.state.logoutRequest ? null : "Logout"}
                  onPress={this.onLogoutPress}
                />
              </View>
            </View>
          </View>
        </ScrollView>
      </View>
    );
  }

  render() {
    const {
      state,
      errorMsg,
      errorType,
    } = this.props.fetchState;
    return (
      <View style={{flex: 1}}>
        <Loading
          state={state}
          errorMsg={errorMsg}
          errorType={errorType}
          load={this.load}
          renderBody={this.renderBody}
          navigation={this.props.navigation}
        />
        <AllFilterableModals
          onSelectSuccess={() => {
            this.props.navigation.navigate({ routeName: 'Explore' });
          }}
        />
      </View>
    );
  }
}

export default connect(
  ({ profile }: RootState) => profile,
  { fetchProfile, removePosition, removeSimpleTrait, infoToast, errorToast },
)(ProfileView);
