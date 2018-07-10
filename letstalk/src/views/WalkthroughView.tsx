import React, {Component} from 'react';
import {View, Text} from 'react-native';
import { NavigationScreenProp, NavigationStackAction } from 'react-navigation';
import Colors from '../services/colors';
import { AsyncStorage } from 'react-native';
import TutorialService from '../services/tutorial_service';

// @ts-ignore
import AppIntroSlider from 'react-native-app-intro-slider';

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class WalkthroughView extends Component<Props> {
  constructor(props: Props) {
    super(props);
  }

  onDone = async () => {
    this.props.navigation.pop();
    await TutorialService.ackShowTutorial();
  }

  render() {
    const background = Colors.HIVE_PRIMARY;
    const pageArray = [
      {
        title: 'Welcome to Hive!',
        text: '',
        image: require('../img/logo_android.png'),
        imageStyle: {
          height: 360,
          width: 360,
        },
        backgroundColor: background,
        fontColor: '#fff',
        level: 10,
      },
      {
        title: 'Meet other cohorts',
        text: 'Connect with upper and lower years in your program.',
        image: require('../assets/images/mentor.jpg'),
        imageStyle: {
          height: 200,
          width: 360,
        },
        backgroundColor: background,
        fontColor: '#fff',
        level: 10,
      },
    {
        title: 'Find new connections',
        text: 'Meet awesome people in your network.',
        image: require('../assets/images/rtm.jpg'),
        imageStyle: {
          height: 400,
          width: 360,
        },
        backgroundColor: background,
        fontColor: '#fff',
        level: 10,
      },
      {
        title: 'Built with privacy in mind',
        text: 'Nobody can view your information unless you both connect',
        image: require('../assets/images/lock.png'),
        imageStyle: {
          height: 400,
          width: 355,
        },
        backgroundColor: background,
        fontColor: '#fff',
        level: 10,
      }
    ];
    return(
      <AppIntroSlider
        slides={pageArray}
        onDone={this.onDone}
      />
    );
  }
}

export default WalkthroughView;
