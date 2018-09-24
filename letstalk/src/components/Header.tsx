import React, { SFC, ReactNode } from 'react';
import { 
  StyleProp, 
  StyleSheet, 
  Text,
  TextStyle
} from 'react-native';

interface Props {
  children: ReactNode;
  textStyle?: StyleProp<TextStyle>;
}

const Header: SFC<Props> = props => {
  const { children } = props
  return (
    <Text style={[styles.text, props.textStyle]}>{children}</Text>
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
