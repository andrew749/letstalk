import React, { SFC, ReactElement } from 'react';
import { ActivityIndicator, StyleSheet, Text, View } from 'react-native';

import { ActionButton } from '../components';

interface Props {
  state: 'prefetch' | 'fetching' | 'error' | 'success';
  errorMsg?: string;
  renderBody(): ReactElement<any>;
  load(): Promise<void>;
}

const Loading: SFC<Props> = props => {
  const { state, errorMsg, load, renderBody } = props;

  switch (state) {
    case 'prefetch':
    case 'fetching':
      // TODO: Separate component for loading pages
      return (
        <View style={styles.centeredContainer}>
          <Text style={styles.headline}>Soon...</Text>
          <ActivityIndicator size="large" />
        </View>
      );
    case 'error':
      // TODO: Separate component for error pages
      return (
        <View style={styles.centeredContainer}>
          <Text style={styles.headline}>Something went wrong :(</Text>
          <Text style={styles.error}>{errorMsg}</Text>
          <ActionButton onPress={load} title="Retry" />
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
