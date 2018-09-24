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
import { State as MatchProfileState, fetchMatchProfile } from '../redux/match-profile/reducer';
import { ActionTypes as MatchProfileActionTypes } from '../redux/match-profile/actions';
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
import TopHeader, { headerStyle, headerTitleStyle, headerTintColor } from './TopHeader';
import AllFilterableModals from './AllFilterableModals';
import { ConnectionRequestWithName } from '../models/bootstrap';

interface DispatchActions {
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  fetchBootstrap: ActionCreator<ThunkAction<Promise<BootstrapActionTypes>, BootstrapState, void>>;
  fetchMatchProfile: ActionCreator<ThunkAction<Promise<MatchProfileActionTypes>, MatchProfileState, void>>;
}

interface Props extends BootstrapState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

interface State {
  refreshing: boolean;
}

enum Directions {
  OUTGOING = 'OUTGOING',
  INCOMING = 'INCOMING',
}

class RequestsView extends Component<Props, State> {
  REQUESTS_VIEW_IDENTIFIER = "RequestsView";

  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: 'Requests',
    headerStyle, 
    headerTitleStyle, 
    headerTintColor 
  })

  constructor(props: Props) {
    super(props);

    this.state = { refreshing: false };

    this.load = this.load.bind(this);
    this.renderHome = this.renderHome.bind(this);
    this.onRefresh = this.onRefresh.bind(this);
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.REQUESTS_VIEW_IDENTIFIER);
    });
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

  private renderRequest(
    requestAndDir: [ConnectionRequestWithName, Directions],
    idx: number,
    arr: Immutable.List<ConnectionRequestWithName>,
  ) {
    const [ request, dir ] = requestAndDir
    const {
      userId,
      firstName,
      lastName,
      intentType,
      searchedTrait,
      message,
    } = request;

    const onCloseAccept = async () => {
      try {
        await requestToMatchService.removeConnection(userId);
        this.props.infoToast('Removed request');
        this.props.fetchBootstrap();
      } catch (e) {
        await this.props.errorToast(e.errorMsg);
      }
    }

    const onClosePress = () => {
      Alert.alert(
        'Remove request',
        'Are you sure you want to remove this request?',
        [
          {text: 'Cancel', onPress: () => null, style: 'cancel'},
          {text: 'Remove', onPress: onCloseAccept, style: 'destructive'},
        ],
      );
    }

    const closeButton = (
      <TouchableOpacity style={styles.deleteRtmMatch} onPress={onClosePress}>
        <MaterialIcons name="close" size={18} />
      </TouchableOpacity>
    );

    const onAccept = async () => {
      try {
        await requestToMatchService.acceptConnection(userId);
        this.props.infoToast('Accepted connection');
        this.props.fetchBootstrap();
      } catch (e) {
        await this.props.errorToast(e.errorMsg);
      }
    }

    const onAcceptPress = () => {
      Alert.alert(
        'Add Connection',
        'Are you sure you want to add ' + firstName + ' ' + lastName + ' as a connection? They ' +
        'will be able to connect with you by email, text, etc.',
        [
          {text: 'Cancel', onPress: () => null, style: 'cancel'},
          {text: 'Add', onPress: onAccept, style: 'default'},
        ],
      );
    }

    const acceptButton = dir === Directions.INCOMING ? (
      <Button
        buttonStyle={{ width: 150 }}
        color={Colors.GREEN}
        title="Accept"
        onPress={onAcceptPress}
      />
    ) : null;

    const extraStyle = arr.size === idx + 1 ? null : {marginBottom: 10};
    let description: ReactNode = null;
    switch (intentType) {
      case IntentTypes.REC_COHORT:
        description = <Text>They are in your cohort</Text>;
        break;
      case IntentTypes.REC_GENERAL:
        if (dir === Directions.INCOMING) {
          description = <Text>You were recommended to them</Text>;
        } else {
          description = <Text>They were recommended to you</Text>;
        }
        break;
      case IntentTypes.SEARCH:
        if (dir === Directions.OUTGOING) {
          description = (
            <Text>
              {'You searched for '}
              <Text style={styles.bold}>{ searchedTrait }</Text>
            </Text>
          );
        } else {
          description = (
            <Text>
              {'They searched for '}
              <Text style={styles.bold}>{ searchedTrait }</Text>
            </Text>
          );
        }
        break;
      default:
        const _: never = intentType
    }

    let messageElement: ReactNode = null;
    if (!!message) {
      messageElement = (
        <Text style={styles.message}>
          <Text style={styles.bold}>{'Message: '}</Text>
          { message }
        </Text>
      );
    }

    const viewProfile = () => {
      this.props.navigation.navigate('MatchProfile', { userId, connectionIntent: request });
    }

    return (
      <Card key={userId} style={[styles.connectionCard, extraStyle]}>
        <View style={{ flexDirection: 'row' }}>
          <View style={styles.cardProfilePicture}>
            <ProfileAvatar userId={userId.toString()} large/>
          </View>
          <View style={{ flex: 2 }}>
            <Text style={styles.name}>{firstName + ' ' + lastName}</Text>
            <Text style={styles.description}>{ description }</Text>
          </View>
        </View>
        { messageElement }
        <View style={{ marginTop: 10, flexDirection: 'row', justifyContent: 'space-between' }}>
          <TouchableOpacity onPress={ viewProfile }>
            <Text style={{ color: Colors.HIVE_PRIMARY, fontSize: 18, paddingTop: 2 }}>
              View Profile
            </Text>
          </TouchableOpacity>
          { acceptButton }
        </View>
        { closeButton }
      </Card>
    );
  }

  private renderHome() {
    const {
      incomingRequests,
      outgoingRequests,
    } = this.props.bootstrap.connections;

    const elements: Array<ReactNode> = [];
    if (!incomingRequests.isEmpty()) {
      elements.push(<Header key={'incoming'}>Incoming Requests</Header>);
      elements.push(incomingRequests.map((req, idx, lst) => {
        return this.renderRequest([req, Directions.INCOMING], idx, lst.toList());
      }).toJS());
    }
    if (!outgoingRequests.isEmpty()) {
      elements.push(<Header key={'outgoing'}>Outgoing Requests</Header>);
      elements.push(outgoingRequests.map((req, idx, lst) => {
        return this.renderRequest([req, Directions.OUTGOING], idx, lst.toList());
      }).toJS());
    }

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
            { elements }
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
  { errorToast, infoToast, fetchBootstrap, fetchMatchProfile })(RequestsView);

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  connectionCard: {
    marginHorizontal: 0,
    marginVertical: 0,
  },
  cardProfilePicture: {
    flex: 1,
    marginRight: 10,
  },
  name: {
    fontWeight: 'bold',
    fontSize: 20,
  },
  scrollContainer: {
    padding: 10,
  },
  description: {
    fontSize: 14,
    color: 'gray',
    marginBottom: 5,
  },
  message: {
    marginTop: 5,
    fontSize: 14,
    color: 'gray',
  },
  bold: {
    fontWeight: 'bold',
  },
  deleteRtmMatch: {
    position: 'absolute',
    top: 10,
    right: 10,
  },
})
