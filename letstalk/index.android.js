import React, { Component } from 'react';
import {
  AppRegistry,
  StyleSheet,
  Text,
  View
} from 'react-native';

import App from './src';

export default class letsTalk extends letsTalk {
  render() {
    return (
      <App />
    );
  }
}

AppRegistry.registerComponent('letsTalk', () => letsTalk);