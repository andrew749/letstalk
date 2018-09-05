import React, { Component, SFC } from 'react';
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
import { MaterialIcons } from '@expo/vector-icons';

import { RootState } from '../redux';
import { errorToast } from '../redux/toast';
import {
  ActionTypes as SearchBarActionTypes,
  SEARCH_LIST_TYPE_CREDENTIALS,
  SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS,
} from '../redux/search-bar/actions';
import {
  State as SearchBarState,
  updateFocus,
} from '../redux/search-bar/reducer';
import {
  ActionButton,
  Button,
  Card,
  Header,
} from '../components';
import Loading from './Loading';
import Colors from '../services/colors';
import TopHeader, { headerStyle } from './TopHeader';
import AllFilterableModals from './AllFilterableModals';
import { AnalyticsHelper } from '../services';

const SCREEN_WIDTH = Dimensions.get('window').width;

interface DispatchActions {
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  updateFocus: ActionCreator<
    ThunkAction<Promise<SearchBarActionTypes>, SearchBarState, void>>;
}

interface Props extends DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

interface State {
  refreshing: boolean;
}

class SearchView extends Component<Props, State> {
  SEARCH_VIEW_IDENTIFIER = "SearchView";

  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: <TopHeader navigation={navigation} />,
    headerStyle,
  })

  constructor(props: Props) {
    super(props);

    this.state = { refreshing: false };

    this.load = this.load.bind(this);
    this.renderBody = this.renderBody.bind(this);
    this.onRefresh = this.onRefresh.bind(this);
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.SEARCH_VIEW_IDENTIFIER);
    });
    this.load();
  }

  private async load() {
  }

  private async onRefresh() {
    this.setState({refreshing: true});
    await this.load();
    this.setState({refreshing: false});
  }

  private renderBody() {
  }

  render() {
    return <Text>yo</Text>;
    // If `this.state.refreshing` is true, it means that we are reloading data using the pull
    // down, which means that we want to still display the ScrollView.
    // return (
    //   <Loading
    //     state={this.state.refreshing ? 'success' : state}
    //     errorMsg={errorMsg}
    //     errorType={errorType}
    //     load={this.load}
    //     renderBody={this.renderBody}
    //     navigation={this.props.navigation}
    //   />
    // );
  }
}

export default connect(
  ({ credentialRequests, credentialOptions, credentials }: RootState) => {
    return { };
  }, {
    errorToast,
    updateFocus,
  })(SearchView);

const styles = StyleSheet.create({
})
