import React, { Component, ReactNode } from 'react';
import { connect, ActionCreator } from 'react-redux';
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

import { RootState } from '../redux';
import {
  State as BootstrapState,
  fetchBootstrap,
  removeRtmMatches,
} from '../redux/bootstrap/reducer';
import {
  State as CredentialOptionsState,
  fetchCredentialOptions,
} from '../redux/credential-options/reducer';
import { ActionTypes as BootstrapActionTypes } from '../redux/bootstrap/actions';
import { ActionTypes as CredentialOptionsActionTypes } from '../redux/credential-options/actions';
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
  programById,
  sequenceById,
} from '../models/cohort';
import Colors from '../services/colors';
import { AnalyticsHelper, AnalyticsActions, logAnalyticsThenExecute } from '../services/analytics';
import TopHeader, { headerStyle } from './TopHeader';
import AllFilterableModals from './AllFilterableModals';

interface DispatchActions {
  fetchBootstrap: ActionCreator<ThunkAction<Promise<BootstrapActionTypes>, BootstrapState, void>>;
  fetchCredentialOptions: ActionCreator<
    ThunkAction<Promise<CredentialOptionsActionTypes>, CredentialOptionsState, void>>;
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
    AnalyticsHelper.getInstance().recordPage(this.HOME_VIEW_IDENTIFIER);
    this.load();
  }

  async componentWillReceiveProps(nextProps: Props) {
    if (nextProps.bootstrap && nextProps.bootstrap.state === 'account_created') {
      // Account not yet setup, so take to onboarding page
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        actions: [NavigationActions.navigate({ routeName: 'Onboarding' })]
      }));
    }
  }

  private async load() {
    await Promise.all([
      this.props.fetchBootstrap(),
      this.props.fetchCredentialOptions(),
    ]);
  }

  private async onRefresh() {
    this.setState({refreshing: true});
    await this.load();
    this.setState({refreshing: false});
  }

  private renderDescription(relationship: Relationship) {
    const {
      userType,
      description,
      cohort,
    } = relationship;
    const cohortText = cohort === null ? 'some unknown program' :
      programById(cohort.programId) + ' ' + cohort.gradYear;

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
        return <Text style={styles.description}>
          They're looking for: <Text style={styles.bold}>{ description }</Text>
        </Text>;
      case USER_TYPE_ANSWERER:
        return <Text style={styles.description}>
          You're looking for: <Text style={styles.bold}>{ description }</Text>
        </Text>;
    }
  }

  private renderContactButton(relationship: Relationship) {
    var icon: string;
    var onPress: () => void;
    const {
      userId,
      email,
      fbId,
      fbLink,
      phoneNumber,
    } = relationship;

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
        getHumanReadableUserType(relationship.userType),
        1,
        () => {
          this.props.navigation.navigate('MatchProfile', { userId });}
      );
      onPress = logAnalyticsThenExecute.bind(
        this,
        "ContactProfile_" + profileType,
        AnalyticsActions.CLICK,
        getHumanReadableUserType(relationship.userType),
        1,
        onPress,
      );

    // TODO: Move into styles
    return (
      <View style={{ marginTop: 10, flexDirection: 'row', justifyContent: 'space-between' }}>
        <TouchableOpacity onPress={ viewProfile }>
          <Text style={{ color: Colors.HIVE_PRIMARY, fontSize: 18, paddingTop: 2 }}>View Profile</Text>
        </TouchableOpacity>
        <Button buttonStyle={{ width: 150 }} icon={icon} title="Contact" onPress={onPress} />
      </View>
    );
  }

  private renderMatch(relationship: Relationship) {
    const {
      userId,
      userType,
      firstName,
      lastName,
      matchingState,
    } = relationship;
    const isVerified = matchingState === MatchingState.Verified;
    const isUnverified = matchingState === MatchingState.Unverified;
    const description = this.renderDescription(relationship);

    const onCloseAccept = async () => {
      await this.props.removeRtmMatches(userId);
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


    // TODO: Handle errors for links
    return (
      <Card key={userId}>
        <View style={{ flexDirection: 'row' }}>
          <View style={styles.cardProfilePicture}>
            <ProfileAvatar userId={userId.toString()} large/>
          </View>
          <View style={{ flex: 2 }}>
            { isVerified && <Text style={styles.verified}>Verified Match</Text> }
            { isUnverified && <Text style={styles.unverified}>Unverified Match</Text> }
            <Text style={styles.name}>{firstName + ' ' + lastName}</Text>
            { description }
          </View>
        </View>
        { this.renderContactButton(relationship) }
        { closeButton }
      </Card>
    );
  }

  private renderMatches() {
    const { relationships } = this.props.bootstrap;

    const mentors = relationships.filter(rel => {
      return rel.userType === USER_TYPE_MENTOR;
    }).map(this.renderMatch).toList();
    const mentees = relationships.filter(rel => {
      return rel.userType === USER_TYPE_MENTEE;
    }).map(this.renderMatch).toList();
    const connections = relationships.filter(rel => {
      return rel.userType === USER_TYPE_ASKER || rel.userType === USER_TYPE_ANSWERER;
    }).map(this.renderMatch).toList();

    const elements: Array<ReactNode> = [];
    const pushRels = (title: string, rels: Immutable.List<ReactNode>) => {
      if (rels.size > 0) {
        elements.push(<Header>{ title + (rels.size > 1 ? 's' : '') }</Header>);
        elements.push(rels.toJS());
      }
    };

    pushRels('Your Mentor', mentors);
    pushRels('Your Mentee', mentees);

    if (connections.size > 0) {
      elements.push(<Header>Your Connections</Header>);
      elements.push(connections.toJS());
    }

    return <View>{ elements }</View>;
  }

  private renderHome() {
    // A little sketchy to be pasting this in all the cases, but haven't found an easy work around
    // yet.
    const allModals = (
      <AllFilterableModals
        onSelectSuccess={() => {
          this.props.navigation.navigate({ routeName: 'Requests' });
        }}
      />
    );
    const { state } = this.props.bootstrap;
    switch (state) {
      case 'account_created':
        // Should in reality not be shown, since we never show home page for account_created
        return (
          <View style={styles.centeredContainer}>
            <Text style={styles.headline}>Waiting for you to finish onboarding</Text>
            <ActionButton onPress={() => this.load()} title="Check again" />
            { allModals }
          </View>
        );
      case 'account_setup':
        return (
          <View style={styles.centeredContainer}>
            <Text style={styles.headline}>Waiting for your match</Text>
            <ActionButton onPress={() => this.load()} title="Check again" />
            { allModals }
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
              { matches }
            </ScrollView>
            { allModals }
          </View>
        );
      default:
        // Ensure exhaustiveness of select
        const _: never = state;
    }
  }

  render() {
    const {
      state,
      errorMsg,
      errorType,
    } = this.props.fetchState;
    // If `this.state.refreshing` is true, it means that we are reloading data using the pull
    // down, which means that we want to still display the ScrollView.
    return (
      <Loading
        state={this.state.refreshing ? 'success' : state}
        errorMsg={errorMsg}
        errorType={errorType}
        load={this.load}
        renderBody={this.renderHome}
        navigation={this.props.navigation}
      />
    );
  }
}

export default connect(({ bootstrap }: RootState) => bootstrap,
  { fetchBootstrap, fetchCredentialOptions, removeRtmMatches })(HomeView);

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
  cardProfilePicture: {
    flex: 1,
    marginRight: 10,
  },
  bold: {
    fontWeight: 'bold',
  },
  verified: {
    color: 'green',
    fontWeight: 'bold',
  },
  unverified: {
    color: 'red',
    fontWeight: 'bold',
  },
  deleteRtmMatch: {
    position: 'absolute',
    top: 10,
    right: 10,
  },
})
