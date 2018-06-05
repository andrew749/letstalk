import React, { SFC } from 'react';
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
} from '../redux/search-bar/reducer';
import { ActionTypes as SearchBarActionTypes } from '../redux/search-bar/actions';
import Colors from '../services/colors';

const SCREEN_WIDTH = Dimensions.get('window').width;

interface DispatchActions {
  updateValue: ActionCreator<
    ThunkAction<Promise<SearchBarActionTypes>, SearchBarState, void>>;
}

interface Props extends DispatchActions, SearchBarState {
  placeholder: string;
}

const TopHeader: SFC<Props> = (props: Props) => {
  // Only show clear icon if there is actually text to be cleared
  const clearIcon = props.value === '' ? null : {
    color: '#86939e',
    name: 'close',
    style: [styles.icon, styles.rightIcon],
  };

  return <View style={styles.header}>
    <SearchBar
      round
      lightTheme
      clearIcon={clearIcon}
      icon={{ style: [styles.icon, styles.leftIcon] }}
      onChangeText={(value: string) => props.updateValue(value)}
      onClearText={() => props.updateValue('')}
      containerStyle={styles.searchBarContainer}
      inputStyle={styles.searchBarTextInput}
      value={props.value}
      placeholder={props.placeholder}
      placeholderTextColor={Colors.HIVE_LIGHT_FONT}
    />
  </View>;
}

export default connect(({ searchBar }: RootState) => searchBar, { updateValue })(TopHeader);

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
