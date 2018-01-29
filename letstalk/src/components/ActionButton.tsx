import React, { Component } from 'react';
import { Dimensions } from 'react-native';
import { Button, ButtonProps } from 'react-native-elements';

const SCREEN_WIDTH = Dimensions.get('window').width;

export default class ActionButton extends Component<ButtonProps> {
  render() {
    return (
      <Button
        style={styles.loginButtonContainerStyle}
        buttonStyle={styles.loginButtonStyle}
        textStyle={styles.loginButtonTextStyle}
        {...this.props}
      />
    );
  }
};

const styles = {
  loginButtonStyle: {
    height: 55,
    width: SCREEN_WIDTH - 40,
    borderRadius: 30,
  },
  loginButtonContainerStyle: {
    marginTop: 10,
  },
  loginButtonTextStyle: {
    fontSize: 20,
    color: 'white',
  },
};
