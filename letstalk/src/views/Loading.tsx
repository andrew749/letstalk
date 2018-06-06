import React, { SFC, ReactElement } from 'react';
import { ActivityIndicator, StyleSheet, Text, View } from 'react-native';
import {
  NavigationActions,
  NavigationScreenProp,
  NavigationStackAction,
} from 'react-navigation';

import auth from '../services/auth';
import { ActionButton } from '../components';
import { ErrorTypes } from '../services/requests';

interface Props {
  state: 'prefetch' | 'fetching' | 'error' | 'success';
  errorMsg: string | null;
  errorType: ErrorTypes | null;
  renderBody(): ReactElement<any>;
  load(): Promise<void>;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

const Loading: SFC<Props> = props => {
  const { state, errorMsg, errorType, load, navigation, renderBody } = props;

  switch (state) {
    case 'prefetch':
    case 'fetching':
      return (
        <View style={styles.centeredContainer}>
          <Text style={styles.headline}>Soon...</Text>
          <ActivityIndicator size="large" />
        </View>
      );
    case 'error':
      const isUnauthorized = !!errorType && errorType === 'UNAUTHORIZED';
      const headline = isUnauthorized ?
        'Your session token has expired' : 'Something went wrong :(';

      const logout = async () => {
        try {
          await auth.logout();
        } catch (error) { }
        navigation.dispatch(NavigationActions.reset({
          index: 0,
          key: null,
          actions: [NavigationActions.navigate({ routeName: 'Login' })]
        }));
      }
      const logoutButton = isUnauthorized ? <ActionButton onPress={logout} title="Logout" /> : null;
      const error = !isUnauthorized ? <Text style={styles.error}>{errorMsg}</Text> : null;

      return (
        <View style={styles.centeredContainer}>
          <Text style={styles.headline}>{headline}</Text>
          {error}
          <ActionButton onPress={load} title="Retry" />
          {logoutButton}
        </View>
      );
    case 'success':
      return renderBody();
    default:
      // Ensure exhaustiveness of select
      const _: never = state;
  }
}

export default Loading;

const styles = StyleSheet.create({
  centeredContainer: {
    flex: 1,
    justifyContent: 'center'
  },
  headline: {
    fontWeight: 'bold',
    fontSize: 18,
    textAlign: 'center',
    marginBottom: 10,
  },
  error: {
    fontWeight: 'bold',
    fontSize: 18,
    color: 'red',
    textAlign: 'center',
  },
});
