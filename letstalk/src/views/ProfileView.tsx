import React, { Component } from 'react';
import {
  ActivityIndicator,
  AppRegistry,
  Button as ReactNativeButton,
  Dimensions,
  FlatList,
  Image,
  Linking,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import { bindActionCreators } from 'redux'
import {
  NavigationScreenDetails,
  NavigationScreenProp,
  NavigationStackAction,
  NavigationActions,
} from 'react-navigation';
import { MaterialIcons } from '@expo/vector-icons';
import Immutable from 'immutable';

import auth from '../services/auth';
import { ActionButton, Card, Header, Loading } from '../components';
import { genderIdToString } from '../models/user';
import { RootState } from '../redux';
import { State as ProfileState, fetchProfile } from '../redux/profile/reducer';
import { ActionTypes } from '../redux/profile/actions';
import { programById, sequenceById } from '../models/cohort';
import {AnalyticsHelper} from '../services/analytics';
import {ProfileAvatar} from '../components';
import Colors from '../services/colors';
import QRCode from "react-native-qrcode";

const SCREEN_WIDTH = Dimensions.get('window').width;

interface DispatchActions {
  fetchProfile: ActionCreator<ThunkAction<Promise<ActionTypes>, ProfileState, void>>;
}

interface Props extends ProfileState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class ProfileView extends Component<Props> {
  PROFILE_VIEW_IDENTIFIER = "ProfileView";

  static navigationOptions = ({ navigation }: NavigationScreenDetails<void>) => ({
    headerTitle: 'Profile',
    headerRight: <ReactNativeButton title="Edit"
      onPress={() => navigation.navigate('ProfileEdit')} />,
    headerStyle: {
      backgroundColor: Colors.HIVE_MAIN_BG,
    },
  })

  constructor(props: Props) {
    super(props);

    this.onLogoutPress = this.onLogoutPress.bind(this);
    this.load = this.load.bind(this);
    this.renderInner = this.renderInner.bind(this);
    this.openQrScannerView = this.openQrScannerView.bind(this);
  }

  private async onLogoutPress() {
    try {
      await auth.logout();
    } catch(error) {}
    this.props.navigation.dispatch(NavigationActions.reset({
      index: 0,
      key: null,
      actions: [NavigationActions.navigate({ routeName: 'Login' })]
    }));
  }

  async componentDidMount() {
    AnalyticsHelper.getInstance().recordPage(this.PROFILE_VIEW_IDENTIFIER);
    this.load();
  }

  private async load() {
    await this.props.fetchProfile();
  }

  private renderInner() {
    const {
      programId,
      gradYear,
      sequenceId,
    } = this.props.profile;

    const {
      gender,
      email,
      birthdate,
      phoneNumber,
      fbId,
    } = this.props.profile;

    const capitalize = (s: string) => s.charAt(0).toUpperCase() + s.slice(1);

    const genderStr = capitalize(genderIdToString(gender));
    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    const birthdayStr = birthdate.toLocaleDateString('en-US', options);
    const sequence = sequenceById(sequenceId);
    const program = programById(programId);

    const hr = <View
      style={{
        borderBottomColor: '#555555',
        borderBottomWidth: 1,
      }}
    />;

    const buildItem = (label: string, value: string) => {
      return (
        <View key={label} style={styles.listItem}>
          <Text style={styles.label}>{label}:</Text>
          <Text style={styles.value}>{value}</Text>
        </View>
      );
    };

    const buildItems = (name_values: Array<[string, string]>) => {
      const items = name_values.map(([label, value]) => {
        return [buildItem(label, value), hr];
      });
      const flatItems = [].concat(...items);
      return flatItems.length > 0 ? flatItems.slice(0, flatItems.length - 1) : flatItems;
    };

    const profileItems = buildItems([
      ['Gender', genderStr],
      ['Birthday', birthdayStr],
      ['Email', email],
    ]);

    if (phoneNumber !== null) {
      profileItems.push(hr);
      profileItems.push(buildItem('Phone Number', phoneNumber));
    }

    if (fbId !== null) {
      const fbLink = 'fb://profile/' + fbId;
      profileItems.push(hr);
      profileItems.push(
        <TouchableOpacity style={styles.listItem} onPress={() => Linking.openURL(fbLink)}>
          <MaterialIcons name="face" size={24} />
          <Text style={styles.value}>Facebook profile</Text>
        </TouchableOpacity>
      );
    }

    const cohortItems = buildItems([
      ['Program', program],
      ['Sequence', sequence],
      ['Grad year', String(gradYear)],
    ]);

    return (
      <View style={styles.contentContainer} >
        <Image style={styles.image} source={require('../img/profile.jpg')} />
        {this.renderQrCode()}
        <ReactNativeButton title="Scan QR Code"
                           onPress={this.openQrScannerView} />
        <Text style={styles.sectionHeader}>Personal Info</Text>
        <View style={styles.sectionContainer}>
          {profileItems}
        </View>
        <Text style={styles.sectionHeader}>Cohort</Text>
        <View style={styles.sectionContainer}>
          {cohortItems}
        </View>
      </View>
    );
  }

  renderBody() {
    const {
      state,
      errorMsg,
    } = this.props.fetchState;
    return (
      <Loading
        state={state}
        errorMsg={errorMsg}
        load={this.load}
        renderBody={this.renderInner}
      />
    );
  }

  renderQrCode = () => {
    const {secret} = this.props.profile;
    return (
      !!secret && <QRCode
        value={secret}
        size={200}
        bgColor='black'
        fgColor='white'
      />
    );
  };

  render() {
    const body = this.renderBody();

    let userId;
    if (this.props.profile) {
        userId = this.props.profile.userId.toString();
    }
    const headerText = this.props.profile ?
      this.props.profile.firstName + ' ' + this.props.profile.lastName : 'Profile';
    return (
      <ScrollView contentContainerStyle={styles.container}>
        <View style={styles.headerContainer}>
          <ProfileAvatar userId={userId} />
          <Header>{headerText}</Header>
        </View>
        {body}
        <ActionButton onPress={this.onLogoutPress} title='LOGOUT'/>
      </ScrollView>
    );
  }

  private async openQrScannerView() {
    this.props.navigation.navigate('QrScanner');
  }
}

export default connect(({profile}: RootState) => profile, { fetchProfile })(ProfileView);

const styles = StyleSheet.create({
  container: {
    paddingBottom: 10,
  },
  contentContainer: {
    alignItems: 'center',
    marginHorizontal: 25
  },
  image: {
    width: 60,
    height: 60,
    margin: 10,
    borderRadius: 30,
  },
  headerContainer: {
    flex: 1,
    flexDirection: 'row',
  },
  listItem: {
    flex: 1,
    flexDirection: 'row',
    marginTop: 8,
    marginBottom: 8,
  },
  sectionHeader: {
    fontWeight: 'bold',
    fontSize: 24,
    marginBottom: 5,
    marginTop: 10,
    alignSelf: 'flex-start',
  },
  sectionContainer: {
    width: SCREEN_WIDTH,
    padding: 10,
    backgroundColor: 'white',
  },
  label: {
    fontWeight: 'bold',
    fontSize: 18,
  },
  value: {
    fontSize: 18,
    marginLeft: 10,
  },
});
