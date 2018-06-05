import React, { Component } from 'react';
import { ThunkAction } from 'redux-thunk';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import {
  Dimensions,
  StyleSheet,
  View,
  Platform,
} from 'react-native';
import { SearchBar } from 'react-native-elements';
import { Constants } from 'expo';

import { RootState } from '../redux';
import {
  State as SearchBarState,
  updateValue,
  updateFocus,
  updateListType
} from '../redux/search-bar/reducer';
import {
  ActionTypes as SearchBarActionTypes,
  SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS,
  SEARCH_LIST_TYPE_CREDENTIALS,
} from '../redux/search-bar/actions';
import Colors from '../services/colors';

const SCREEN_WIDTH = Dimensions.get('window').width;

interface DispatchActions {
  updateValue: ActionCreator<
    ThunkAction<Promise<SearchBarActionTypes>, SearchBarState, void>>;
  updateFocus: ActionCreator<
    ThunkAction<Promise<SearchBarActionTypes>, SearchBarState, void>>;
  updateListType: ActionCreator<
    ThunkAction<Promise<SearchBarActionTypes>, SearchBarState, void>>;
}

interface Props extends DispatchActions, SearchBarState {}

class TopHeader extends Component<Props> {
  private searchBar: SearchBar;

  constructor(props: Props) {
    super(props);

    this.searchBar = null;
  }

  componentWillReceiveProps(props: Props) {
    if (this.props.hasFocus && !props.hasFocus) {
      this.searchBar.blur();
    } else if (!this.props.hasFocus && props.hasFocus) {
      this.searchBar.focus();
    }
  }

  render() {
    // Only show clear icon if the element has focus
    const clearIcon = this.props.hasFocus ? {
      color: '#86939e',
      name: 'close',
      style: [styles.icon, styles.rightIcon],
    } : null;

    let placeholder = '';
    switch (this.props.listType) {
      case SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS:
        placeholder = 'Find someone who is a...';
        break;
      case SEARCH_LIST_TYPE_CREDENTIALS:
        placeholder = 'I am a...';
        break
      default:
        // Ensure exhaustiveness of select
        const _: never = this.props.listType;
    }

    return <View style={styles.header}>
      <SearchBar
        round
        lightTheme
        ref={(ref: SearchBar) => this.searchBar = ref}
        clearIcon={clearIcon}
        icon={{ style: [styles.icon, styles.leftIcon] }}
        onChangeText={(value: string) => this.props.updateValue(value)}
        onClearText={() => this.searchBar.blur()}
        containerStyle={styles.searchBarContainer}
        inputStyle={styles.searchBarTextInput}
        value={this.props.value}
        placeholder={placeholder}
        placeholderTextColor={Colors.HIVE_LIGHT_FONT}
        onFocus={() => this.props.updateFocus(true)}
        onBlur={() => {
          this.props.updateFocus(false);
          this.props.updateValue('');
          this.props.updateListType(SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS);
        }}
      />
    </View>;
  }
}

export default connect(({ searchBar }: RootState) => searchBar,
  { updateValue, updateFocus, updateListType })(TopHeader);

const styles = StyleSheet.create({
  searchBarTextInput: {
    margin: 4, // from a default of 8
    height: 30,
    fontSize: 16,
    backgroundColor: Colors.HIVE_PRIMARY_LIGHT,
    color: Colors.HIVE_MAIN_FONT,
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
    width: SCREEN_WIDTH,
    backgroundColor: Colors.HIVE_PRIMARY,
    borderBottomWidth: 0,
    borderTopWidth: 0,
  },
  header: {
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