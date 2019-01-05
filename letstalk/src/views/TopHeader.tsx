import React, { Component } from 'react';
import { ThunkAction } from 'redux-thunk';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import {
  Dimensions,
  StyleSheet,
  View,
  Platform,
  TouchableOpacity,
  BackHandler,
} from 'react-native';
import { SearchBar } from 'react-native-elements';
import { Constants } from 'expo';
import {
  NavigationScreenProp,
  NavigationRoute,
  NavigationParams,
} from 'react-navigation';
import _ from 'underscore';
import Immutable from 'immutable';

import { RootState } from '../redux';
import {
  State as SearchBarState,
  updateValue,
  updateFocus,
  updateSuggestions,
  updateError,
} from '../redux/search-bar/reducer';
import {
  ActionTypes as SearchBarActionTypes,
  SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS,
  SEARCH_LIST_TYPE_CREDENTIALS,
} from '../redux/search-bar/actions';
import Colors from '../services/colors';
import Window from '../services/window';
import { MaterialIcons } from '@expo/vector-icons';
import { logAnalyticsThenExecute, AnalyticsActions } from '../services/analytics';
import autocompleteService from '../services/autocomplete-service';
import { MultiTrait } from '../models/multi-trait';

const THROTTLE_TIME = 250; // ms

const onQueryChange = async (
  query: string,
  setData: (data: Immutable.List<MultiTrait>) => void,
  onError?: (e: any) => void,
) => {
  let res: Immutable.List<MultiTrait> = Immutable.List();
  if (query !== '') {
    try {
      res = await autocompleteService.autocompleteMultiTrait(query, 10);
    } catch (e) {
      if (!!onError) onError(e);
      throw e;
    }
  }
  setData(res);
}
const onQueryChangeThrottled = _.throttle(onQueryChange, THROTTLE_TIME);

interface DispatchActions {
  updateValue: ActionCreator<
    ThunkAction<Promise<SearchBarActionTypes>, SearchBarState, void>>;
  updateFocus: ActionCreator<
    ThunkAction<Promise<SearchBarActionTypes>, SearchBarState, void>>;
  updateSuggestions: ActionCreator<
    ThunkAction<Promise<SearchBarActionTypes>, SearchBarState, void>>;
  updateError: ActionCreator<
    ThunkAction<Promise<SearchBarActionTypes>, SearchBarState, void>>;
}

interface Props extends DispatchActions, SearchBarState {
  navigation: NavigationScreenProp<NavigationRoute<NavigationParams>, NavigationParams>;
}

class TopHeader extends Component<Props> {
  private searchBar: SearchBar;

  constructor(props: Props) {
    super(props);

    this.searchBar = null;

    this.onFocus = this.onFocus.bind(this);
    this.onBlur = this.onBlur.bind(this);
    this.onFocusAfterLog = this.onFocusAfterLog.bind(this);
    this.onChangeText = this.onChangeText.bind(this);
    this.onClearText = this.onClearText.bind(this);
  }

  componentWillReceiveProps(props: Props) {
    if (this.props.hasFocus && !props.hasFocus) {
      BackHandler.removeEventListener('hardwareBackPress', this.handleBackPress);
      this.searchBar.blur();
    } else if (!this.props.hasFocus && props.hasFocus) {
      BackHandler.addEventListener('hardwareBackPress', this.handleBackPress);
      this.searchBar.focus();
    }
  }

  private onFocusAfterLog() {
    this.props.updateFocus(true);
  }

  private onFocus() {
    logAnalyticsThenExecute("SearchBar", AnalyticsActions.FOCUS, "", 1, this.onFocusAfterLog);
  }

  private handleBackPress = () => {
    this.props.updateFocus(false);
    return true;
  };

  private onBlur() {
    this.props.updateFocus(false);
    this.props.updateValue('');
    this.props.updateSuggestions(Immutable.List());
  }

  private onChangeText(value: string) {
    this.props.updateValue(value);
    onQueryChangeThrottled(value, data => {
      this.props.updateError(null);
      this.props.updateSuggestions(data);
    }, e => {
      this.props.updateError(e.errorMsg);
    });
  }

  private onClearText() {
    this.searchBar.blur();
  }

  render() {
    // Only show clear icon if the element has focus
    const clearIcon = this.props.hasFocus ? {
      color: '#86939e',
      name: 'close',
      style: [styles.icon, styles.rightIcon],
    } : null;

    let placeholder = 'Meet someone...';

    const openQr = () => {
      this.props.navigation.navigate({routeName: 'QrScanner'});
    };

    const openNotif = () => {
      this.props.navigation.navigate({routeName: 'NotificationView'});
    };

    return (
      <View style={styles.header}>
        <SearchBar
          round
          lightTheme
          ref={(ref: SearchBar) => this.searchBar = ref}
          clearIcon={clearIcon}
          icon={{ style: [styles.icon, styles.leftIcon] }}
          onChangeText={this.onChangeText}
          onClearText={this.onClearText}
          containerStyle={styles.searchBarContainer}
          inputStyle={styles.searchBarTextInput}
          value={this.props.value}
          placeholder={placeholder}
          placeholderTextColor={Colors.WHITE}
          onFocus={this.onFocus}
          onBlur={this.onBlur}
        />
        <TouchableOpacity style={styles.qrButton} onPress={openQr}>
          <MaterialIcons name="camera-enhance" color={Colors.WHITE} size={28} />
        </TouchableOpacity>
        <TouchableOpacity style={styles.notifButton} onPress={openNotif}>
          <MaterialIcons name="notifications" color={Colors.WHITE} size={28} />
        </TouchableOpacity>
      </View>
    );
  }
}

export default connect(({ searchBar }: RootState) => searchBar,
  { updateValue, updateFocus, updateSuggestions, updateError })(TopHeader);

const SEARCH_BAR_LEFT_MARGIN = 50;
const SEARCH_BAR_RIGHT_MARGIN = 50;

const styles = StyleSheet.create({
  searchBarTextInput: {
    margin: 4, // from a default of 8
    height: 30,
    fontSize: 16,
    backgroundColor: Colors.HIVE_PRIMARY_LIGHT,
    color: Colors.WHITE,
  },
  qrButton: {
    position: 'absolute',
    padding: 4.7,
    marginLeft: 10,
    marginRight: 10
  },
  notifButton: {
    position: 'absolute',
    padding: 4.7,
    marginLeft: 10,
    marginRight: 10,
    right: 0,
  },
  icon: {
    top: 11.5, // 15.5 - 4 (hard coded top - decrease in margin)
    color: Colors.WHITE,
  },
  leftIcon: {
    left: 12, // 16 - 4 (hard coded left - decrease in margin)
  },
  rightIcon: {
    right: 12, // 16 - 4 (hard coded right - decrease in margin)
  },
  searchBarContainer: {
    width: Window.WIDTH - SEARCH_BAR_LEFT_MARGIN - SEARCH_BAR_RIGHT_MARGIN,
    backgroundColor: Colors.HIVE_PRIMARY,
    borderBottomWidth: 0,
    borderTopWidth: 0,
    marginLeft: SEARCH_BAR_LEFT_MARGIN,
    marginRight: SEARCH_BAR_RIGHT_MARGIN,
  },
  header: {
    flex: 1,
    flexDirection: 'row',
    backgroundColor: Colors.HIVE_PRIMARY,
  },
  topLevelHeader: {
    backgroundColor: Colors.HIVE_PRIMARY,
    // React Navigation on Android adds extra top margin for nested views
    marginTop: Platform.OS === 'ios' ? 0 : -Constants.statusBarHeight,
  },
  topLevelHeaderTitleStyle: {
    color: Colors.WHITE
  }
})

const headerStyle = styles.topLevelHeader;
const headerTitleStyle = styles.topLevelHeaderTitleStyle;
const headerTintColor = Colors.WHITE;

export { headerStyle, headerTitleStyle, headerTintColor };
