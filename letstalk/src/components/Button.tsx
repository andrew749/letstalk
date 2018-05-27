import React, { ReactNode } from 'react';
import {
  Dimensions,
  StyleSheet,
  Text,
  TextStyle,
  TouchableOpacity,
  View,
  ViewStyle,
  StyleProp,
} from 'react-native';
import { MaterialIcons } from '@expo/vector-icons';
import { Button as ReactButton, ButtonProps } from 'react-native-elements';
import Colors from '../services/colors';

const SCREEN_WIDTH = Dimensions.get('window').width;

interface Props {
  textStyle?: StyleProp<TextStyle>;
  containerStyle?: StyleProp<ViewStyle>;
  buttonStyle?: StyleProp<ViewStyle>;
  onPress(): void;
  title: string;
  icon?: string;
}

const Button: React.SFC<Props> = props => {
  const icon = props.icon ?
    <MaterialIcons
      style={{ position: 'absolute', left: 2, top: 2 }}
      color={Colors.HIVE_MAIN_BG}
      name={props.icon}
      size={24}
    /> : null;
  return (
    <TouchableOpacity style={[styles.loginButtonStyle, props.buttonStyle]} onPress={props.onPress}>
      { icon }
      <Text style={[styles.loginButtonTextStyle, props.textStyle]}>
        { props.title }
      </Text>
    </TouchableOpacity>
  );
};

export default Button;

const styles = StyleSheet.create({
  loginButtonStyle: {
    flexDirection: 'row',
    borderColor: Colors.HIVE_MAIN_BG,
    borderWidth: 0.5,
    borderRadius: 5,
    height: 28,
    backgroundColor: 'white',
    alignItems: 'center',
    justifyContent: 'center',
  },
  loginButtonTextStyle: {
    fontSize: 14,
    color: Colors.HIVE_MAIN_BG,
  },
});
