import React, { Component } from 'react';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import {
  ActivityIndicator,
  Button as ReactNativeButton,
  Linking,
  ScrollView,
  StyleSheet,
  Text,
  TouchableOpacity,
  TouchableWithoutFeedback,
  View,
} from 'react-native';
import {
  NavigationScreenProp,
  NavigationScreenDetails,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import { MaterialIcons } from '@expo/vector-icons';
import Menu, { MenuOptions, MenuOption, MenuTrigger } from 'react-native-menu';

import { RootState } from '../redux';
import { State as BootstrapState, fetchBootstrap } from '../redux/bootstrap/reducer';
import { ActionTypes } from '../redux/bootstrap/actions';
import { ActionButton, Card, Header, Loading } from '../components';

import Colors from '../services/colors';
import {AnalyticsHelper} from '../services/analytics';
import navigate = NavigationActions.navigate;

interface DispatchActions {
  fetchBootstrap: ActionCreator<ThunkAction<Promise<ActionTypes>, BootstrapState, void>>;
}

interface Props extends BootstrapState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class HomeView extends Component<Props> {
  HOME_VIEW_IDENTIFIER = "HomeView";

  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: 'Home',
    headerRight: <ReactNativeButton title="Profile"
      onPress={() => navigation.navigate('Profile')} />,
  })

  constructor(props: Props) {
    super(props);

    this.load = this.load.bind(this);
    this.renderHome = this.renderHome.bind(this);
    // TODO(aklen): remove this.menu when QR code moves to a different page
    this.menu = {};
  }

  async componentDidMount() {
    AnalyticsHelper.getInstance().recordPage(this.HOME_VIEW_IDENTIFIER);
    this.load();
  }

  async componentWillReceiveProps(nextProps: Props) {
    if (nextProps.bootstrap && nextProps.bootstrap.state === 'account_created') {
      // Account not yet setup, so take to onboarding page
      this.props.navigation.dispatch(NavigationActions.reset({
        index: 0,
        actions: [NavigationActions.navigate({ routeName: 'Onboarding' })]
      }));
    }
  }

  private async load() {
    await this.props.fetchBootstrap();
  }

  private requestToMatch() {
    this.props.navigation.navigate('RequestToMatch');
  }

  private credentialEdit() {
    this.props.navigation.navigate('CredentialEdit');
  }

  private static openQrCodePage(userId: number) {
    navigate('QrCode', { userId: userId });
  }

  private renderMatches() {
    const { relationships } = this.props.bootstrap;
    const matchCards = relationships.map(relationship => {
      const {
        userId,
        userType,
        firstName,
        lastName,
        email,
        fbId,
        phoneNumber,
      } = relationship;

      const fbLink = 'fb://profile/' + fbId;
      const emailLink = 'mailto:' + email;
      const smsLink = 'sms:' + phoneNumber;
      const fb = fbId === null ? null : (
        <TouchableOpacity style={styles.emailContainer} onPress={() => Linking.openURL(fbLink)}>
          <MaterialIcons name="face" size={24} />
          <Text style={styles.email}>Talk on Messenger</Text>
        </TouchableOpacity>
      );
      const sms = phoneNumber === null ? null : (
        <TouchableOpacity style={styles.emailContainer} onPress={() => Linking.openURL(smsLink)}>
          <MaterialIcons name="textsms" size={24} />
          <Text style={styles.email}>Talk on SMS</Text>
        </TouchableOpacity>
      );
      // TODO: Handle errors for links
      return (
        // TODO(ajklen): remove this TouchableWithoutFeedback and Menu when QR code moves to a different page
        <TouchableWithoutFeedback onLongPress={()=>this.menu[userId].open()}>
          <Card key={userId}>
            <Text style={styles.name}>{firstName + ' ' + lastName}</Text>
            <Text style={styles.relationshipType}>{userType}</Text>
            <TouchableOpacity style={styles.emailContainer}
              onPress={() => Linking.openURL(emailLink)}
            >
              <MaterialIcons name="email" size={24} />
              <Text style={styles.email}>{email}</Text>
            </TouchableOpacity>
            {fb}
            {sms}
          </Card>
          <Menu ref={c => (this.menu[userId] = c)} onSelect={callback => callback(userId)}>
            <MenuTrigger>
              <Text style={{ fontSize: 20 }}>&#8942;</Text>
            </MenuTrigger>
            <MenuOptions value={HomeView.openQrCodePage}>
              <Text>Open QR Code</Text>
            </MenuOptions>
          </Menu>
        </TouchableWithoutFeedback>
      );
    });
    return (
      <View>
        <Header>Matches</Header>
        {matchCards}
      </View>
    );
  }

  private renderHome() {
    const { state } = this.props.bootstrap;
    switch (state) {
      case 'account_created':
        // Should in reality not be shown, since we never show home page for account_created
        return (
          <View style={styles.centeredContainer}>
            <Text style={styles.headline}>Waiting for you to finish onboarding</Text>
            <ActionButton onPress={() => this.load()} title="Check again" />
          </View>
        );
      case 'account_setup':
        return (
          <View style={styles.centeredContainer}>
            <Text style={styles.headline}>Waiting for your match</Text>
            <ActionButton onPress={() => this.load()} title="Check again" />
          </View>
        );
      case 'account_matched':
        const matches = this.renderMatches();
        return (
          <ScrollView>
            { matches }
            <ActionButton onPress={() => this.requestToMatch()} title="Request To Match" />
            <ActionButton onPress={() => this.credentialEdit()} title="Edit Credentials" />
          </ScrollView>
        );
      default:
        // Ensure exhaustiveness of select
        const _: never = state;
    }
  }

  render() {
    const {
      state,
      errorMsg,
    } = this.props.fetchState;
    return (
      <Loading
        state={state}
        errorMsg={errorMsg}
        load={this.load}
        renderBody={this.renderHome}
      />
    );
  }
}

export default connect(({ bootstrap }: RootState) => bootstrap, { fetchBootstrap })(HomeView);

const styles = StyleSheet.create({
  centeredContainer: {
    flex: 1,
    justifyContent: 'center'
  },
  headline: {
    fontWeight: 'bold',
    fontSize: 18,
    textAlign: 'center',
    marginBottom: 10,
  },
  name: {
    fontWeight: 'bold',
    fontSize: 18,
  },
  relationshipType: {
    fontSize: 12,
    color: 'gray',
    marginBottom: 5,
  },
  emailContainer: {
    flex: 1,
    flexDirection: 'row',
  },
  email: {
    paddingTop: 2,
    marginLeft: 5,
    fontSize: 16,
  },
})
