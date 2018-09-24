import React, { Component, ReactNode } from 'react';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import {
  ActivityIndicator,
  Alert,
  Button as ReactNativeButton,
  Linking,
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
  NavigationActions
} from 'react-navigation';
import { MaterialIcons } from '@expo/vector-icons';
import Immutable from 'immutable';

import requestToMatchService from '../services/request-to-match-service';
import { BootstrapConnection } from '../models/bootstrap';
import { RootState } from '../redux';
import {
  State as BootstrapState,
  fetchBootstrap,
  removeRtmMatches,
} from '../redux/bootstrap/reducer';
import { errorToast, infoToast } from '../redux/toast';
import { ActionTypes as BootstrapActionTypes } from '../redux/bootstrap/actions';
import { ActionButton, Button, Card, Header, ProfileAvatar } from '../components';
import Loading from './Loading';
import {MatchingState, Relationship, getHumanReadableUserType} from '../models/bootstrap';
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

interface DispatchActions {
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  fetchBootstrap: ActionCreator<ThunkAction<Promise<BootstrapActionTypes>, BootstrapState, void>>;
  removeRtmMatches: ActionCreator<ThunkAction<Promise<BootstrapActionTypes>, BootstrapState, void>>;
}

interface Props extends BootstrapState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

interface State {
  refreshing: boolean;
}

class HomeView extends Component<Props, State> {
  HOME_VIEW_IDENTIFIER = "HomeView";

  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: <TopHeader navigation={navigation} />,
    headerStyle, 
    headerTitleStyle, 
    headerTintColor 
  })

  constructor(props: Props) {
    super(props);

    this.state = { refreshing: false };

    this.load = this.load.bind(this);
    this.renderHome = this.renderHome.bind(this);
    this.renderMatch = this.renderMatch.bind(this);
    this.onRefresh = this.onRefresh.bind(this);
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.HOME_VIEW_IDENTIFIER);
    });

    this.load();
    await TutorialService.launchTutorial(this.props.navigation);
  }

  async componentWillReceiveProps(nextProps: Props) {
    if (nextProps.bootstrap && nextProps.bootstrap.state === 'account_created') {
      // Email not yet verified, so take to email verification page
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        actions: [NavigationActions.navigate({ routeName: 'VerifyEmail' })]
      }));
    } else if (nextProps.bootstrap && nextProps.bootstrap.state === 'account_email_verified') {
      // Account not yet setup, so take to onboarding page
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        actions: [NavigationActions.navigate({ routeName: 'Onboarding' })]
      }));
    }
  }

  private async load() {
    await this.props.fetchBootstrap();
  }

  private async onRefresh() {
    this.setState({refreshing: true});
    await this.load();
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
              <Text>{'They requested to connect with you for: '}
                <Text style={styles.bold}>{ searchedTrait }</Text>
              </Text>
            );
            break;
          case USER_TYPE_ANSWERER:
            description = (
              <Text>{'You requested to connect with them for: '}
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

  private renderContactButton(connection: BootstrapConnection) {
    const { userProfile } = connection;
    var icon: string;
    var onPress: () => void;
    const {
      userId,
      email,
      fbId,
      fbLink,
      phoneNumber,
    } = userProfile;

    let profileType;
    if (fbLink != null) {
      icon = 'face';
      profileType = "Facebook";
      onPress = () => Linking.openURL(fbLink);
    } else if (phoneNumber !== null) {
      const smsLink = 'sms:' + phoneNumber;
      profileType = "Phone";
      icon = 'textsms';
      onPress = () => Linking.openURL(smsLink);
    } else {
      const emailLink = 'mailto:' + email;
      profileType = "Email";
      icon = 'email';
      onPress = () => Linking.openURL(emailLink);
    }

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
      onPress = logAnalyticsThenExecute.bind(
        this,
        "ContactProfile_" + profileType,
        AnalyticsActions.CLICK,
        getHumanReadableUserType(userProfile.userType),
        1,
        onPress,
      );

    // TODO: Move into styles
    return (
      <View style={{ marginTop: 10, flexDirection: 'row', justifyContent: 'space-between' }}>
        <TouchableOpacity onPress={ viewProfile }>
          <Text style={{ color: Colors.HIVE_PRIMARY, fontSize: 18, paddingTop: 2 }}>
            View Profile
          </Text>
        </TouchableOpacity>
        <Button buttonStyle={{ width: 150 }} icon={icon} title="Contact" onPress={onPress} />
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
          title="View Connection Requests"
          onPress={() => this.props.navigation.navigate('Requests')}
        />
      </View>
    );
  }

  private renderHome() {
    // A little sketchy to be pasting this in all the cases, but haven't found an easy work around
    // yet.
    const feedbackPrompt = (
      <View>
        <Text style={styles.feedbackText}>
          Thank you for joining Hive! If you notice any bugs or
          have any suggestions, please let us know!
        </Text>
        <Button
          buttonStyle={styles.feedbackButton}
          textStyle={styles.feedbackButtonText}
          title="Submit Feedback"
          onPress={() => Linking.openURL('https://goo.gl/forms/dkZf8AcgPPCNW7xe2')}
        />
      </View>
    );
    const requestsButton = this.renderRequestsButton();

    const { state } = this.props.bootstrap;
    switch (state) {
      case 'account_created':
        // fallthrough
      case 'account_email_verified':
        // Should in reality not be shown, since we never show home page until state account_setup.
        return (
          <View style={styles.centeredContainer}>
            <Text style={styles.headline}>Waiting for you to finish onboarding</Text>
            <ActionButton onPress={() => this.load()} title="Check again" />
          </View>
        );
      case 'account_setup':
        return (

          <View style={styles.container}>
            <View style={styles.centeredContainer}>
              <Text style={styles.headline}>Waiting for your match</Text>
              <ActionButton onPress={() => this.load()} title="Check again" />
            </View>
            { feedbackPrompt }
          </View>
        );
      case 'account_matched':
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
                { feedbackPrompt }
                { requestsButton }
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
    return (
      <View style={{flex: 1}}>
        <Loading
          state={this.state.refreshing ? 'success' : state}
          errorMsg={errorMsg}
          errorType={errorType}
          load={this.load}
          renderBody={this.renderHome}
          navigation={this.props.navigation}
        />
        { allModals }
      </View>
    );
  }
}

export default connect(({ bootstrap }: RootState) => bootstrap,
  { errorToast, infoToast, fetchBootstrap, removeRtmMatches })(HomeView);

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
  }
})
