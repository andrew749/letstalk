import React, { Component, ReactNode } from 'react';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import {
  ActivityIndicator,
  Button as ReactNativeButton,
  Dimensions,
  Linking,
  RefreshControl,
  RefreshControlProps,
  StyleSheet,
  Text,
  TouchableOpacity,
  View,
  ScrollView,
  Platform,
} from 'react-native';
import {
  NavigationScreenProp,
  NavigationScreenDetails,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import { MaterialIcons, MaterialCommunityIcons } from '@expo/vector-icons';

import { RootState } from '../redux';
import { errorToast } from '../redux/toast';
import {
  ActionTypes as UserSearchActionTypes,
  QueryTypes,
} from '../redux/user-search/actions';
import {
  searchByCohort,
  searchByMyCohort,
  searchByPosition,
  searchBySimpleTrait,
  State as UserSearchState,
} from '../redux/user-search/reducer';
import {
  ActionButton,
  Button,
  Card,
  Header,
  ProfileAvatar,
} from '../components';
import Loading from './Loading';
import Colors from '../services/colors';
import TopHeader, { headerStyle, headerTitleStyle, headerTintColor  } from './TopHeader';
import AllFilterableModals from './AllFilterableModals';
import { AnalyticsHelper } from '../services';
import { DEFAULT_SEARCH_SIZE } from '../services/user-search-service';
import { UserSearchResult } from '../models/user-search';
import { IntentTypes } from '../models/connection';

const SCREEN_WIDTH = Dimensions.get('window').width;

interface DispatchActions {
  searchByCohort: ActionCreator<ThunkAction<Promise<UserSearchActionTypes>, UserSearchState, void>>;
  searchByMyCohort:
    ActionCreator<ThunkAction<Promise<UserSearchActionTypes>, UserSearchState, void>>;
  searchByPosition:
    ActionCreator<ThunkAction<Promise<UserSearchActionTypes>, UserSearchState, void>>;
  searchBySimpleTrait:
    ActionCreator<ThunkAction<Promise<UserSearchActionTypes>, UserSearchState, void>>;
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

interface Props extends DispatchActions, UserSearchState {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

interface State {
  refreshing: boolean;
}

class ExploreView extends Component<Props, State> {
  EXPLORE_VIEW_IDENTIFIER = "ExploreView";

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
    this.renderBody = this.renderBody.bind(this);
    this.onRefresh = this.onRefresh.bind(this);
    this.renderResult = this.renderResult.bind(this);
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.EXPLORE_VIEW_IDENTIFIER);
    });
    if (this.props.fetchState.state === 'prefetch' &&
      this.props.currentQuery.type === QueryTypes.YOUR_COHORT) {
      this.props.searchByMyCohort(DEFAULT_SEARCH_SIZE);
    }
  }

  private async load() {
    const {
      currentQuery,
    } = this.props;

    switch (currentQuery.type) {
      case QueryTypes.YOUR_COHORT:
        await this.props.searchByMyCohort(DEFAULT_SEARCH_SIZE);
        break;
      case QueryTypes.SEARCH_COHORT:
        await this.props.searchByCohort(currentQuery.cohortId, DEFAULT_SEARCH_SIZE);
        break;
      case QueryTypes.SEARCH_POSITION:
        await this.props.searchByPosition(currentQuery.roleId, currentQuery.organizationId,
          DEFAULT_SEARCH_SIZE);
        break;
      case QueryTypes.SEARCH_SIMPLE_TRAIT:
        await this.props.searchBySimpleTrait(currentQuery.simpleTraitId, DEFAULT_SEARCH_SIZE);
        break;
      default:
        const _: never = currentQuery;
    }
  }

  private async onRefresh() {
    this.setState({refreshing: true});
    await this.load();
    this.setState({refreshing: false});
  }

  private renderQueryInfo() {
    const {
      currentQuery,
    } = this.props;

    let header: ReactNode = null;
    switch (currentQuery.type) {
      case QueryTypes.YOUR_COHORT:
        header = (
          <Text>
            <Text>{'Users in '}</Text>
            <Text style={styles.boldText}>Your Cohort</Text>
          </Text>
        );
        break;
      case QueryTypes.SEARCH_COHORT:
        const {
          programName,
          sequenceName,
          gradYear,
        } = currentQuery;
        let cohortText = programName + ' ' + gradYear;
        if (!!sequenceName) cohortText = cohortText + ' ' + sequenceName;
        header = (
          <Text>
            <Text>{'Users in "'}</Text>
            <Text style={styles.boldText}>{ cohortText }</Text>
            <Text>{'"'}</Text>
          </Text>
        );
        break;
      case QueryTypes.SEARCH_POSITION:
        const {
          roleName,
          organizationName,
        } = currentQuery;
        let positionText = roleName + ' @ ' + organizationName;
        header = (
          <Text>
            <Text>{'Users with the position "'}</Text>
            <Text style={styles.boldText}>{ positionText }</Text>
            <Text>{'"'}</Text>
          </Text>
        );
        break;
      case QueryTypes.SEARCH_SIMPLE_TRAIT:
        const { simpleTraitName } = currentQuery;
        header = (
          <Text>
            <Text>{'Users with the trait "'}</Text>
            <Text style={styles.boldText}>{ simpleTraitName }</Text>
            <Text>{'"'}</Text>
          </Text>
        );
        break;
      default:
        const _: never = currentQuery;
    }

    const { numResults } = this.props.response;
    const numResText = 'Found ' +
      (numResults === DEFAULT_SEARCH_SIZE ? 'at least ' : '') + numResults + ' user' +
      (numResults === 1 ? '' : 's');

    return (
      <View>
        <Text style={styles.sectionHeader}>{ header }</Text>
        <Text style={styles.numResultsText}>{ numResText }</Text>
      </View>
    );
  }

  private renderResult(res: UserSearchResult) {
    const { userId, firstName, lastName, cohort, profilePic, reason } = res;

    let cohortText = null;
    if (!!cohort) {
      cohortText = cohort.programName + ' ' + cohort.gradYear;
      if (!!cohort.sequenceName) cohortText = cohortText + ' ' + cohort.sequenceName;
    }

    const onPress = () => {
      const { currentQuery } = this.props;
      let searchedTrait: string = null;
      let intentType: IntentTypes = IntentTypes.REC_COHORT;
      switch (currentQuery.type) {
        case QueryTypes.YOUR_COHORT:
          intentType = IntentTypes.REC_COHORT;
          break;
        case QueryTypes.SEARCH_COHORT:
          intentType = IntentTypes.SEARCH;
          const {
            programName,
            sequenceName,
            gradYear,
          } = currentQuery;
          searchedTrait = programName + ' ' + gradYear;
          if (!!sequenceName) searchedTrait = searchedTrait + ' ' + sequenceName;
          break;
        case QueryTypes.SEARCH_POSITION:
          intentType = IntentTypes.SEARCH;
          const {
            roleName,
            organizationName,
          } = currentQuery;
          searchedTrait = roleName + ' @ ' + organizationName;
          break;
        case QueryTypes.SEARCH_SIMPLE_TRAIT:
          intentType = IntentTypes.SEARCH;
          const { simpleTraitName } = currentQuery;
          searchedTrait = simpleTraitName;
          break;
        default:
          const _: never = currentQuery;
      }

      const connectionIntent = { intentType, searchedTrait };
      this.props.navigation.navigate('MatchProfile', { userId, connectionIntent });
    }

    return (
      <TouchableOpacity onPress={onPress} key={userId} style={styles.resultContainer}>
        <ProfileAvatar medium userId={userId.toString()} />
        <View style={styles.descriptionContainer}>
          <Text style={styles.nameText}>{ firstName } { lastName }</Text>
          {cohortText && <Text style={styles.cohortText}>{ cohortText }</Text>}
          {reason && <Text style={styles.reasonText}>{ reason }</Text>}
        </View>
        <MaterialIcons
          style={styles.rightArrow}
          name="chevron-right"
          size={32}
          color={Colors.HIVE_SUBDUED}
        />
      </TouchableOpacity>
    );
  }

  private renderResults() {
    const { results, numResults, isAnonymous } = this.props.response;
    if (isAnonymous) {
      return (
        <View style={styles.explanationContainer}>
          <MaterialIcons name="priority-high" size={64} color={Colors.HIVE_SUBDUED} />
          <Text style={styles.footerText}>
            You are searching for a sensitive trait. In order to
            protect users' privacy, we currently do not support searching for users by traits that
            are sensitive.
          </Text>
        </View>
      );
    } else if (numResults === 0) {
      return (
        <View style={styles.explanationContainer}>
          <MaterialCommunityIcons name="emoticon-sad" size={64} color={Colors.HIVE_SUBDUED} />
          <Text style={styles.footerText}>No users found. Try again later.</Text>
        </View>
      );
    }

    const userCards = results.map(this.renderResult);
    const endText = numResults === DEFAULT_SEARCH_SIZE ?
      'There may be more than ' + numResults + ' users. Refresh by pulling down to see other users'
      : '';

    return (
      <View>
        <View style={styles.resultsContainer}>{ userCards }</View>
        {endText && <View style={styles.footerContainer}>
          <Text style={styles.footerText}>
            { endText }
          </Text>
        </View>}
      </View>
    );
  }

  private renderBody() {
    return (
      <ScrollView
        style={styles.container}
        keyboardShouldPersistTaps={'always'}
        refreshControl={
          <RefreshControl
            refreshing={this.state.refreshing}
            onRefresh={this.onRefresh}
          /> as React.ReactElement<RefreshControlProps>
        }
      >
        <View style={styles.headerContainer}>{ this.renderQueryInfo() }</View>
        { this.renderResults() }
      </ScrollView>
    );
  }

  render() {
    const {
      errorMsg,
      errorType,
      state,
    } = this.props.fetchState;

    return (
      <View style={{flex: 1}}>
        <Loading
          state={this.state.refreshing ? 'success' : state}
          errorMsg={errorMsg}
          errorType={errorType}
          load={this.load}
          renderBody={this.renderBody}
          navigation={this.props.navigation}
        />
        <AllFilterableModals />
      </View>
    );
  }
}

export default connect(
  ({ userSearch }: RootState) => {
    return userSearch;
  }, {
    errorToast,
    searchByCohort,
    searchByMyCohort,
    searchByPosition,
    searchBySimpleTrait,
  })(ExploreView);

const styles = StyleSheet.create({
  container: {
  },
  headerContainer: {
    padding: 10,
  },
  sectionHeader: {
    fontSize: 18,
    paddingBottom: 5,
  },
  boldText: {
    fontWeight: '900',
  },
  numResultsText: {
    color: Colors.HIVE_SUBDUED,
  },
  resultContainer: {
    flexDirection: 'row',
    padding: 10,
    alignItems: 'center',
    borderBottomWidth: 0.5,
    borderColor: Colors.HIVE_SUBDUED,
  },
  resultsContainer: {
    backgroundColor: Colors.WHITE,
    borderTopWidth: 0.5,
    borderColor: Colors.HIVE_SUBDUED,
  },
  nameText: {
    fontSize: 16,
  },
  cohortText: {
    fontSize: 14,
    color: Colors.HIVE_SUBDUED,
  },
  reasonText: {
    fontSize: 14,
    color: Colors.HIVE_ACCENT,
  },
  descriptionContainer: {
    flexDirection: 'column',
    width: SCREEN_WIDTH - 100,
    paddingLeft: 10,
  },
  rightArrow: {
    position: 'absolute',
    right: 10,
  },
  footerContainer: {
    alignItems: 'center',
    padding: 10,
  },
  footerText: {
    fontSize: 14,
    color: Colors.HIVE_SUBDUED,
  },
  explanationContainer: {
    paddingHorizontal: 10,
    alignItems: 'center',
  },
})
