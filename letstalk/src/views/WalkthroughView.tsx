import React, {Component} from 'react';
import {View, Text} from 'react-native';
import AppIntroSlider from 'react-native-app-intro-slider';
import { NavigationScreenProp, NavigationStackAction } from 'react-navigation';
import Colors from '../services/colors';
import { AsyncStorage } from 'react-native';
import TutorialService from '../services/tutorial_service';

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
        title: 'Explore your program',
        text: 'Connect with upper years mentors in your program.',
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
        text: 'Meet people you want to meet. Search for people based on interests, experiences, etc.',
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
        text: 'Nobody can view your information until you are matched with them',
        image: require('../img/profile.jpg'),
        imageStyle: {
          height: 360,
          width: 360,
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
