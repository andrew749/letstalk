import React, { Component } from 'react';
import { StyleSheet, Text } from 'react-native';

interface Props {
  title: string;
}

const Header: React.SFC<Props> = props => {
  const { title } = props
  return (
    <Text style={styles.text}>{title}</Text>
  );
};

export default Header;

const styles = StyleSheet.create({
  text: {
    padding: 10,
    fontWeight: "900",
    fontSize: 28,
  },
});
