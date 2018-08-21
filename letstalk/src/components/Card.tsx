import React, { ReactNode } from 'react';
import { Dimensions, StyleProp, StyleSheet, View, ViewStyle } from 'react-native';

import Colors from '../services/colors';

interface Props {
  children?: ReactNode;
  style?: StyleProp<ViewStyle>;
}

const Card: React.SFC<Props> = props => {
  const { children, style } = props;
  return (
    <View style={[styles.card, style]}>{children}</View>
  );
};

export default Card;

// Defaults
const MARGIN_HOR = 20;
const MARGIN_VER = 10;
const PADDING = 20;

const styles = StyleSheet.create({
  card: {
    flex: 1,
    backgroundColor: Colors.WHITE,
    padding: PADDING,
    marginHorizontal: MARGIN_HOR,
    marginVertical: MARGIN_VER,
    borderRadius: 10,
  },
});
