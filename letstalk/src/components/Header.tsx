import React, { SFC, ReactNode } from 'react';
import { StyleSheet, Text } from 'react-native';

interface Props {
  children: ReactNode;
}

const Header: SFC<Props> = props => {
  const { children } = props
  return (
    <Text style={styles.text}>{children}</Text>
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
