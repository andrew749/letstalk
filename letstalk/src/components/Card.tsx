import React, { ReactNode } from 'react';
import { Dimensions, StyleSheet, View, ViewStyle } from 'react-native';

const SCREEN_WIDTH = Dimensions.get('window').width;

interface Props {
  children?: ReactNode;
  style?: ViewStyle;
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
    backgroundColor: 'white',
    padding: PADDING,
    marginHorizontal: MARGIN_HOR,
    marginVertical: MARGIN_VER,
    borderRadius: 10,
  },
});
