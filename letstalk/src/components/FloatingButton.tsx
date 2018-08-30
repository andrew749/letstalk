import React from 'react';
import { TouchableHighlight, Dimensions, StyleSheet, TextStyle, ViewStyle, StyleProp } from 'react-native';
import { Button, ButtonProps } from 'react-native-elements';
import Colors from '../services/colors';

const SCREEN_WIDTH = Dimensions.get('window').width;

interface Props extends ButtonProps {
  textStyle?: StyleProp<TextStyle>,
  containerStyle?: StyleProp<ViewStyle>,
  buttonStyle?: StyleProp<ViewStyle>,
}

const FloatingButton: React.SFC<Props> = props => {
  return (
    <Button
      {...props}
      buttonStyle={[styles.loginButtonStyle, props.buttonStyle]}
      textStyle={[styles.loginButtonTextStyle, props.textStyle]}
      containerViewStyle={[styles.containerStyle, props.containerStyle]}
      component={TouchableHighlight}
    />
  );
};

export default FloatingButton;

const styles = StyleSheet.create({
  containerStyle: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
  },
  loginButtonStyle: {
    height: 55,
    borderRadius: 5,
    margin: 10,
    backgroundColor: Colors.HIVE_PRIMARY,
  },
  loginButtonTextStyle: {
    fontSize: 20,
    color: 'white',
  },
});
