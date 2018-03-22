import React, { Component, ReactNode } from 'react';
import { StyleSheet, Text } from 'react-native';

interface Props {
  children: ReactNode
}

const InfoText: React.SFC<Props> = props => {
  const { children } = props
  return <Text style={styles.text}>{children}</Text>;
};

export default InfoText;

const styles = StyleSheet.create({
  text: {
    padding: 10,
    fontWeight: '500',
    fontSize: 18,
  },
});
