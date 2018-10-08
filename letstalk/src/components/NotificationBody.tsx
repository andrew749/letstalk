import React from 'react';
import {
  TouchableHighlight,
  StatusBar,
  View,
  Text,
  Image,
  Vibration,
  Platform,
  StyleSheet,
} from 'react-native';
import GestureRecognizer, { swipeDirections } from 'react-native-swipe-gestures';
import NotificationService from '../services/notification-service';
import { Linking } from 'expo';

const styles = StyleSheet.create({
  root: {
    flex: 1,
    backgroundColor: 'white',
  },
  container: {
    position: 'absolute',
    top: 0,
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: 'white',
  },
  content: {
    flex: 1,
    flexDirection: 'row',
    backgroundColor: 'white',
    borderRadius: 20,
    margin: 10,
  },
  iconApp: {
    marginTop: 10,
    marginLeft: 20,
    resizeMode: 'contain',
    width: 24,
    height: 24,
    borderRadius: 5,
  },
  icon: {
    marginTop: 10,
    marginLeft: 10,
    resizeMode: 'contain',
    width: 48,
    height: 48,
  },
  textContainer: {
    alignSelf: 'center',
    marginLeft: 10,
  },
  title: {
    color: 'black',
    fontWeight: 'bold',
  },
  message: {
    color: 'gray',
    marginTop: 5,
  },
  footer: {
    backgroundColor: '#696969',
    borderRadius: 5,
    alignSelf: 'center',
    height: 5,
    width: 35,
    margin: 5,
  },
});

// TODO: We don't do anything with `icon` yet
interface Props {
  title: string;
  message: string;
  onPress: () => void;
  onClose: () => void;
  icon: string;
  vibrate: boolean;
  isOpen: boolean;
}

class NotificationBody extends React.Component<Props> {
  constructor(props: Props) {
    super(props);

    this.onSwipe = this.onSwipe.bind(this);
  }

  componentWillReceiveProps(nextProps: Props) {
    if (Platform.OS === 'ios' && nextProps.isOpen !== this.props.isOpen) {
      StatusBar.setHidden(nextProps.isOpen);
    }

    if ((this.props.vibrate || nextProps.vibrate) && nextProps.isOpen && !this.props.isOpen) {
      Vibration.vibrate(400, false);
    }
  }

  onSwipe(direction: any) {
    const { onClose } = this.props;

    if (Platform.OS === 'ios') {
      const { SWIPE_UP } = swipeDirections;
      if (direction === SWIPE_UP) this.props.onClose();
    } else {
      const { SWIPE_LEFT, SWIPE_RIGHT } = swipeDirections;
      if (direction === SWIPE_RIGHT || direction === SWIPE_LEFT) onClose();
    }
  }

  render() {
    const {
      title,
      message,
      onPress,
      onClose,
    } = this.props;

    const footer = Platform.OS === 'ios' ? <View style={styles.footer} /> : null;

    return (
      <View style={styles.root}>
        <GestureRecognizer onSwipe={this.onSwipe} style={styles.container}>
          <TouchableHighlight
            style={styles.content}
            activeOpacity={0.3}
            underlayColor="transparent"
            onPress={() => {
              onClose();
              onPress();
            }}
          >
            <View style={styles.textContainer}>
              <Text numberOfLines={1} style={styles.title}>{title}</Text>
              <Text numberOfLines={1} style={styles.message}>{message}</Text>
            </View>
          </TouchableHighlight>
          { footer }
        </GestureRecognizer>
      </View>
    );
  }
}

export default NotificationBody;
