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
import { programById, sequenceById } from '../models/cohort';
import { AnalyticsHelper } from '../services/analytics';
import { ProfileAvatar } from '../components';
import Colors from '../services/colors';
import TopHeader, { headerStyle } from './TopHeader';
import AllFilterableModals from './AllFilterableModals';
import { UserPosition } from '../models/position';
import { UserSimpleTrait } from '../models/simple-trait';

const SCREEN_WIDTH = Dimensions.get('window').width;

const MAX_NUMBER_POSITIONS_SHOWN = 3;
const MAX_NUMBER_SIMPLE_TRAITS_SHOWN = 5;

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
    this.renderPosition = this.renderPosition.bind(this);
    this.renderSimpleTrait = this.renderSimpleTrait.bind(this);
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

  private renderProfile(bio: string | null) {
    const bioStr = bio === null ? 'Add bio by editing profile' : bio;
    return (
      <View style={styles.sectionContainer}>
        <Text style={styles.description}>{ bioStr }</Text>
      </View>
    )
  }

  private renderCohortInfo() {
    const {
      programId,
      gradYear,
      sequenceId,
    } = this.props.profile;
    const program = programById(programId);
    const changeCohort = () => this.props.navigation.navigate('ChangeCohort', {
      programId,
      gradYear,
      sequenceId,
    });
    return (
      <View style={styles.sectionContainer}>
        <TouchableOpacity onPress={changeCohort} style={styles.addTraitButton}>
          <MaterialIcons name="edit" size={32} color={Colors.HIVE_PRIMARY} />
        </TouchableOpacity>
        <Text style={styles.sectionHeader}>Cohort</Text>
        <Text style={styles.cohortText}>{ program + ', ' + gradYear }</Text>
      </View>
    );
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

  private static renderShowLessMore(isShown: boolean, show: () => void, hide: () => void) {
    if (isShown) {
      return (
        <TouchableOpacity style={styles.listItem} onPress={hide}>
          <Text style={styles.value}>Show less...</Text>
        </TouchableOpacity>
      );
    } else {
      return (
        <TouchableOpacity style={styles.listItem} onPress={show}>
          <Text style={styles.value}>Show more...</Text>
        </TouchableOpacity>
      );
    }
  }

  private renderPosition(pos: UserPosition) {
    const dateFmt = "MMM YYYY";
    const until = !pos.endDate ? 'present' : Moment(pos.endDate).format(dateFmt);
    const frm = Moment(pos.startDate).format(dateFmt);

    const onRemoveAccept = async () => {
      try {
        await this.props.removePosition(pos.id);
      } catch (e) {
        await this.props.errorToast(e.errorMsg);
      }
    }

    const onRemovePress = () => {
      Alert.alert(
        'Remove Position',
        `Are you sure you want to remove your position as ${pos.roleName} at ${pos.organizationName}?`,
        [
          {text: 'Cancel', onPress: () => null, style: 'cancel'},
          {text: 'Remove', onPress: onRemoveAccept, style: 'destructive'},
        ],
      );
    }

    return (
      <View key={ pos.id } style={styles.positionContainer}>
        <Text style={styles.positionText}>
          <Text style={styles.positionBold}>{ pos.roleName }</Text>
          <Text> @ </Text>
          <Text style={styles.positionBold}>{ pos.organizationName }</Text>
          <Text>{'\n'}({ frm } - { until })</Text>
        </Text>
        <TouchableOpacity style={styles.traitDelete} onPress={onRemovePress}>
          <MaterialIcons color={Colors.WHITE} name="close" size={18} />
        </TouchableOpacity>
      </View>
    );
  }

  private renderPositions() {
    let { userPositions } = this.props.profile;
    const { showAllPositions } = this.state;

    const addPosition = () => this.props.navigation.navigate('AddPosition');

    let bottomAction: ReactNode = null;
    if (userPositions.isEmpty()) {
      bottomAction = [
        <Text key={'text'} style={styles.noTraitText}>You don't have any positions</Text>,
        <Button
          key={'button'}
          buttonStyle={styles.noTraitButton}
          title="Add position"
          onPress={addPosition}
          color={Colors.HIVE_ACCENT}
        />,
      ];
    } else if (userPositions.size > MAX_NUMBER_POSITIONS_SHOWN) {
      bottomAction = ProfileView.renderShowLessMore(showAllPositions,
        () => this.setState({ showAllPositions: true }),
        () => this.setState({ showAllPositions: false }));
    }

    userPositions = userPositions.sortBy(pos => pos.startDate).reverse().toList();
    if (!showAllPositions) {
      userPositions = userPositions.take(MAX_NUMBER_POSITIONS_SHOWN).toList();
    }
    const positionItems = userPositions.map(this.renderPosition).toJS();

    return (
      <View style={styles.sectionContainer}>
        <Text style={styles.sectionHeader}>Positions</Text>
        <TouchableOpacity onPress={addPosition} style={styles.addTraitButton}>
          <MaterialIcons name="add-circle" size={32} color={Colors.HIVE_ACCENT} />
        </TouchableOpacity>
        { positionItems }
        <View style={styles.traitBottomActionContainer}>
          { bottomAction }
        </View>
      </View>
    );
  }

  private renderSimpleTrait(trait: UserSimpleTrait) {
    const onRemoveAccept = async () => {
      try {
        await this.props.removeSimpleTrait(trait.id);
      } catch (e) {
        await this.props.errorToast(e.errorMsg);
      }
    };

    const onRemovePress = () => {
      Alert.alert(
        'Remove Trait',
        `Are you sure you want to remove "${trait.simpleTraitName}" as a trait?`,
        [
          {text: 'Cancel', onPress: () => null, style: 'cancel'},
          {text: 'Remove', onPress: onRemoveAccept, style: 'destructive'},
        ],
      );
    }

    return (
      <View key={ trait.id } style={styles.simpleTraitContainer}>
        <Text style={styles.simpleTraitText}>{ trait.simpleTraitName }</Text>
        <TouchableOpacity style={styles.traitDelete} onPress={onRemovePress}>
          <MaterialIcons color={Colors.WHITE} name="close" size={18} />
        </TouchableOpacity>
      </View>
    );
  }

  private renderSimpleTraits() {
    let { userSimpleTraits } = this.props.profile;
    const { showAllSimpleTraits } = this.state;

    const addSimpleTrait = () => this.props.navigation.navigate('AddSimpleTrait');

    let bottomAction: ReactNode = null;
    if (userSimpleTraits.isEmpty()) {
      bottomAction = [
        <Text key={'text'} style={styles.noTraitText}>You don't have any traits</Text>,
        <Button
          key={'button'}
          buttonStyle={styles.noTraitButton}
          title="Add trait"
          onPress={addSimpleTrait}
          color={Colors.HIVE_PRIMARY}
        />,
      ];
    } else if (userSimpleTraits.size > MAX_NUMBER_SIMPLE_TRAITS_SHOWN) {
      bottomAction = ProfileView.renderShowLessMore(showAllSimpleTraits,
        () => this.setState({ showAllSimpleTraits: true }),
        () => this.setState({ showAllSimpleTraits: false }));
    }

    if (!showAllSimpleTraits) {
      userSimpleTraits = userSimpleTraits.take(MAX_NUMBER_SIMPLE_TRAITS_SHOWN).toList();
    }
    const traitItems = userSimpleTraits.map(this.renderSimpleTrait).toJS();

    return (
      <View style={styles.sectionContainer}>
        <Text style={styles.sectionHeader}>Traits</Text>
        <TouchableOpacity onPress={addSimpleTrait} style={styles.addTraitButton}>
          <MaterialIcons name="add-circle" size={32} color={Colors.HIVE_PRIMARY} />
        </TouchableOpacity>
        <View style={styles.simpleTraitOuterContainer}>
          { traitItems }
        </View>
        <View style={styles.traitBottomActionContainer}>
          { bottomAction }
        </View>
      </View>
    );
  }

  private renderBody() {
    const { navigate } = this.props.navigation;

    const {
      userPositions,
      userSimpleTraits,
    } = this.props.profile;

    const {
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

    const capitalize = (s: string) => s.charAt(0).toUpperCase() + s.slice(1);

    let userId;
    if (this.props.profile) {
      userId = this.props.profile.userId.toString();
    }
    const headerText = this.props.profile ?
      this.props.profile.firstName + ' ' + this.props.profile.lastName : 'Profile';

    const genderStr = capitalize(genderIdToString(gender));
    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    const timeDiff = new Date().valueOf() - new Date(birthdate).valueOf();
    const age = Math.floor(timeDiff / (1000 * 60 * 60 * 24 * 365));

    const hometownStr = hometown === null || hometown === '' ? 'Some place on Earth' : hometown;

    return (
      <View>
        <ScrollView contentContainerStyle={styles.container}>
          <View style={styles.contentContainer} >
            <ProfileAvatar userId={userId} xlarge containerStyle={styles.profilePicture} />
            <Header>{headerText}</Header>
            <Text style={styles.subHeaderText}>{age}{genderStr[0]} - {hometownStr}</Text>
            <TouchableOpacity style={styles.listItem} onPress={() => navigate('QrCode', { secret })}>
              <Text style={styles.value}>Show QR Code</Text>
            </TouchableOpacity>
            {this.renderProfile(bio)}
            {this.renderCohortInfo()}
            {this.renderContactInfo(email, fbId, fbLink, phoneNumber)}
            {this.renderPositions()}
            {this.renderSimpleTraits()}
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

const BUTTON_WIDTH = SCREEN_WIDTH - 80;
const styles = StyleSheet.create({
  container: {
    paddingTop: 10,
    paddingBottom: 65,
    backgroundColor: 'white',
    minHeight: '100%'
  },
  contentContainer: {
    alignItems: 'center',
    flex: 1,
    flexDirection: 'column',
    padding: 20,
  },
  description: {
    fontSize: 18,
    color: Colors.HIVE_SUBDUED
  },
  listItem: {
    flex: 1,
    flexDirection: 'row',
    marginTop: 5,
  },
  changePassButton: {
    width: BUTTON_WIDTH,
    marginTop: 10,
  },
  logoutButtonText: {
    color: 'white',
  },
  logoutButton: {
    width: BUTTON_WIDTH,
    marginTop: 10,
    backgroundColor: "gray",
    borderWidth: 0,
  },
  profileTitle: {
    fontSize: 18,
    marginTop: 10,
    alignSelf: 'flex-end'
  },
  profilePicture: {
    margin: 20
  },
  sectionHeader: {
    fontWeight: 'bold',
    fontSize: 24,
    alignSelf: 'flex-start',
  },
  sectionContainer: {
    width: "100%",
    flex: 1,
    flexDirection: 'column',
    marginTop: 20
  },
  subHeaderText: {
    fontSize: 18
  },
  label: {
    fontSize: 16,
  },
  value: {
    fontSize: 16,
    color: Colors.HIVE_ACCENT
  },
  cohortText: {
    fontSize: 16,
    marginTop: 5,
  },
  positionContainer: {
    backgroundColor: Colors.HIVE_ACCENT,
    marginTop: 5,
    padding: 5,
    paddingRight: 25,
    borderRadius: 5,
  },
  positionText: {
    color: Colors.WHITE,
    fontSize: 14,
  },
  positionBold: {
    fontWeight: '700',
  },
  traitDelete: {
    position: 'absolute',
    padding: 5,
    top: 0,
    right: 0,
  },
  simpleTraitOuterContainer: {
    flexWrap: 'wrap',
    flexDirection: 'row',
    justifyContent: 'center',
    paddingVertical: 2.5,
  },
  simpleTraitContainer: {
    backgroundColor: Colors.HIVE_PRIMARY,
    marginVertical: 2.5,
    marginHorizontal: 2.5,
    padding: 5,
    paddingRight: 25,
    paddingLeft: 10,
    borderRadius: 20,
  },
  simpleTraitText: {
    color: Colors.WHITE,
    fontWeight: '700',
    fontSize: 14,
  },
  traitBottomActionContainer: {
    flexDirection: 'column',
    alignItems: 'center',
  },
  addTraitButton: {
    position: 'absolute',
    top: 0,
    right: 0,
  },
  noTraitText: {
    marginTop: 10,
  },
  noTraitButton: {
    width: 200,
    marginTop: 10,
  },
});
