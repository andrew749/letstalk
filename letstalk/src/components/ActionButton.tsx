import React from 'react';
import { Dimensions, StyleSheet, TextStyle, ViewStyle, StyleProp } from 'react-native';
import { Button, ButtonProps } from 'react-native-elements';

const SCREEN_WIDTH = Dimensions.get('window').width;

interface Props extends ButtonProps {
  textStyle?: StyleProp<TextStyle>,
  containerStyle?: StyleProp<ViewStyle>,
  buttonStyle?: StyleProp<ViewStyle>,
}

const ActionButton: React.SFC<Props> = props => {
  return (
    <Button
      {...props}
      buttonStyle={[styles.loginButtonStyle, props.buttonStyle]}
      textStyle={[styles.loginButtonTextStyle, props.textStyle]}
    />
  );
};

export default ActionButton;

const styles = StyleSheet.create({
  loginButtonStyle: {
    height: 55,
    width: SCREEN_WIDTH - 40,
    borderRadius: 5,
    margin: 10,
  },
  loginButtonTextStyle: {
    fontSize: 20,
    color: 'white',
  },
});
