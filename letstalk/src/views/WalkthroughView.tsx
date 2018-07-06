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
    const pageArray = [
      {
        title: 'Welcome to Hive!',
        text: 'Connect with interesting new people around you',
        image: require('../img/logo_android.png'),
        imageStyle: {
          height: 103 * 2.5,
          width: 103 * 2.5,
        },
        backgroundColor: Colors.HIVE_PRIMARY,
        fontColor: '#fff',
        level: 10,
      },
      {
        title: 'View matches',
        text: 'Connect with upper years in your program',
        image: require('../img/profile.jpg'),
        imageStyle: {
          height: 103 * 2.5,
          width: 103 * 2.5,
        },
        backgroundColor: '#a4b602',
        fontColor: '#fff',
        level: 10,
      },
    {
        title: 'Find new connections',
        text: 'Meet people you want to meet. Search for people based on interests, experiences, etc.',
        image: require('../img/profile.jpg'),
        imageStyle: {
          height: 103 * 2.5,
          width: 103 * 2.5,
        },
        backgroundColor: '#a4b602',
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
