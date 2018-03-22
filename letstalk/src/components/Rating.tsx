import React, { ReactNode } from 'react';
import {
  Dimensions,
  TouchableOpacity,
  Text,
  View,
  StyleSheet,
  RegisteredStyle,
  ViewStyle,
  TextStyle,
  ImageStyle,
} from 'react-native';
import { WrappedFieldProps } from 'redux-form';

const SCREEN_WIDTH = Dimensions.get('window').width;
const MARGIN = 10;
const HEIGHT = 60;
const BORDER_RADIUS = 20;

type Props = WrappedFieldProps & {
  label: string;
  numElems: number;
};

const styles = StyleSheet.create({
  label: {
    marginLeft: 10,
    marginRight: 10,
    padding: 5,
    fontWeight: 'bold',
    fontSize: 18,
    backgroundColor: 'transparent',
  },
  container: {
    flex: 1,
    flexDirection: 'column',
    backgroundColor: '#F5FCFF',
    margin: MARGIN,
    borderRadius: BORDER_RADIUS,
  },
  labelContainer: {
    flex: 1,
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  opacityContainer: {
    flex: 1,
    flexDirection: 'row',
  },
  firstRect: {
    borderBottomLeftRadius: BORDER_RADIUS,
  },
  lastRect: {
    borderBottomRightRadius: BORDER_RADIUS,
  },
  active: {
    backgroundColor: '#F5CE13',
  },
});

const Rating: React.SFC<Props> = props => {
  const { label, numElems } = props;
  const { onChange, value } = props.input;
  const width = (SCREEN_WIDTH - 2 * MARGIN) / numElems;
  const specificStyles = StyleSheet.create({
    rect: {
      width,
      backgroundColor: '#FFEE9C',
      height: HEIGHT,
    },
  });
  const opacities = (new Array(numElems).fill(0)).map((_, index) => {
    const style: Array<RegisteredStyle<ViewStyle | TextStyle | ImageStyle>> = [specificStyles.rect];
    if (index === 0)            style.push(styles.firstRect);
    if (index === numElems - 1) style.push(styles.lastRect);
    if (index < value)          style.push(styles.active);
    const onPress = index + 1 === value ?
      () => onChange(0) :
      () => onChange(index + 1);
    return <TouchableOpacity activeOpacity={0.8} key={index} style={style} onPress={onPress} />;
  });
  const valueLabel = value + '/' + numElems;
  return (
    <View style={styles.container}>
      <View style={styles.labelContainer}>
        <Text style={styles.label}>{label}</Text>
        <Text style={styles.label}>{valueLabel}</Text>
      </View>
      <View style={styles.opacityContainer}>
        { opacities }
      </View>
    </View>
  );
}

export default Rating;
