import React, { Component, ReactNode } from 'react';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import {
  ActivityIndicator,
  Button as ReactNativeButton,
  Linking,
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
import { State as BootstrapState, fetchBootstrap } from '../redux/bootstrap/reducer';
import { ActionTypes } from '../redux/bootstrap/actions';
import { ActionButton, Card, Header, Loading, ProfileAvatar } from '../components';
import { Relationship } from '../models/bootstrap';
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
import { AnalyticsHelper } from '../services/analytics';

interface DispatchActions {
  fetchBootstrap: ActionCreator<ThunkAction<Promise<ActionTypes>, BootstrapState, void>>;
}

interface Props extends BootstrapState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class HomeView extends Component<Props> {
  HOME_VIEW_IDENTIFIER = "HomeView";

  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: 'Home',
    headerStyle: {
      backgroundColor: Colors.HIVE_MAIN_BG,
    },
  })

  constructor(props: Props) {
    super(props);

    this.load = this.load.bind(this);
    this.renderHome = this.renderHome.bind(this);
    this.renderMatch = this.renderMatch.bind(this);
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
    await this.props.fetchBootstrap();
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

  private renderMatch(relationship: Relationship) {
    const {
      userId,
      firstName,
      lastName,
      email,
      fbId,
      phoneNumber,
    } = relationship;

    const fbLink = 'fb://profile/' + fbId;
    const emailLink = 'mailto:' + email;
    const smsLink = 'sms:' + phoneNumber;
    const fb = fbId === null ? null : (
      <TouchableOpacity style={styles.emailContainer} onPress={() => Linking.openURL(fbLink)}>
        <MaterialIcons name="face" size={24} />
        <Text style={styles.email}>Facebook profile</Text>
      </TouchableOpacity>
    );
    const sms = phoneNumber === null ? null : (
      <TouchableOpacity style={styles.emailContainer} onPress={() => Linking.openURL(smsLink)}>
        <MaterialIcons name="textsms" size={24} />
        <Text style={styles.email}>Talk on SMS</Text>
      </TouchableOpacity>
    );
    const description = this.renderDescription(relationship);
    // TODO: Handle errors for links
    return (
      <Card style={styles.card} key={userId}>
        <View style={styles.cardProfilePicture}>
          <ProfileAvatar userId={userId.toString()}/>
        </View>
        <View style={{ flex: 2 }}>
          <Text style={styles.name}>{firstName + ' ' + lastName}</Text>
          { description }
          <TouchableOpacity style={styles.emailContainer}
            onPress={() => Linking.openURL(emailLink)}
          >
              <MaterialIcons name="email" size={24} />
              <Text style={styles.email}>{email}</Text>
          </TouchableOpacity>
          {fb}
          {sms}
      </View>
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
    pushRels('Your Mentor', mentees);

    if (connections.size > 0) {
      elements.push(<Header>Your Connections</Header>);
      elements.push(connections.toJS());
    }

    return <View>{ elements }</View>;
  }

  private renderHome() {
    const { state } = this.props.bootstrap;
    switch (state) {
      case 'account_created':
        // Should in reality not be shown, since we never show home page for account_created
        return (
          <View style={styles.centeredContainer}>
            <Text style={styles.headline}>Waiting for you to finish onboarding</Text>
            <ActionButton onPress={() => this.load()} title="Check again" />
          </View>
        );
      case 'account_setup':
        return (
          <View style={styles.centeredContainer}>
            <Text style={styles.headline}>Waiting for your match</Text>
            <ActionButton onPress={() => this.load()} title="Check again" />
          </View>
        );
      case 'account_matched':
        const matches = this.renderMatches();
        return (
          <ScrollView>
            { matches }
          </ScrollView>
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
    } = this.props.fetchState;
    return (
      <Loading
        state={state}
        errorMsg={errorMsg}
        load={this.load}
        renderBody={this.renderHome}
      />
    );
  }
}

export default connect(({ bootstrap }: RootState) => bootstrap, { fetchBootstrap })(HomeView);

const styles = StyleSheet.create({
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
    fontSize: 18,
  },
  description: {
    fontSize: 12,
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
  card: {
    flex: 1,
    flexDirection: 'row',
  },
  cardProfilePicture: {
    flex: 1,
    marginRight: 10,
  },
  bold: {
    fontWeight: 'bold',
  },
})
