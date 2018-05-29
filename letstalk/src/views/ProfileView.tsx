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
import { Icon } from 'react-native-elements';
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
import {fbLogin} from '../services/fb';
import { ActionButton, Card, Header, Loading } from '../components';
import { genderIdToString } from '../models/user';
import { RootState } from '../redux';
import { State as ProfileState, fetchProfile } from '../redux/profile/reducer';
import { ActionTypes } from '../redux/profile/actions';
import { programById, sequenceById } from '../models/cohort';
import { AnalyticsHelper } from '../services/analytics';
import { ProfileAvatar } from '../components';
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

  constructor(props: Props) {
    super(props);

    this.onLogoutPress = this.onLogoutPress.bind(this);
    this.load = this.load.bind(this);
    this.renderInner = this.renderInner.bind(this);
    // this.openQrScannerView = this.openQrScannerView.bind(this);
  }

  private async onLogoutPress() {
    try {
      await auth.logout();
    } catch (error) { }
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

  private renderProfile(gradYear: string, program: string, bio: string | null) {
    const bioStr = bio === null ? 'Add bio by editing profile' : bio;
    return (
      <View style={styles.sectionContainer}>
        <Text style={styles.description}>{ bioStr }</Text>
        <Text style={styles.profileTitle}>{program}, {gradYear}</Text>
      </View>
    )
  }

  private renderContactInfo(email: string, fbId: string, phoneNumber: string) {

    const buildItem = (label: string, value: string) => {
      return (
        <View key={label} style={styles.listItem}>
          <Text style={styles.label}>{label}: </Text>
          <Text style={styles.value}>{value}</Text>
        </View>
      );
    };

    const buildItems = (name_values: Array<[string, string]>) => {
      const items = name_values.map(([label, value]) => {
        return buildItem(label, value);
      });
      return items;
    };

    const contactItems = buildItems([
      ['Phone', phoneNumber],
      ['Email', email]
    ]);


    if (fbId !== null) {
      const fbLink = 'fb://profile/' + fbId;
      contactItems.push(
        <TouchableOpacity style={styles.listItem} onPress={() => Linking.openURL(fbLink)}>
          <MaterialIcons name="face" size={24} />
          <Text style={styles.label}>Facebook: </Text>
          <Text style={styles.value}>{fbId}</Text>
        </TouchableOpacity>
      );
    } else {
      // link fb profile
      contactItems.push(
        <TouchableOpacity
          style={styles.listItem}
          onPress={async () => {
            await auth.linkFB();
            await this.props.fetchProfile();
          }}
        >
          <Text style={styles.value}>Link your Facebook account</Text>
        </TouchableOpacity>
      );
    }

    return (
      <View style={styles.sectionContainer}>
        <Text style={styles.sectionHeader}>Contact Info</Text>
        {contactItems}
      </View>
    )
  }

  private renderInner() {

    const { navigate } = this.props.navigation;

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
      bio,
      hometown,
    } = this.props.profile;

    const capitalize = (s: string) => s.charAt(0).toUpperCase() + s.slice(1);

    let userId;
    if (this.props.profile) {
      userId = this.props.profile.userId.toString();
    }
    const headerText = this.props.profile ?
      this.props.profile.firstName + ' ' + this.props.profile.lastName : 'Profile';

    const genderStr = capitalize(genderIdToString(gender));
    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    const timeDiff = new Date().valueOf() - birthdate.valueOf();
    const age = Math.floor(timeDiff / (1000 * 60 * 60 * 24 * 365));

    const sequence = sequenceById(sequenceId);
    const program = programById(programId);

    const hometownStr = hometown === null || hometown === '' ? 'Some place on Earth' : hometown;

    return (
      <View style={styles.contentContainer} >
        {this.renderQrCode()}
        <ProfileAvatar userId={userId} xlarge containerStyle={styles.profilePicture} />
        <Header>{headerText}</Header>
        <Icon
          name='pencil'
          type='font-awesome'
          color={Colors.HIVE_PRIMARY}
          containerStyle={styles.editButton}
          onPress={() => navigate('ProfileEdit')} />
        {/* <ReactNativeButton title="Scan QR Code" onPress={this.openQrScannerView} /> */}
        <Text style={styles.subHeaderText}>{age}{genderStr[0]} - {hometownStr}</Text>
        {this.renderProfile(String(gradYear), program, bio)}
        {this.renderContactInfo(email, fbId, phoneNumber)}
        <View style={styles.sectionContainer}>
        </View>
        <ActionButton buttonStyle={styles.logoutButton} onPress={this.onLogoutPress} title='LOGOUT' />
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
        size={150}
        bgColor='black'
        fgColor='white'
      />
    );
  };

  render() {
    const body = this.renderBody();
    return (
      <ScrollView contentContainerStyle={styles.container}>
        {body}
      </ScrollView>
    );
  }

  // private async openQrScannerView() {
  //   this.props.navigation.navigate('QrScanner');
  // }
}

export default connect(({ profile }: RootState) => profile, { fetchProfile })(ProfileView);

const styles = StyleSheet.create({
  container: {
    paddingBottom: 10,
    minHeight: '100%'
  },
  contentContainer: {
    alignItems: 'center',
    flex: 1,
    flexDirection: 'column',
    margin: 20,
    padding: 20,
    borderWidth: 2,
    borderColor: Colors.HIVE_PRIMARY,
    backgroundColor: Colors.WHITE
  },
  description: {
    fontSize: 18,
    color: Colors.HIVE_SUBDUED
  },
  editButton: {
    position: 'absolute',
    right: 0,
    margin: 20
  },
  listItem: {
    flex: 1,
    flexDirection: 'row',
    marginTop: 8,
    marginBottom: 8,
  },
  logoutButton: {
    color: Colors.HIVE_ERROR
  },
  profileTitle: {
    fontSize: 18,
    marginTop: 10,
    alignSelf: 'flex-end'
  },
  profilePicture: {
    borderWidth: 2,
    borderColor: Colors.HIVE_PRIMARY,
    margin: 20
  },
  sectionHeader: {
    fontWeight: 'bold',
    fontSize: 24,
    alignSelf: 'flex-start',
  },
  sectionContainer: {
    width: "100%",
    backgroundColor: 'white',
    flex: 1,
    flexDirection: 'column',
    marginTop: 20
  },
  subHeaderText: {
    fontSize: 18
  },
  label: {
    fontSize: 18,
  },
  value: {
    fontSize: 18,
    color: Colors.HIVE_ACCENT
  },
});
