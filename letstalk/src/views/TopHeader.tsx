import React, { SFC } from 'react';

import {
  Dimensions,
  StyleSheet,
  View,
  Platform,
} from 'react-native';
import { SearchBar } from 'react-native-elements';
import { Constants } from 'expo';

import Colors from '../services/colors';

const SCREEN_WIDTH = Dimensions.get('window').width;

const TopHeader: SFC<{}> = (props: {}) => {
  return <View style={styles.header}>
    <SearchBar
      round
      lightTheme
      clearIcon={{ color: '#86939e', name: 'close', style: [styles.icon, styles.rightIcon] }}
      icon={{ style: [styles.icon, styles.leftIcon] }}
      onClearText={() => console.log('cleared')}
      containerStyle={styles.searchBarContainer}
      inputStyle={styles.searchBarTextInput}
    />
  </View>;
}

const styles = StyleSheet.create({
  searchBarTextInput: {
    margin: 4, // from a default of 8
    height: 30,
    fontSize: 16,
  },
  icon: {
    top: 11.5, // 15.5 - 4 (hard coded top - decrease in margin)
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

export default TopHeader;

export { headerStyle };
