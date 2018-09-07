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
  ActionTypes as UserSearchActionTypes,
} from '../redux/user-search/actions';
import {
  State as SearchBarState,
} from '../redux/user-search/reducer';
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
}

interface Props extends DispatchActions, SearchBarState {
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
      AnalyticsHelper.getInstance().recordPage(this.EXPLORE_VIEW_IDENTIFIER);
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
    return <Text>a</Text>;
  }

  render() {
    const {
      errorMsg,
      errorType,
      state,
    } = this.props.fetchState;

    console.log(this.props.response);

    return (
      <View style={{flex: 1}}>
        <Loading
          state={state}
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
  })(ExploreView);

const styles = StyleSheet.create({
})
