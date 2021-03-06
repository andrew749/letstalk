import React, { Component } from 'react';
import {
  Button,
  Keyboard,
  EmitterSubscription,
  StyleSheet,
  Platform,
  Text,
  View,
} from 'react-native';
import { Provider } from 'react-redux';
import { combineReducers, compose, createStore, applyMiddleware } from 'redux';
import { MaterialIcons } from '@expo/vector-icons';
import { Notifications, Font, Linking } from 'expo';
import createLogger from 'redux-logger';
import thunk from 'redux-thunk';
import {
  NavigationContainerComponent,
  StackNavigator,
  TabNavigator,
  TabBarBottomProps,
  TabBarBottom,
} from 'react-navigation';
import NotificationComponent from 'react-native-in-app-notification';
import Sentry from 'sentry-expo';
import { Toast } from 'react-native-redux-toast';
import { YellowBox } from 'react-native'

import appReducer from './redux';
import auth from './services/auth';
import HomeView from './views/HomeView';
import LoginView from './views/LoginView';
import MatchProfileView from './views/MatchProfileView';
import ProfileView from './views/ProfileView';
import ProfileEditView from './views/ProfileEditView';
import SignupView from './views/SignupView';
import OnboardingView from './views/OnboardingView';
import ExploreView from './views/ExploreView';
import NotificationView from './views/NotificationView';
import ForgotPasswordView from './views/ForgotPasswordView';
import WalkthroughView from './views/WalkthroughView';
import QrScannerView from "./views/QrScannerView";
import NotificationContentView from './views/NotificationContentView';
import EditProfileSelectorView from './views/EditProfileSelectorView';
import ChangeCohortView from './views/trait-editors/ChangeCohortView';
import AddSimpleTraitView from './views/trait-editors/AddSimpleTraitView';
import AddPositionView from './views/trait-editors/AddPositionView';
import AddGroupView from './views/trait-editors/AddGroupView';
import QrCodeView from './views/QrCodeView';
import VerifyEmailView from './views/VerifyEmailView';
import RequestConnectionView from './views/RequestConnectionView';
import RequestsView from './views/RequestsView';
import UpdatePersonalView from './views/UpdatePersonalView';
import SurveyView from './views/SurveyView';

import NotificationService, { RawNotification } from './services/notification-service';
import navService from './services/navigation-service';
import Colors from './services/colors';
import { NotificationBody } from './components';
import { AsyncStorage } from 'react-native';

YellowBox.ignoreWarnings(['Warning: isMounted(...) is deprecated']);

Sentry.config('https://a851f386142f48c1956ca9f7e06e8a4a:d43e2ccf892644508f4a18e22d8e0df4@sentry.io/305511').install();
Sentry.setTagsContext({environment: (__DEV__) ? "development" : "production" });

// turn off annoying warnings
console.disableYellowBox = true;

interface TabBarIcon {
  tintColor: (string | null),
  focused: boolean,
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: Colors.HIVE_PRIMARY,
    alignItems: 'center',
    justifyContent: 'center'
  },
  toastMessageStyle: {
    color: 'white',
  },
  toastContainerStyle: {
    backgroundColor: Colors.HIVE_PRIMARY,
  },
  toastErrorStyle: {
    backgroundColor: Colors.HIVE_ERROR,
  },
});

interface TabBarState {
  visible: boolean;
}

const prefix =  Platform.OS == 'android' ? 'hive://hive/' : 'hive://';

class TabBar extends Component<TabBarBottomProps, TabBarState> {
  private keyboardEventListeners: Array<EmitterSubscription>;

  constructor(props: TabBarBottomProps) {
    super(props);

    this.state = { visible: true };

    this.visible = this.visible.bind(this);
  }

  componentDidMount() {
    if (Platform.OS !== 'ios') {
      this.keyboardEventListeners = [
        Keyboard.addListener('keyboardDidShow', this.visible(false)),
        Keyboard.addListener('keyboardDidHide', this.visible(true))
      ];
    } else {
      this.keyboardEventListeners = [];
    }
  }

  componentWillUnmount() {
    this.keyboardEventListeners.forEach((eventListener) => eventListener.remove());
  }

  visible = (visible: boolean) => () => this.setState({ visible });

  render() {
    if (!this.state.visible) return null;
    return <TabBarBottom {...this.props}/>;
  }
}

const createTabView = () => TabNavigator({
  'Home': {
    screen: HomeView,
  },
  'Explore': {
    screen: ExploreView,
  },
  'Profile': {
    screen: ProfileView,
  },
}, {
  tabBarPosition: 'bottom',
  navigationOptions: ({ navigation }) => ({
    tabBarIcon: ({ focused, tintColor }) => {
      const { routeName } = navigation.state;
      let iconName;
      if (routeName === 'Home') {
        iconName = 'home';
      } else if (routeName === 'Profile') {
        iconName = 'person';
      } else if (routeName === 'Explore') {
        iconName = 'search';
      }

      return <MaterialIcons name={iconName} size={24} color={tintColor} />;
    },
  }),
  tabBarOptions: {
    showLabel: Platform.OS === 'ios',
    showIcon: true,
    activeTintColor: Colors.HIVE_PRIMARY,
    inactiveTintColor: 'gray',
    style: {
      backgroundColor: 'white',
    },
  },
  tabBarComponent: TabBar,
});

const createAppNavigation = (initialRouteName: string) => StackNavigator({
  AddPosition: {
    screen: AddPositionView,
    path: 'AddPosition',
  },
  AddGroup: {
    screen: AddGroupView,
    path: 'AddGroup',
  },
  AddSimpleTrait: {
    screen: AddSimpleTraitView,
    path: 'AddSimpleTrait',
  },
  BlankDoNotUse: {
    screen: View
  },
  ChangeCohort: {
    screen: ChangeCohortView
  },
  EditProfileSelector: {
    screen: EditProfileSelectorView
  },
  ForgotPassword: {
    screen: ForgotPasswordView,
    path: 'ForgotPassword'
  },
  Login: {
    screen: LoginView,
    path: 'Login'
  },
  MatchProfile: {
    screen: MatchProfileView,
    path: 'MatchProfile/:userId'
  },
  NotificationContent: {
    screen: NotificationContentView,
    path: 'NotificationContent/:notificationId'
  },
  NotificationView: {
    screen: NotificationView,
    path: 'NotificationView'
  },
  Onboarding: {
    screen: OnboardingView
  },
  ProfileEdit: {
    screen: ProfileEditView // TODO: needs work on the profile edit view to work properly.
    // need to add ability to load without props.
    // path: 'ProfileEdit',

  },
  QrCode: {
    screen: QrCodeView,
    path: 'QrCode'
  },
  QrScanner: {
    screen: QrScannerView,
    path: 'QrScanner'
  },
  RequestConnection: {
    screen: RequestConnectionView
  },
  Requests: {
    screen: RequestsView
  },
  Signup: {
    screen: SignupView,
    path: 'Signup'
  },
  Tabbed: {
    screen: createTabView(),
    navigationOptions: {}
  },
  UpdatePersonal: {
    screen: UpdatePersonalView,
    path: 'UpdatePersonal',
  },
  VerifyEmail: {
    screen: VerifyEmailView
  },
  WalkthroughView: {
    screen: WalkthroughView,
    navigationOptions: {
      header: null
    }
  },
  SurveyView: {
    screen: SurveyView,
  },
}, {
  initialRouteName,
}, );

const store = createStore(appReducer, applyMiddleware(thunk));

interface AppState {
  loggedIn: null | boolean;
}

type Props = {};

class App extends React.Component<Props, AppState> {
  private notificationService: NotificationService;

  constructor(props: Props) {
    super(props);
    this.state = {
      loggedIn: null,
    };

    this.handleNotification = this.handleNotification.bind(this);
    this.notificationService = new NotificationService(store);
    Linking.addEventListener('url', this.linkHandler);
  }

  async handleNotification(notification: any) {
    await this.notificationService.handleNotification(notification as RawNotification);
  }

  private linkHandler = (event: {url: string}) => {
    let { path, queryParams } = Linking.parse(event.url);
    navService.navigate(path, queryParams);
  }

  async componentWillMount() {
    const sessionToken = await auth.getSessionToken();
    this.setState({ loggedIn: sessionToken !== null });
    Notifications.addListener(this.handleNotification);
  }

  render() {
    const { loggedIn } = this.state;

    let initialRouteName = 'BlankDoNotUse';
    if (loggedIn === false) initialRouteName = 'Login';
    else if (loggedIn === true) initialRouteName = 'Tabbed';

    const AppNavigation = createAppNavigation(initialRouteName);

    const addNavContainer = (navContainer: NavigationContainerComponent) => {
      navService.setTopLevelNavigator(navContainer);
      // TODO: refactor to use navigation service
      this.notificationService.setNavContainer(navContainer);
    }
    return (
      <Provider store={store}>
        <View style={{ flex: 1, backgroundColor: Colors.HIVE_BG }}>
          <AppNavigation uriPrefix={prefix} ref={addNavContainer} />
          <NotificationComponent
            ref={(ref: any) => this.notificationService.setNotifContainer(ref)}
            notificationBodyComponent={ NotificationBody }
          />
          <Toast
            containerStyle={styles.toastContainerStyle}
            messageStyle={styles.toastMessageStyle}
            errorStyle={styles.toastErrorStyle} />
        </View>
      </Provider>
    );
  }
}

export default App;
