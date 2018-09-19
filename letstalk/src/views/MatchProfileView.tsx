import React, { Component } from 'react';
import {
  ActivityIndicator,
  AppRegistry,
  Button as ReactNativeButton,
  Dimensions,
  FlatList,
  Image,
  Linking,
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

import { infoToast, errorToast } from '../redux/toast';
import auth from '../services/auth';
import {fbLogin} from '../services/fb';
import { FloatingButton, Card, Header } from '../components';
import Loading from './Loading';
import { genderIdToString } from '../models/user';
import { RootState } from '../redux';
import { State as MatchProfileState, fetchMatchProfile } from '../redux/match-profile/reducer';
import { ActionTypes as MatchProfileActionTypes } from '../redux/match-profile/actions';
import { programById, sequenceById } from '../models/cohort';
import { AnalyticsHelper } from '../services/analytics';
import { ProfileAvatar } from '../components';
import Colors from '../services/colors';
import { headerStyle, headerTitleStyle, headerTintColor } from './TopHeader';
import {
  CohortInfo,
  PersonalInfo,
  UserPositions,
  UserSimpleTraits,
  styles,
} from './profile-components/ProfileComponents';
import {
  RelationshipTypes,
} from '../models/profile';
import {
  IntentTypes,
  ConnectionIntent,
} from '../models/connection';
import { ActionTypes as BootstrapActionTypes } from '../redux/bootstrap/actions';
import requestToMatchService from '../services/request-to-match-service';
import {
  State as BootstrapState,
  fetchBootstrap,
} from '../redux/bootstrap/reducer';

const SCREEN_WIDTH = Dimensions.get('window').width;

interface DispatchActions {
  fetchMatchProfile: ActionCreator<ThunkAction<Promise<MatchProfileActionTypes>, MatchProfileState, void>>;
  fetchBootstrap: ActionCreator<ThunkAction<Promise<BootstrapActionTypes>, BootstrapState, void>>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

interface NavigationParams {
  readonly userId: number,
  readonly connectionIntent?: ConnectionIntent,
}

interface Props extends MatchProfileState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction & NavigationParams>;
}

class MatchProfileView extends Component<Props> {
  PROFILE_VIEW_IDENTIFIER = "MatchProfileView";

  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: 'Profile',
    headerStyle, 
    headerTitleStyle, 
    headerTintColor 
  })

  constructor(props: Props) {
    super(props);

    this.load = this.load.bind(this);
    this.renderBody = this.renderBody.bind(this);
  }

  async componentDidMount() {
    AnalyticsHelper.getInstance().recordPage(this.PROFILE_VIEW_IDENTIFIER);
    this.load();
  }

  private async load() {
    const userId = this.props.navigation.getParam('userId', 0);
    await this.props.fetchMatchProfile(userId);
  }

  private renderContactInfo() {
    const {
      email,
      phoneNumber,
      fbId,
      fbLink,
    } = this.props.profile;

    const buildItem = (label: string, value: string, link: string) => {
      const onPress = () => Linking.openURL(link);
      return (
        <TouchableOpacity key={label} onPress={onPress} style={styles.listItem}>
          <Text style={styles.label}>{label}: </Text>
          <Text style={styles.value}>{value}</Text>
        </TouchableOpacity>
      );
    };

    const buildItems = (name_values: Array<[string, string, string]>) => {
      const items = name_values.map(([label, value, link]) => {
        return buildItem(label, value, link);
      });
      return items;
    };

    const emailLink = 'mailto:' + email;
    const contactInfos: Array<[string, string, string]> = [
      ['Email', email, emailLink],
    ];

    if (!!phoneNumber) {
      const smsLink = 'sms:' + phoneNumber;
      contactInfos.push(['Phone', phoneNumber, smsLink]);
    }
    const contactItems = buildItems(contactInfos);

    if (fbLink) {
      contactItems.push(
        <TouchableOpacity style={styles.listItem} onPress={() => Linking.openURL(fbLink)}>
          <MaterialIcons name="face" size={24} />
          <Text style={styles.label}>Facebook</Text>
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

  private renderRequestButton(connectionIntent: ConnectionIntent) {
    const { relationshipType } = this.props.profile;
    let title = 'Request to connect';
    let color = Colors.HIVE_PRIMARY;
    let disabled = false;
    let onPress = () => {
      this.props.navigation.navigate('RequestConnection', {
        profile: this.props.profile,
        connectionIntent,
      });
    };

    const userId = this.props.profile.userId;

    switch (relationshipType) {
      case RelationshipTypes.CONNECTED:
        // Already connected so we don't show button
        return null;
      case RelationshipTypes.YOU_REQUESTED:
        title = 'Request already sent';
        disabled = true;
        onPress = null;
        break;
      case RelationshipTypes.THEY_REQUESTED:
        title = 'Accept request';
        color = Colors.GREEN;
        onPress = async () => {
          try {
            await requestToMatchService.acceptConnection(userId);
            this.props.infoToast('Accepted connection');
            this.props.fetchMatchProfile(userId);
            this.props.fetchBootstrap();
          } catch(e) {
            await this.props.errorToast(e.errorMsg);
            throw e;
          }
        }
        break;
    }

    return (
      <FloatingButton
        title={title}
        onPress={onPress}
        buttonStyle={{ backgroundColor: color }}
        disabled={disabled}
      />
    );
  }

  private renderBody() {
    const { navigate } = this.props.navigation;

    const isConnected = this.props.profile.relationshipType === RelationshipTypes.CONNECTED;

    let userId;
    if (this.props.profile) {
      userId = this.props.profile.userId.toString();
    }

    const connectionIntent = this.props.navigation.getParam('connectionIntent', null)
    const extraStyle = !!connectionIntent && !isConnected ? { paddingBottom: 80 } : null;

    return (
      <View>
        <ScrollView contentContainerStyle={styles.container}>
          <View style={[styles.contentContainer, extraStyle]} >
            <ProfileAvatar userId={userId} xlarge containerStyle={styles.profilePicture} />
            <PersonalInfo
              {...this.props.profile}
              navigation={this.props.navigation}
              showConnectedBadge={isConnected}
            />
            <CohortInfo
              programId={this.props.profile.programId}
              sequenceId={this.props.profile.sequenceId}
              gradYear={this.props.profile.gradYear}
              navigation={this.props.navigation}
            />
            {isConnected ? this.renderContactInfo() : null}
            <UserPositions
              userPositions={this.props.profile.userPositions}
              navigation={this.props.navigation}
            />
            <UserSimpleTraits
              userSimpleTraits={this.props.profile.userSimpleTraits}
              navigation={this.props.navigation}
            />
          </View>
        </ScrollView>
        { connectionIntent && this.renderRequestButton(connectionIntent) }
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
  ({ matchProfile }: RootState) => matchProfile,
  { infoToast, errorToast, fetchBootstrap, fetchMatchProfile },
)(MatchProfileView);
