import React, { Component } from 'react';
import {
  ActivityIndicator,
  Image,
  ImageBackground,
  Platform,
  StyleSheet,
  View,
} from 'react-native';
import {
  NavigationActions,
  NavigationScreenProp,
  NavigationStackAction,
} from 'react-navigation';
import { connect, ActionCreator, Dispatch } from 'react-redux';
import { ThunkAction } from 'redux-thunk';

import { RootState } from '../redux';
import { ActionTypes as BootstrapActionTypes } from '../redux/bootstrap/actions';
import {
  State as BootstrapState,
  fetchBootstrap,
} from '../redux/bootstrap/reducer';
import { errorToast, infoToast } from '../redux/toast';
import { AnalyticsHelper } from '../services/analytics';
import profileService from '../services/profile-service';
import auth from '../services/auth';
import Color from '../services/colors';
import Window from '../services/window';

interface DispatchActions {
  errorToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  infoToast(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  fetchBootstrap: ActionCreator<ThunkAction<Promise<BootstrapActionTypes>, BootstrapState, void>>;
}

interface Props extends DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

// Since this view gets rendered multiple times, we only want to add the device token once per app
// launch.
let addedExpoToken = false;

const PROBABILITY_SHOW_IOS_MODEL = 0.4;

function showModal(): boolean {
  return Math.random() < PROBABILITY_SHOW_IOS_MODEL;
}

class SplashScreenView extends Component<Props> {
  SPLASH_SCREEN_VIEW_IDENTIFIER = "SplashScreenView";

  static navigationOptions = () => ({
    header: null as any,
  });

  private triedAddingToken: boolean = false;

  private async maybeAddExpoToken() {
    if (!addedExpoToken) {
      addedExpoToken = true;
      try {
        const token = await auth.registerForPushNotificationsAsync(showModal());
        if (token !== null) {
          await profileService.addExpoDeviceToken(token);
        }
      } catch(e){
        console.log("Failed to register for notification " + e);
      }
    }
  }

  private async load() {
    const sessionToken = await auth.getSessionToken();
    if (sessionToken !== null) {
      try {
        await Promise.all([
          this.props.fetchBootstrap(),
          this.maybeAddExpoToken(),
        ]);
        await this.props.navigation.dispatch(NavigationActions.reset({
          index: 0,
          actions: [NavigationActions.navigate({ routeName: 'Tabbed' })]
        }));
      } catch (e) {
        await this.props.errorToast(e.errorMsg);
      }
    }
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.SPLASH_SCREEN_VIEW_IDENTIFIER);
    });
    await this.load();
  }

  render() {
    return (
      <View style={styles.container}>
        <ImageBackground
          style={styles.logo}
          source={require('../img/logo_android.png')}
        >
          <ActivityIndicator
            color={Color.HIVE_PRIMARY}
            size="large"
            style={styles.spinner}
          />
        </ImageBackground>
        <Image
          style={styles.name}
          source={require('../img/name_white.png')}
        />
      </View>
    );
  }
};

export default connect(null,
  { errorToast, infoToast, fetchBootstrap })(SplashScreenView);

const styles = StyleSheet.create({
  container: {
    width: Window.WIDTH,
    height: Window.HEIGHT,
    backgroundColor: Color.HIVE_PRIMARY,
    justifyContent: 'center',
    alignItems: 'center',
  },
  logo: {
    width: 200,
    height: 200,
  },
  name: {
    marginTop: 20,
    width: 101,
    height: 57,
  },
  spinner: {
    position: 'absolute',
    // Bit of a hack that makes it easier to position the spinner
    width: 0,
    height: 0,
    left: Platform.OS === 'ios' ? 107 : 105,
    top: Platform.OS === 'ios' ? 123 : 121,
  },
});