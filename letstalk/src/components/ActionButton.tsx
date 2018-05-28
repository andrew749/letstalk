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
      containerViewStyle={[styles.containerStyle, props.containerStyle]}
    />
  );
};

export default ActionButton;

const styles = StyleSheet.create({
  containerStyle: {
    width: "90%"
  },
  loginButtonStyle: {
    height: 55,
    borderRadius: 5,
    margin: 10,
  },
  loginButtonTextStyle: {
    fontSize: 20,
    color: 'white',
  },
});
