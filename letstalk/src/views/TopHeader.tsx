import React, { Component } from 'react';
import { ThunkAction } from 'redux-thunk';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import {
  Dimensions,
  StyleSheet,
  View,
  Platform,
  TouchableOpacity,
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
} from '../redux/search-bar/reducer';
import {
  ActionTypes as SearchBarActionTypes,
  SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS,
  SEARCH_LIST_TYPE_CREDENTIALS,
} from '../redux/search-bar/actions';
import Colors from '../services/colors';
import { MaterialIcons } from '@expo/vector-icons';
import { logAnalyticsThenExecute, AnalyticsActions } from '../services/analytics';
import autocompleteService from '../services/autocomplete-service';
import { MultiTrait } from '../models/multi-trait';

const SCREEN_WIDTH = Dimensions.get('window').width;
const THROTTLE_TIME = 250; // ms

const onQueryChange = async (query: string, setData: (data: Immutable.List<MultiTrait>) => void) => {
  let res: Immutable.List<MultiTrait> = Immutable.List();
  if (query !== '') res = await autocompleteService.autocompleteMultiTrait(query, 10);
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
      this.searchBar.blur();
    } else if (!this.props.hasFocus && props.hasFocus) {
      this.searchBar.focus();
    }
  }

  private onFocusAfterLog() {
    this.props.updateFocus(true);
  }

  private onFocus() {
    logAnalyticsThenExecute("SearchBar", AnalyticsActions.FOCUS, "", 1, this.onFocusAfterLog);
  }

  private onBlur() {
    this.props.updateFocus(false);
    this.props.updateValue('');
  }

  private onChangeText(value: string) {
    this.props.updateValue(value);
    onQueryChangeThrottled(value, this.props.updateSuggestions);
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

    let placeholder = 'Find someone...';

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
          placeholderTextColor={Colors.HIVE_LIGHT_FONT}
          onFocus={this.onFocus}
          onBlur={this.onBlur}
        />
        <TouchableOpacity style={styles.qrButton} onPress={openQr}>
          <MaterialIcons name="camera-enhance" color={Colors.HIVE_PRIMARY_LIGHT} size={24} />
        </TouchableOpacity>
        <TouchableOpacity style={styles.notifButton} onPress={openNotif}>
          <MaterialIcons name="notifications" color={Colors.HIVE_PRIMARY_LIGHT} size={24} />
        </TouchableOpacity>
      </View>
    );
  }
}

export default connect(({ searchBar }: RootState) => searchBar,
  { updateValue, updateFocus, updateSuggestions })(TopHeader);

const SEARCH_BAR_LEFT_MARGIN = 36;
const SEARCH_BAR_RIGHT_MARGIN = 36;

const styles = StyleSheet.create({
  searchBarTextInput: {
    margin: 4, // from a default of 8
    height: 30,
    fontSize: 16,
    backgroundColor: Colors.HIVE_PRIMARY_LIGHT,
    color: Colors.HIVE_MAIN_FONT,
  },
  qrButton: {
    position: 'absolute',
    padding: 8,
  },
  notifButton: {
    position: 'absolute',
    padding: 8,
    right: 0,
  },
  icon: {
    top: 11.5, // 15.5 - 4 (hard coded top - decrease in margin)
    color: Colors.HIVE_MAIN_FONT,
  },
  leftIcon: {
    left: 12, // 16 - 4 (hard coded left - decrease in margin)
  },
  rightIcon: {
    right: 12, // 16 - 4 (hard coded right - decrease in margin)
  },
  searchBarContainer: {
    width: SCREEN_WIDTH - SEARCH_BAR_LEFT_MARGIN - SEARCH_BAR_RIGHT_MARGIN,
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
  }
})

const headerStyle = styles.topLevelHeader;

export { headerStyle };
