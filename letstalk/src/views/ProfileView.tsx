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
import { Button, Card, FloatingButton, Header } from '../components';
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
import { ProfileAvatar } from '../components';
import Colors from '../services/colors';
import TopHeader, { headerStyle } from './TopHeader';
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
};

const initialState: State = {
  showAllPositions: false,
  showAllSimpleTraits: false,
};

class ProfileView extends Component<Props, State> {
  PROFILE_VIEW_IDENTIFIER = "ProfileView";

  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: <TopHeader navigation={navigation} />,
    headerStyle,
  });

  constructor(props: Props) {
    super(props);

    this.state = initialState;

    this.onLogoutPress = this.onLogoutPress.bind(this);
    this.onChangePasswordPress = this.onChangePasswordPress.bind(this);
    this.onEditTraitsButtonPress = this.onEditTraitsButtonPress.bind(this);
    this.load = this.load.bind(this);
    this.renderBody = this.renderBody.bind(this);
  }

  private async onLogoutPress() {
    try {
      await auth.logout();
    } catch (error) { }
    await this.props.navigation.dispatch(NavigationActions.reset({
      index: 0,
      key: null,
      actions: [NavigationActions.navigate({ routeName: 'Login' })]
    }));
  }

  private async onChangePasswordPress() {
    try {
      await auth.forgotPassword(this.props.profile.email);
      await this.props.infoToast("Sent an email with reset instructions.");
    } catch(e) {
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

  private renderContactInfo(email: string, fbId: string, fbLink: string, phoneNumber: string) {
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
          <MaterialIcons name="face" size={24} />
          <Text style={styles.label}>Facebook</Text>
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
    const { navigate } = this.props.navigation;

    const {
      userPositions,
      userSimpleTraits,
    } = this.props.profile;

    const {
      firstName,
      lastName,
      gender,
      email,
      birthdate,
      phoneNumber,
      fbId,
      fbLink,
      bio,
      hometown,
      secret,
    } = this.props.profile;

    let userId;
    if (this.props.profile) {
      userId = this.props.profile.userId.toString();
    }

    return (
      <View>
        <ScrollView contentContainerStyle={styles.container}>
          <View style={styles.contentContainer} >
            <ProfileAvatar userId={userId} xlarge containerStyle={styles.profilePicture} />
            <PersonalInfo
              {...this.props.profile}
              navigation={this.props.navigation}
              allowQrCode={true}
            />
            <CohortInfo
              programId={this.props.profile.programId}
              sequenceId={this.props.profile.sequenceId}
              gradYear={this.props.profile.gradYear}
              navigation={this.props.navigation}
              allowEditing={true}
            />
            {this.renderContactInfo(email, fbId, fbLink, phoneNumber)}
            <UserPositions
              userPositions={userPositions}
              navigation={this.props.navigation}
              allowEditing={true}
              removePosition={this.props.removePosition}
              errorToast={this.props.errorToast}
            />
            <UserSimpleTraits
              userSimpleTraits={userSimpleTraits}
              navigation={this.props.navigation}
              allowEditing={true}
              removeSimpleTrait={this.props.removeSimpleTrait}
              errorToast={this.props.errorToast}
            />
            <View style={styles.sectionContainer}>
              <Text style={styles.sectionHeader}>Account Actions</Text>
              <View style={{ alignItems: 'center' }}>
                <Button
                  buttonStyle={styles.changePassButton}
                  onPress={this.onChangePasswordPress}
                  title='Change Password'
                />
                <Button
                  buttonStyle={styles.logoutButton}
                  textStyle={styles.logoutButtonText}
                  onPress={this.onLogoutPress}
                  title='Logout'
                />
              </View>
            </View>
          </View>
        </ScrollView>
        <FloatingButton title="Edit Profile" onPress={() => {
          navigate('EditProfileSelector', { profile: this.props.profile });
        }} />
        <AllFilterableModals
          onSelectSuccess={() => {
            this.props.navigation.navigate({ routeName: 'Requests' });
          }}
        />
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

export default connect(
  ({ profile }: RootState) => profile,
  { fetchProfile, removePosition, removeSimpleTrait, infoToast, errorToast },
)(ProfileView);
