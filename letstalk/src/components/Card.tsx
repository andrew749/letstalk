import React, { ReactNode } from 'react';
import { Dimensions, StyleSheet, View } from 'react-native';

const SCREEN_WIDTH = Dimensions.get('window').width;

interface Props {
  children?: ReactNode,
}

const Card: React.SFC<Props> = props => {
  const { children } = props;
  return (
    <View style={styles.card}>{children}</View>
  );
};

export default Card;

const MARGIN_HOR = 20;
const MARGIN_VER = 10;
const PADDING = 20;

const styles = StyleSheet.create({
  card: {
    backgroundColor: 'white',
    width: SCREEN_WIDTH - MARGIN_HOR - PADDING,
    padding: PADDING,
    marginLeft: MARGIN_HOR,
    marginRight: MARGIN_HOR,
    marginTop: MARGIN_VER,
    marginBottom: MARGIN_VER,
    borderRadius: 10,
  },
});
