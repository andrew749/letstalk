import React, { Component, ReactNode } from 'react';
import { ActionCreator, Dispatch } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import {
  Alert,
  Dimensions,
  FlatList,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  TouchableOpacity,
  View,
} from 'react-native';
import {
  NavigationScreenDetails,
  NavigationScreenProp,
  NavigationStackAction,
  NavigationActions,
} from 'react-navigation';
import Immutable from 'immutable';
import { MaterialIcons } from '@expo/vector-icons';
import Moment from 'moment';

import { programById, sequenceById } from '../../models/cohort';
import { Button, Header } from '../../components';
import { UserPosition } from '../../models/position';
import { UserSimpleTrait } from '../../models/simple-trait';
import Colors from '../../services/colors';
import { ActionTypes } from '../../redux/profile/actions';
import { State as ProfileState } from '../../redux/profile/reducer';
import { RootState } from '../../redux';
import { genderIdToString, GenderId } from '../../models/user';

const SCREEN_WIDTH = Dimensions.get('window').width;

interface PersonalInfoProps {
	firstName: string;
	lastName: string;
	gender: GenderId;
	birthdate: string
  secret?: string // only required if `allowQrCode` is true
  bio?: string;
  hometown?: string;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  allowQrCode?: boolean; // default false
}

export class PersonalInfo extends Component<PersonalInfoProps> {
  render() {
    const {
      firstName,
      lastName,
      gender,
      secret,
      birthdate,
      hometown,
      bio,
      allowQrCode
    } = this.props;

    const capitalize = (s: string) => s.charAt(0).toUpperCase() + s.slice(1);
    const genderStr = capitalize(genderIdToString(gender));

    const hometownStr = hometown === null || hometown === '' ? 'Some place on Earth' : hometown;

    const options = { year: 'numeric', month: 'long', day: 'numeric' };
    const timeDiff = new Date().valueOf() - new Date(birthdate).valueOf();
    const age = Math.floor(timeDiff / (1000 * 60 * 60 * 24 * 365));

    const headerText = firstName + ' ' + lastName;

    const bioStr = bio === null ? 'Add bio by editing profile' : bio;

    return (
      <View style={styles.personalInfoContainer}>
        <Header>{headerText}</Header>
        <Text style={styles.subHeaderText}>{age}{genderStr[0]} - {hometownStr}</Text>
        {!!allowQrCode && <TouchableOpacity style={styles.listItem} onPress={() => {
          this.props.navigation.navigate('QrCode', { secret });
        }}>
          <Text style={styles.value}>Show QR Code</Text>
        </TouchableOpacity>}
        <View style={styles.sectionContainer}>
          <Text style={styles.description}>{ bioStr }</Text>
        </View>
      </View>
    );
  }

}

interface CohortInfoProps {
  programId: string;
  sequenceId: string;
  gradYear: number;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  allowEditing?: boolean; // default false
}

export class CohortInfo extends Component<CohortInfoProps> {
  render() {
    const {
      programId,
      gradYear,
      sequenceId,
    } = this.props;
    const program = programById(programId);
    const changeCohort = () => this.props.navigation.navigate('ChangeCohort', {
      programId,
      gradYear,
      sequenceId,
    });
    return (
      <View style={styles.sectionContainer}>
        <TouchableOpacity onPress={changeCohort} style={styles.addTraitButton}>
          <MaterialIcons name="edit" size={32} color={Colors.HIVE_PRIMARY} />
        </TouchableOpacity>
        <Text style={styles.sectionHeader}>Cohort</Text>
        <Text style={styles.cohortText}>{ program + ', ' + gradYear }</Text>
      </View>
    );
  }
}

// Could add props to override these defaults although hasn't been required yet.
const MAX_NUMBER_POSITIONS_SHOWN = 3;
const MAX_NUMBER_SIMPLE_TRAITS_SHOWN = 5;

function renderShowLessMore(isShown: boolean, show: () => void, hide: () => void) {
  if (isShown) {
    return (
      <TouchableOpacity style={styles.listItem} onPress={hide}>
        <Text style={styles.value}>Show less...</Text>
      </TouchableOpacity>
    );
  } else {
    return (
      <TouchableOpacity style={styles.listItem} onPress={show}>
        <Text style={styles.value}>Show more...</Text>
      </TouchableOpacity>
    );
  }
}

interface UserPositionsProps {
  userPositions: Immutable.List<UserPosition>;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  allowEditing?: boolean; // default false
  // Below are required when allowEditing is true
  removePosition?: ActionCreator<ThunkAction<Promise<ActionTypes>, ProfileState, void>>;
  errorToast?(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

interface UserTraitsState {
  showAll: boolean,
}

export class UserPositions extends Component<UserPositionsProps, UserTraitsState> {
  constructor(props: UserPositionsProps) {
    super(props);

    this.state = { showAll: false }

    this.renderPosition = this.renderPosition.bind(this);
  }

  private renderPosition(pos: UserPosition) {
    const dateFmt = "MMM YYYY";
    const until = !pos.endDate ? 'present' : Moment(pos.endDate).format(dateFmt);
    const frm = Moment(pos.startDate).format(dateFmt);

    const onRemoveAccept = async () => {
      try {
        await this.props.removePosition(pos.id);
      } catch (e) {
        await this.props.errorToast(e.errorMsg);
      }
    }

    const onRemovePress = () => {
      Alert.alert(
        'Remove Position',
        `Are you sure you want to remove your position as ${pos.roleName} at ${pos.organizationName}?`,
        [
          {text: 'Cancel', onPress: () => null, style: 'cancel'},
          {text: 'Remove', onPress: onRemoveAccept, style: 'destructive'},
        ],
      );
    }

    return (
      <View key={ pos.id } style={styles.positionContainer}>
        <Text style={styles.positionText}>
          <Text style={styles.positionBold}>{ pos.roleName }</Text>
          <Text> @ </Text>
          <Text style={styles.positionBold}>{ pos.organizationName }</Text>
          <Text>{'\n'}({ frm } - { until })</Text>
        </Text>
        <TouchableOpacity style={styles.traitDelete} onPress={onRemovePress}>
          <MaterialIcons color={Colors.WHITE} name="close" size={18} />
        </TouchableOpacity>
      </View>
    );
  }

  render() {
    let { userPositions } = this.props;
    const { showAll } = this.state;

    const addPosition = () => this.props.navigation.navigate('AddPosition');

    let bottomAction: ReactNode = null;
    if (userPositions.isEmpty()) {
      bottomAction = [
        <Text key={'text'} style={styles.noTraitText}>You don't have any positions</Text>,
        <Button
          key={'button'}
          buttonStyle={styles.noTraitButton}
          title="Add position"
          onPress={addPosition}
          color={Colors.HIVE_ACCENT}
        />,
      ];
    } else if (userPositions.size > MAX_NUMBER_POSITIONS_SHOWN) {
      bottomAction = renderShowLessMore(showAll,
        () => this.setState({ showAll: true }),
        () => this.setState({ showAll: false }));
    }

    userPositions = userPositions.sortBy(pos => pos.startDate).reverse().toList();
    if (!showAll) {
      userPositions = userPositions.take(MAX_NUMBER_POSITIONS_SHOWN).toList();
    }
    const positionItems = userPositions.map(this.renderPosition).toJS();

    return (
      <View style={styles.sectionContainer}>
        <Text style={styles.sectionHeader}>Positions</Text>
        <TouchableOpacity onPress={addPosition} style={styles.addTraitButton}>
          <MaterialIcons name="add-circle" size={32} color={Colors.HIVE_ACCENT} />
        </TouchableOpacity>
        { positionItems }
        <View style={styles.traitBottomActionContainer}>
          { bottomAction }
        </View>
      </View>
    );
  }
}

interface UserSimpleTraitsProps {
  userSimpleTraits: Immutable.List<UserSimpleTrait>;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  allowEditing?: boolean; // default false
  // Below are required when allowEditing is true
  removeSimpleTrait?: ActionCreator<ThunkAction<Promise<ActionTypes>, ProfileState, void>>;
  errorToast?(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

export class UserSimpleTraits extends Component<UserSimpleTraitsProps, UserTraitsState> {
  constructor(props: UserSimpleTraitsProps) {
    super(props);

    this.state = { showAll: false }

    this.renderSimpleTrait = this.renderSimpleTrait.bind(this);
  }

  private renderSimpleTrait(trait: UserSimpleTrait) {
    const onRemoveAccept = async () => {
      try {
        await this.props.removeSimpleTrait(trait.id);
      } catch (e) {
        await this.props.errorToast(e.errorMsg);
      }
    };

    const onRemovePress = () => {
      Alert.alert(
        'Remove Trait',
        `Are you sure you want to remove "${trait.simpleTraitName}" as a trait?`,
        [
          {text: 'Cancel', onPress: () => null, style: 'cancel'},
          {text: 'Remove', onPress: onRemoveAccept, style: 'destructive'},
        ],
      );
    }

    return (
      <View key={ trait.id } style={styles.simpleTraitContainer}>
        <Text style={styles.simpleTraitText}>{ trait.simpleTraitName }</Text>
        <TouchableOpacity style={styles.traitDelete} onPress={onRemovePress}>
          <MaterialIcons color={Colors.WHITE} name="close" size={18} />
        </TouchableOpacity>
      </View>
    );
  }

  render() {
    let { userSimpleTraits } = this.props;
    const { showAll } = this.state;

    const addSimpleTrait = () => this.props.navigation.navigate('AddSimpleTrait');

    let bottomAction: ReactNode = null;
    if (userSimpleTraits.isEmpty()) {
      bottomAction = [
        <Text key={'text'} style={styles.noTraitText}>You don't have any traits</Text>,
        <Button
          key={'button'}
          buttonStyle={styles.noTraitButton}
          title="Add trait"
          onPress={addSimpleTrait}
          color={Colors.HIVE_PRIMARY}
        />,
      ];
    } else if (userSimpleTraits.size > MAX_NUMBER_SIMPLE_TRAITS_SHOWN) {
      bottomAction = renderShowLessMore(showAll,
        () => this.setState({ showAll: true }),
        () => this.setState({ showAll: false }));
    }

    if (!showAll) {
      userSimpleTraits = userSimpleTraits.take(MAX_NUMBER_SIMPLE_TRAITS_SHOWN).toList();
    }
    const traitItems = userSimpleTraits.map(this.renderSimpleTrait).toJS();

    return (
      <View style={styles.sectionContainer}>
        <Text style={styles.sectionHeader}>Traits</Text>
        <TouchableOpacity onPress={addSimpleTrait} style={styles.addTraitButton}>
          <MaterialIcons name="add-circle" size={32} color={Colors.HIVE_PRIMARY} />
        </TouchableOpacity>
        <View style={styles.simpleTraitOuterContainer}>
          { traitItems }
        </View>
        <View style={styles.traitBottomActionContainer}>
          { bottomAction }
        </View>
      </View>
    );
  }
};

const BUTTON_WIDTH = SCREEN_WIDTH - 80;

export const styles = StyleSheet.create({
  container: {
    paddingTop: 10,
    paddingBottom: 65,
    backgroundColor: 'white',
    minHeight: '100%'
  },
  contentContainer: {
    alignItems: 'center',
    flex: 1,
    flexDirection: 'column',
    padding: 20,
  },
  description: {
    fontSize: 18,
    color: Colors.HIVE_SUBDUED
  },
  listItem: {
    flex: 1,
    flexDirection: 'row',
    marginTop: 5,
  },
  changePassButton: {
    width: BUTTON_WIDTH,
    marginTop: 10,
  },
  logoutButtonText: {
    color: 'white',
  },
  logoutButton: {
    width: BUTTON_WIDTH,
    marginTop: 10,
    backgroundColor: "gray",
    borderWidth: 0,
  },
  profileTitle: {
    fontSize: 18,
    marginTop: 10,
    alignSelf: 'flex-end'
  },
  profilePicture: {
    margin: 20
  },
  sectionHeader: {
    fontWeight: 'bold',
    fontSize: 24,
    alignSelf: 'flex-start',
  },
  sectionContainer: {
    width: "100%",
    flex: 1,
    flexDirection: 'column',
    marginTop: 20
  },
  subHeaderText: {
    fontSize: 18
  },
  label: {
    fontSize: 16,
  },
  value: {
    fontSize: 16,
    color: Colors.HIVE_ACCENT
  },
  cohortText: {
    fontSize: 16,
    marginTop: 5,
  },
  positionContainer: {
    backgroundColor: Colors.HIVE_ACCENT,
    marginTop: 5,
    padding: 5,
    paddingRight: 25,
    borderRadius: 5,
  },
  positionText: {
    color: Colors.WHITE,
    fontSize: 14,
  },
  positionBold: {
    fontWeight: '700',
  },
  traitDelete: {
    position: 'absolute',
    padding: 5,
    top: 0,
    right: 0,
  },
  simpleTraitOuterContainer: {
    flexWrap: 'wrap',
    flexDirection: 'row',
    justifyContent: 'center',
    paddingVertical: 2.5,
  },
  simpleTraitContainer: {
    backgroundColor: Colors.HIVE_PRIMARY,
    marginVertical: 2.5,
    marginHorizontal: 2.5,
    padding: 5,
    paddingRight: 25,
    paddingLeft: 10,
    borderRadius: 20,
  },
  simpleTraitText: {
    color: Colors.WHITE,
    fontWeight: '700',
    fontSize: 14,
  },
  traitBottomActionContainer: {
    flexDirection: 'column',
    alignItems: 'center',
  },
  addTraitButton: {
    position: 'absolute',
    top: 0,
    right: 0,
  },
  noTraitText: {
    marginTop: 10,
  },
  noTraitButton: {
    width: 200,
    marginTop: 10,
  },
  personalInfoContainer: {
    flex: 1,
    alignItems: 'center',
  },
});
