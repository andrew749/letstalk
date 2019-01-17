import React, { Component, ReactNode } from 'react';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import {
  ActivityIndicator,
  AsyncStorage,
  Alert,
  Button as ReactNativeButton,
  Image,
  ImageBackground,
  Linking,
  Modal as ReactNativeModal,
  Picker,
  Platform,
  RefreshControl,
  RefreshControlProps,
  ScrollView,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
} from 'react-native';
import {
  NavigationScreenProp,
  NavigationScreenDetails,
  NavigationStackAction,
  NavigationActions,
  NavigationEventSubscription,
} from 'react-navigation';
import { MaterialIcons, FontAwesome } from '@expo/vector-icons';
import Immutable from 'immutable';
import Modal from 'react-native-modal';

import requestToMatchService from '../services/request-to-match-service';
import { BootstrapConnection } from '../models/bootstrap';
import { RootState } from '../redux';
import {
  State as BootstrapState,
  fetchBootstrap,
} from '../redux/bootstrap/reducer';
import { State as SurveyState, fetchSurvey } from '../redux/survey/reducer';
import { errorToast, infoToast } from '../redux/toast';
import { ActionTypes as BootstrapActionTypes } from '../redux/bootstrap/actions';
import { ActionTypes as SurveyActionTypes } from '../redux/survey/actions';
import { ActionButton, Button, Card, Header, ProfileAvatar } from '../components';
import { ClubDayInformationCard, ProfileFillCallToAction } from '../components/InformationCard';
import Loading from './Loading';
import {
  USER_STATE_ACCOUNT_CREATED,
  USER_STATE_ACCOUNT_EMAIL_VERIFIED,
  USER_STATE_ACCOUNT_HAS_BASIC_INFO,
  USER_STATE_ACCOUNT_SETUP,
  USER_STATE_ACCOUNT_MATCHED,
  MatchingState,
  Relationship,
  getHumanReadableUserType,
} from '../models/bootstrap';
import {
  USER_TYPE_MENTOR,
  USER_TYPE_MENTEE,
  USER_TYPE_ASKER,
  USER_TYPE_ANSWERER,
} from '../models/user';
import {
  IntentTypes,
} from '../models/connection';
import {
  programById,
  sequenceById,
} from '../models/cohort';
import Colors from '../services/colors';
import { AnalyticsHelper, AnalyticsActions, logAnalyticsThenExecute } from '../services/analytics';
import TutorialService from '../services/tutorial_service';
import TopHeader, { headerStyle, headerTitleStyle, headerTintColor  } from './TopHeader';
import AllFilterableModals from './AllFilterableModals';
import { GROUP_GENERIC } from '../services/survey';
import { FETCH_STATE_PREFETCH, FETCH_STATE_FETCHING, FETCH_STATE_SUCCESS } from '../redux/actions';
import Window from '../services/window';
import Color from '../services/colors';
import profileService from '../services/profile-service';
import auth from '../services/auth';

interface ContactModalProps {
  relationship: Relationship;
}

interface ContactModalState {
  isVisible: boolean;
}

class ContactModal extends Component<ContactModalProps, ContactModalState> {

  constructor(props: ContactModalProps) {
    super(props);
    this.state = { isVisible: false }
  }

  private getButton(
    title: string,
    profileType: string,
    onClick: () => void,
    icon?: string,
    iconComponent?: ReactNode,
  ): ReactNode {
    const executable = () => {
      onClick();
      this.setState({ isVisible: false });
    }
    const onPress = logAnalyticsThenExecute.bind(
      this,
      "ContactProfile_" + profileType,
      AnalyticsActions.CLICK,
      getHumanReadableUserType(this.props.relationship.userType),
      1,
      executable);

    return <Button
      key={ title }
      title={ title }
      icon={ icon || null }
      iconComponent={ iconComponent || null }
      buttonStyle={ contactModalStyles.button }
      onPress={ onPress }
    />;
  }

  private getButtons(): Array<ReactNode> {
    const buttonOpts = [];
    const {
      email,
      fbId,
      fbLink,
      phoneNumber,
    } = this.props.relationship;

    if (fbLink != null) {
      buttonOpts.push({
        title: 'Add as a friend on FB',
        icon: null,
        iconComponent: <FontAwesome
          style={{ position: 'absolute', left: 2, top: 1, margin: 4 }}
          color={Colors.HIVE_PRIMARY}
          name={'facebook-official'}
          size={24}
        />,
        profileType: 'Facebook',
        onPress: () => Linking.openURL(fbLink),
      });
    }
    if (phoneNumber !== null) {
      const smsLink = 'sms:' + phoneNumber;
      buttonOpts.push({
        title: 'Send a text',
        icon: 'textsms',
        iconComponent: null,
        profileType: 'Sms',
        onPress: () => Linking.openURL(smsLink),
      });
    }
    const emailLink = 'mailto:' + email;
    buttonOpts.push({
      title: 'Send an email',
      icon: 'email',
      iconComponent: null,
      profileType: 'Email',
      onPress: () => Linking.openURL(emailLink),
    });

    return buttonOpts.map(opt => this.getButton(
      opt.title,
      opt.profileType,
      opt.onPress,
      opt.icon,
      opt.iconComponent,
    ));
  }

  render() {
    const onPress = () => this.setState({ isVisible: true });
    const icon = 'textsms';
    const buttons = this.getButtons();
    return (
      <View>
        <Modal isVisible={ this.state.isVisible }>
          <View style={ contactModalStyles.container }>
            <Text style={contactModalStyles.header}>Contact Options</Text>
            { buttons }
            <Button
              title={ 'Close' }
              buttonStyle={ contactModalStyles.closeButton }
              textStyle={ contactModalStyles.closeButtonText }
              onPress={ () => this.setState({ isVisible: false }) }
            />
          </View>
        </Modal>
        <Button buttonStyle={{ width: 150 }} icon={ icon } title="Contact" onPress={ onPress } />
      </View>
    );
  }
}

const contactModalStyles = StyleSheet.create({
  container: {
    padding: 10,
    borderRadius: 10,
    backgroundColor: Colors.WHITE,
  },
  header: {
    fontWeight: 'bold',
    fontSize: 20,
  },
  button: {
    marginTop: 10,
  },
  closeButton: {
    marginTop: 20,
    borderColor: Colors.RED,
    backgroundColor: Colors.RED,
  },
  closeButtonText: {
    color: Colors.WHITE,
  },
});

interface DispatchActions {
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  fetchBootstrap: ActionCreator<ThunkAction<Promise<BootstrapActionTypes>, BootstrapState, void>>;
  fetchSurvey: ActionCreator<ThunkAction<Promise<SurveyActionTypes>, SurveyState, void>>;
}

interface Props extends BootstrapState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

interface State {
  refreshing: boolean;
  focused: boolean;
}

// Since this view gets rendered multiple times, we only want to add the device token once per app
// launch.
let addedExpoToken = false;

const PROBABILITY_SHOW_IOS_MODAL = 0.4;

function showModal(): boolean {
  return Math.random() < PROBABILITY_SHOW_IOS_MODAL;
}

class HomeView extends Component<Props, State> {
  HOME_VIEW_IDENTIFIER = "HomeView";

  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: <TopHeader navigation={navigation} />,
    headerStyle,
    headerTitleStyle,
    headerTintColor
  })

  private willFocusHandler: NavigationEventSubscription;
  private willBlurHandler: NavigationEventSubscription;

  private async maybeAddExpoToken() {

    if (!addedExpoToken) {
      addedExpoToken = true;
      try {
        const token = await auth.registerForPushNotificationsAsync(showModal());
        if (token !== null) {
          await profileService.addExpoDeviceToken(token);
        }
      } catch(e){
        console.log("Failed to register for notification " + e);
      }
    }
  }

  constructor(props: Props) {
    super(props);

    this.state = { refreshing: false, focused: false };

    this.loadBootstrap = this.loadBootstrap.bind(this);
    this.renderHome = this.renderHome.bind(this);
    this.renderMatch = this.renderMatch.bind(this);
    this.onRefresh = this.onRefresh.bind(this);
  }

  async componentDidMount() {
    // Hack since for some reason state is of type void, when it clearly has things inside
    if ((this.props.navigation.state as any).routeName === 'Home') {
      this.setState({ focused: true })
    } else {
      this.setState({ focused: false })
    }

    this.willBlurHandler = this.props.navigation.addListener('willBlur', (route) => {
      this.setState({ focused: false })
    });
    this.willFocusHandler = this.props.navigation.addListener('willFocus', (route) => {
      this.setState({ focused: true });
      AnalyticsHelper.getInstance().recordPage(this.HOME_VIEW_IDENTIFIER);
    });

    await Promise.all([
      this.props.fetchBootstrap(),
      this.maybeAddExpoToken(),
    ]);
    await TutorialService.launchTutorial(this.props.navigation);
  }

  componentWillUnmount() {
    this.willFocusHandler.remove();
    this.willBlurHandler.remove();
  }

  async componentDidUpdate() {
    await this.maybeNavigateRequired();
  }

  // Depending on the user's state, we may need to navigate to another view to get more info
  private async maybeNavigateRequired() {
    // Need to check if we are currently focused because of the following edge case:
    // User is not done onboarding but gets a push notification.
    // When the app is opened, the notification is opened but then overriden by one of the below
    // navigation actions.
    // The user never actually ends up seeing the notification, resulting in a bad user experience.
    if (!!this.props.bootstrap && this.state.focused) {
      switch (this.props.bootstrap.state) {
        case USER_STATE_ACCOUNT_CREATED:
          await this.props.navigation.dispatch(NavigationActions.reset({
            index: 0,
            actions: [NavigationActions.navigate({ routeName: 'VerifyEmail' })]
          }));
          break;
        case USER_STATE_ACCOUNT_EMAIL_VERIFIED:
          await this.props.navigation.dispatch(NavigationActions.reset({
            index: 0,
            actions: [NavigationActions.navigate({ routeName: 'Onboarding' })]
          }));
          break;
        case USER_STATE_ACCOUNT_HAS_BASIC_INFO:
          await this.props.fetchSurvey(GROUP_GENERIC);
          await this.props.navigation.dispatch(NavigationActions.reset({
            index: 0,
            actions: [NavigationActions.navigate({ routeName: 'SurveyView' })]
          }));
          break;
        case USER_STATE_ACCOUNT_SETUP:
          // fallthrough
        case USER_STATE_ACCOUNT_MATCHED:
          // fallthrough
      }
    }
  }

  private async loadBootstrap() {
    await this.props.fetchBootstrap();
  }

  private async onRefresh() {
    this.setState({refreshing: true});
    await this.loadBootstrap();
    this.setState({refreshing: false});
  }

  private renderDescription(connection: BootstrapConnection) {
    const {
      userType,
      cohort,
    } = connection.userProfile;
    const {
      intentType,
      searchedTrait,
    } = connection.request;
    const cohortText = cohort === null ? 'some unknown program' :
      programById(cohort.programId) + ' ' + cohort.gradYear;

    let description: ReactNode = null;
    switch (intentType) {
      case IntentTypes.REC_COHORT:
        description = <Text>You are both in the same cohort</Text>;
        break;
      case IntentTypes.SCAN_CODE:
        description = <Text>You connected by QR code</Text>;
        break;
      case IntentTypes.REC_GENERAL:
        switch (userType) {
          case USER_TYPE_ASKER:
            description = (
              <Text>They requested to connect with you after you were recommended to them</Text>
            );
            break;
          case USER_TYPE_ANSWERER:
            description = (
              <Text>You requested to connect with them after they were recommended to you</Text>
            );
            break;
        }
        break;
      case IntentTypes.SEARCH:
        switch (userType) {
          case USER_TYPE_ASKER:
            description = (
              <Text>{'They connected with you for: '}
                <Text style={styles.bold}>{ searchedTrait }</Text>
              </Text>
            );
            break;
          case USER_TYPE_ANSWERER:
            description = (
              <Text>{'You connected with them for: '}
                <Text style={styles.bold}>{ searchedTrait }</Text>
              </Text>
            );
            break;
        }
        break;
      default:
        const _: never = intentType;
    }

    switch (userType) {
      case USER_TYPE_MENTOR:
        return <Text style={styles.description}>
          Your mentor in <Text style={styles.bold}>{ cohortText }</Text>
        </Text>;
      case USER_TYPE_MENTEE:
        return <Text style={styles.description}>
          Your mentee in <Text style={styles.bold}>{ cohortText }</Text>
        </Text>;
      case USER_TYPE_ASKER:
      case USER_TYPE_ANSWERER:
        return <Text style={styles.description}>{ description }</Text>;
    }
    return null;
  }

  private renderInformationCards() {
    return (
      <View>
        <ClubDayInformationCard />
        <ProfileFillCallToAction />
      </View>
    );
  }

  private renderContactButton(connection: BootstrapConnection) {
    const { userProfile } = connection;
    var icon: string;
    var onPress: () => void;
    const { userId } = userProfile;

    // record user clicking to view mentor/mentee profile
    const viewProfile =
      logAnalyticsThenExecute.bind(
        this,
        "Profile",
        AnalyticsActions.CLICK,
        getHumanReadableUserType(userProfile.userType),
        1,
        () => {
          this.props.navigation.navigate('MatchProfile', { userId });}
      );

    // TODO: Move into styles
    return (
      <View style={{ marginTop: 10, flexDirection: 'row', justifyContent: 'space-between' }}>
        <TouchableOpacity onPress={ viewProfile }>
          <Text style={{ color: Colors.HIVE_PRIMARY, fontSize: 18, paddingTop: 2 }}>
            View Profile
          </Text>
        </TouchableOpacity>
        <ContactModal relationship={userProfile} />
      </View>
    );
  }

  private renderMatch(
    connection: BootstrapConnection,
    idx: number,
    arr: Immutable.List<BootstrapConnection>,
  ) {
    const {
      userId,
      userType,
      firstName,
      lastName,
    } = connection.userProfile;
    const description = this.renderDescription(connection);

    const onCloseAccept = async () => {
      try {
        await requestToMatchService.removeConnection(userId);
        this.props.infoToast('Removed connection');
        this.props.fetchBootstrap();
      } catch (e) {
        await this.props.errorToast(e.errorMsg);
      }
    }

    const onClosePress = () => {
      Alert.alert(
        'Unmatch',
        'Are you sure you want to unmatch? This will permanently remove ' + firstName +
        ' from your list of connections',
        [
          {text: 'Cancel', onPress: () => null, style: 'cancel'},
          {text: 'Unmatch', onPress: onCloseAccept, style: 'destructive'},
        ],
      );
    }

    const closeButton = (userType === USER_TYPE_ASKER || userType === USER_TYPE_ANSWERER) ? (
      <TouchableOpacity style={styles.deleteRtmMatch} onPress={onClosePress}>
        <MaterialIcons name="close" size={18} />
      </TouchableOpacity>
    ) : null;

    const extraStyle = arr.size === idx + 1 ? null : {marginBottom: 10};

    // TODO: Handle errors for links
    return (
      <Card key={userId} style={[styles.connectionCard, extraStyle]}>
        <View style={{ flexDirection: 'row' }}>
          <View style={styles.cardProfilePicture}>
            <ProfileAvatar userId={userId.toString()} large/>
          </View>
          <View style={{ flex: 2 }}>
            <Text style={styles.name}>{firstName + ' ' + lastName}</Text>
            { description }
          </View>
        </View>
        { this.renderContactButton(connection) }
        { closeButton }
      </Card>
    );
  }

  private renderPeerMatches() {
    const { connections } = this.props.bootstrap;
    const peers = connections.peers.map(this.renderMatch).toList();

    const elements: Array<ReactNode> = [];

    if (peers.size > 0) {
      elements.push(<Header key={'Your Connections'}>Your Connections</Header>);
      elements.push(peers.toJS());
    }

    return <View>{ elements }</View>;
  }

  private renderMatches() {
    const { connections } = this.props.bootstrap;

    const mentors = connections.mentors.map(this.renderMatch).toList();
    const mentees = connections.mentees.map(this.renderMatch).toList();
    const peers = connections.peers.map(this.renderMatch).toList();

    const elements: Array<ReactNode> = [];
    const pushRels = (title: string, rels: Immutable.List<ReactNode>) => {
      if (rels.size > 0) {
        elements.push(<Header key={title}>{ title + (rels.size > 1 ? 's' : '') }</Header>);
        elements.push(rels.toJS());
      }
    };

    pushRels('Your Mentor', mentors);
    pushRels('Your Mentee', mentees);

    if (peers.size > 0) {
      elements.push(<Header key={'Your Connections'}>Your Connections</Header>);
      elements.push(peers.toJS());
    }

    return <View>{ elements }</View>;
  }

  private renderRequestsButton() {
    const {
      incomingRequests,
      outgoingRequests,
    } = this.props.bootstrap.connections;
    if (incomingRequests.isEmpty() && outgoingRequests.isEmpty()) return null;

    let description = `You have ${incomingRequests.size} incoming and ${outgoingRequests.size} ` +
      `outgoing requests to connect.`;

    return (
      <View>
        <Text style={styles.requestsButtonText}>{ description }</Text>
        <Button
          buttonStyle={styles.feedbackButton}
          textStyle={styles.feedbackButtonText}
          title="See Requests"
          onPress={() => this.props.navigation.navigate('Requests')}
        />
      </View>
    );
  }

  private renderHome() {
    // A little sketchy to be pasting this in all the cases, but haven't found an easy work around
    // yet.
    const requestsButton = this.renderRequestsButton();
    const informationCards = this.renderInformationCards();

    const { state } = this.props.bootstrap;
    switch (state) {
      case USER_STATE_ACCOUNT_CREATED:
        // fallthrough
      case USER_STATE_ACCOUNT_EMAIL_VERIFIED:
        // fallthrough
      case USER_STATE_ACCOUNT_HAS_BASIC_INFO:
        // Should in reality not be shown, since we never show home page until state account_setup.
        return (
          <View style={styles.centeredContainer}>
            <Text style={styles.headline}>Waiting for you to finish onboarding</Text>
            <ActionButton onPress={() => this.loadBootstrap()} title="Check again" />
          </View>
        );
      case USER_STATE_ACCOUNT_SETUP:
        const peerMatches = this.renderPeerMatches();
        return (

          <View style={styles.container}>
            <ScrollView
              refreshControl={
                <RefreshControl
                  refreshing={this.state.refreshing}
                  onRefresh={this.onRefresh}
                /> as React.ReactElement<RefreshControlProps>
              }
            >
              <View style={styles.scrollContainer}>
                <View style={styles.centeredContainer}>
                  <Text style={styles.headline}>Waiting for your mentorship match</Text>
                  <ActionButton onPress={() => this.loadBootstrap()} title="Check again" />
                </View>
                { requestsButton }
                { informationCards }
                { peerMatches }
              </View>
            </ScrollView>
          </View>
        );
      case USER_STATE_ACCOUNT_MATCHED:
        const matches = this.renderMatches();
        // Watch out! Typescript hack below.
        return (
          <View style={styles.container}>
            <ScrollView
              refreshControl={
                <RefreshControl
                  refreshing={this.state.refreshing}
                  onRefresh={this.onRefresh}
                /> as React.ReactElement<RefreshControlProps>
              }
            >
              <View style={styles.scrollContainer}>
                { requestsButton }
                { informationCards }
                { matches }
              </View>
            </ScrollView>
          </View>
        );
      default:
        // Ensure exhaustiveness of select
        const _: never = state;
    }
  }

  render() {
    const allModals = (
      <AllFilterableModals
        onSelectSuccess={() => {
          this.props.navigation.navigate({ routeName: 'Explore' });
        }}
      />
    );
    const {
      state,
      errorMsg,
      errorType,
    } = this.props.fetchState;
    // If `this.state.refreshing` is true, it means that we are reloading data using the pull
    // down, which means that we want to still display the ScrollView.
    const refreshingState = this.state.refreshing ? FETCH_STATE_SUCCESS : state;
    if (refreshingState === FETCH_STATE_FETCHING || refreshingState === FETCH_STATE_PREFETCH) {
      return (
        <ReactNativeModal visible={true} animationType={'none'}>
          <View style={styles.splashScreenContainer}>
            <ImageBackground
              style={styles.hiveLogo}
              source={require('../img/logo_android.png')}
            >
              <ActivityIndicator
                color={Color.HIVE_PRIMARY}
                size="large"
                style={styles.spinner}
              />
            </ImageBackground>
            <Image
              style={styles.hiveName}
              source={require('../img/name_white.png')}
            />
          </View>
        </ReactNativeModal>
      );
    } else {
      return (
        <View style={{flex: 1}}>
          <Loading
            state={refreshingState}
            errorMsg={errorMsg}
            errorType={errorType}
            load={this.loadBootstrap}
            renderBody={this.renderHome}
            navigation={this.props.navigation}
          />
          { allModals }
        </View>
      );
    }
  }
}

export default connect(({ bootstrap }: RootState) => bootstrap,
  { errorToast, infoToast, fetchBootstrap, fetchSurvey })(HomeView);

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  centeredContainer: {
    flex: 1,
    justifyContent: 'center'
  },
  headline: {
    fontWeight: 'bold',
    fontSize: 18,
    textAlign: 'center',
    marginBottom: 10,
  },
  name: {
    fontWeight: 'bold',
    fontSize: 20,
  },
  description: {
    fontSize: 14,
    color: 'gray',
    marginBottom: 5,
  },
  emailContainer: {
    flex: 1,
    flexDirection: 'row',
  },
  email: {
    paddingTop: 2,
    marginLeft: 5,
    fontSize: 16,
  },
  feedbackButton: {
    backgroundColor: Colors.WHITE,
    alignSelf: 'center',
    width: 200,
    marginVertical: 15,
  },
  feedbackButtonText: {
    color: Colors.HIVE_PRIMARY,
    fontSize: 16
  },
  feedbackText: {
    fontSize: 14,
    marginHorizontal: 15
  },
  requestsButtonText: {
    marginTop: 10,
    fontSize: 14,
  },
  cardProfilePicture: {
    flex: 1,
    marginRight: 10,
  },
  bold: {
    fontWeight: 'bold',
  },
  verified: {
    color: Colors.GREEN,
    fontWeight: 'bold',
  },
  unverified: {
    color: Colors.RED,
    fontWeight: 'bold',
  },
  deleteRtmMatch: {
    position: 'absolute',
    top: 10,
    right: 10,
  },
  scrollContainer: {
    padding: 10,
  },
  connectionCard: {
    marginHorizontal: 0,
    marginVertical: 0,
  },
  splashScreenContainer: {
    width: Window.WIDTH,
    height: Window.HEIGHT,
    backgroundColor: Color.HIVE_PRIMARY,
    justifyContent: 'center',
    alignItems: 'center',
  },
  hiveLogo: {
    width: 200,
    height: 200,
  },
  hiveName: {
    marginTop: 20,
    width: 101,
    height: 57,
  },
  spinner: {
    position: 'absolute',
    // Bit of a hack that makes it easier to position the spinner
    width: 0,
    height: 0,
    left: Platform.OS === 'ios' ? 107 : 105,
    top: Platform.OS === 'ios' ? 123 : 121,
  },
})
