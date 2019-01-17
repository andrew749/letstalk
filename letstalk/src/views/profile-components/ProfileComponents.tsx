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
import { UserGroupSurvey } from '../../models/profile';
import Colors from '../../services/colors';
import Window from '../../services/window';
import { RootState } from '../../redux';
import { genderIdToString, GenderId } from '../../models/user';
import { State as SurveyState } from '../../redux/survey/reducer';
import { ActionTypes as SurveyActionTypes } from '../../redux/survey/actions';
import { State as ProfileState } from '../../redux/profile/reducer';
import { ActionTypes as ProfileActionTypes } from '../../redux/profile/actions';

interface PersonalInfoProps {
	firstName: string;
	lastName: string;
	gender: GenderId;
	birthdate?: string
  secret?: string // only required if `allowQrCode` is true
  bio?: string;
  hometown?: string;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  allowQrCode?: boolean; // default false
  showConnectedBadge?: boolean; // default false
  allowEditing?: boolean;
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
      allowQrCode,
      showConnectedBadge,
      allowEditing
    } = this.props;

    const badge = !!showConnectedBadge ? (
      <View style={styles.badgeContainer}>
        <Text style={styles.badgeText}>Connected</Text>
      </View>
    ) : null;

    const capitalize = (s: string) => s.charAt(0).toUpperCase() + s.slice(1);
    const genderStr = capitalize(genderIdToString(gender));

    const hometownStr = hometown === null || hometown === '' ? 'Some place on Earth' : hometown;

    let age;

    if (birthdate) {
      const timeDiff = new Date().valueOf() - new Date(birthdate).valueOf();
      age = Math.floor(timeDiff / (1000 * 60 * 60 * 24 * 365));
    } else {
      age = "";
    }

    const headerText = firstName + ' ' + lastName;

    // If allow QrCode is false, we assume that it's not your profile, but someone else's.
    const bioStr = !!bio ? bio : (
      !allowQrCode ? 'An awesome person that forgot to write a bio' : 'Add a bio through editing your profile!'
    );

    const updatePersonal = () => this.props.navigation.navigate('UpdatePersonal', {});

    return (
      <View style={styles.personalInfoContainer}>
        <Header textStyle={styles.headerText}>{headerText}</Header>
        {badge}
        <Text style={styles.subHeaderText}>{age}{gender != GenderId.Unspecified ? genderStr[0] : ''}{(!age && gender == GenderId.Unspecified) ? '' : ' - '}{hometownStr}</Text>
        {!!allowQrCode && <TouchableOpacity style={styles.listItem} onPress={() => {
          this.props.navigation.navigate('QrCode', { secret });
        }}>
          <Text style={styles.value}>Show QR Code</Text>
        </TouchableOpacity>}
        <View style={styles.sectionContainer}>
          <Text style={styles.description}>{ bioStr }</Text>
        </View>
        {!!allowEditing && <TouchableOpacity onPress={updatePersonal} style={styles.addTraitButton}>
          <MaterialIcons name="edit" size={25} color={Colors.HIVE_PRIMARY} />
        </TouchableOpacity>}
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
      allowEditing,
    } = this.props;
    const program = programById(programId);
    const changeCohort = () => this.props.navigation.navigate('ChangeCohort', {
      programId,
      gradYear,
      sequenceId,
    });
    return (
      <View style={styles.sectionContainer}>
        <Text style={styles.sectionHeader}>Cohort</Text>
        <Text style={styles.cohortText}>{ program + ', ' + gradYear }</Text>
        {!!allowEditing && <TouchableOpacity onPress={changeCohort} style={styles.addTraitButton}>
          <MaterialIcons name="edit" size={25} color={Colors.HIVE_PRIMARY} />
        </TouchableOpacity>}
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

interface UserGroupSurveysProps {
  userGroupSurveys: Immutable.List<UserGroupSurvey>;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  fetchSurvey: ActionCreator<ThunkAction<Promise<SurveyActionTypes>, SurveyState, void>>;
  removeGroup: ActionCreator<ThunkAction<Promise<ProfileActionTypes>, ProfileState, void>>;
  errorToast?(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
  allowEditing?: boolean; // default false
}

interface UserGroupSurveysState {
  showAll: boolean,
}

export class UserGroupSurveys extends Component<UserGroupSurveysProps, UserGroupSurveysState> {
  constructor(props: UserGroupSurveysProps) {
    super(props);

    this.state = { showAll: false }
    this.renderGroupSurvey = this.renderGroupSurvey.bind(this);
  }

  private renderGroupSurvey(groupSurvey: UserGroupSurvey) {
    const { allowEditing } = this.props;

    const numQuestions = groupSurvey.survey.questions.size;
    const numAnswers = groupSurvey.survey.responses ? groupSurvey.survey.responses.size : 0;

    const color = numQuestions === numAnswers ? Colors.HIVE_ACCENT : Colors.HIVE_ERROR;
    const onEditPress = () => {
      this.props.fetchSurvey(groupSurvey.survey.group);
      this.props.navigation.navigate('SurveyView', {});
    }

    const onRemoveAccept = async () => {
      try {
        await this.props.removeGroup(groupSurvey.userGroup.id);
      } catch (e) {
        await this.props.errorToast(e.errorMsg);
      }
    }

    const onRemovePress = () => {
      Alert.alert(
        'Leave Group',
        `Are you sure you want to leave the ${groupSurvey.userGroup.groupName} group`,
        [
          {text: 'Cancel', onPress: () => null, style: 'cancel'},
          {text: 'Leave', onPress: onRemoveAccept, style: 'destructive'},
        ],
      );
    }

    const padRight = allowEditing ? 50 : 25;
    const extraStyle = { paddingRight: padRight };
    const editSurveyRight = allowEditing ? 25 : 0;
    const editSurveyExtraStyle = { right: editSurveyRight };
    return (
      <View
        key={ groupSurvey.survey.group }
        style={[styles.pillContainer, { backgroundColor: color, borderColor: color }, extraStyle]}
      >
        <Text style={styles.pillText}>
          <Text style={styles.boldText}>{ groupSurvey.userGroup.groupName }</Text>
          <Text>{ `\n${numAnswers} out of ${numQuestions} questions completed` }</Text>
        </Text>
        <TouchableOpacity style={[styles.editSurvey, editSurveyExtraStyle]} onPress={onEditPress}>
          <MaterialIcons name="edit" size={18} color={Colors.WHITE} />
        </TouchableOpacity>
        {!!allowEditing && <TouchableOpacity style={styles.traitDelete} onPress={onRemovePress}>
          <MaterialIcons color={Colors.WHITE} name="close" size={18} />
        </TouchableOpacity>}
      </View>
    );
  }

  render() {
    let { userGroupSurveys, allowEditing } = this.props;
    const { showAll } = this.state;

    const addGroup = () => this.props.navigation.navigate('AddGroup');
    let bottomAction: ReactNode = null;
    if (userGroupSurveys.isEmpty()) {
      bottomAction = [
        <Text key={'text'} style={styles.noTraitText}>
          { 'No group surveys for you to complete' }
        </Text>,
      ];
    } else if (userGroupSurveys.size > MAX_NUMBER_POSITIONS_SHOWN) {
      bottomAction = renderShowLessMore(showAll,
        () => this.setState({ showAll: true }),
        () => this.setState({ showAll: false }));
    }

    if (!showAll) {
      userGroupSurveys = userGroupSurveys.take(MAX_NUMBER_POSITIONS_SHOWN).toList();
    }
    const surveyItems = userGroupSurveys.map(this.renderGroupSurvey).toJS();

    return (
      <View style={styles.sectionContainer}>
        <Text style={styles.sectionHeader}>Group Surveys</Text>
        <View style={styles.pillsOuterContainer}>
          { surveyItems }
        </View>
        <View style={styles.traitBottomActionContainer}>
          { bottomAction }
        </View>
        {!!allowEditing && <TouchableOpacity onPress={addGroup} style={styles.addTraitButton}>
          <MaterialIcons name="add-circle" size={32} color={Colors.HIVE_PRIMARY} />
        </TouchableOpacity>}
      </View>
    );
  }
}

interface UserPositionsProps {
  userPositions: Immutable.List<UserPosition>;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  allowEditing?: boolean; // default false
  // Below are required when allowEditing is true
  removePosition?: ActionCreator<ThunkAction<Promise<ProfileActionTypes>, ProfileState, void>>;
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
    const { allowEditing } = this.props;

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

    const padRight = allowEditing ? { paddingRight: 25 } : null;
    return (
      <View key={ pos.id } style={[styles.pillContainer, padRight]}>
        <Text style={styles.pillText}>
          <Text style={styles.boldText}>{ pos.roleName } @ {pos.organizationName }</Text>
          <Text>{'\n'}({ frm } - { until })</Text>
        </Text>
        {!!allowEditing && <TouchableOpacity style={styles.traitDelete} onPress={onRemovePress}>
          <MaterialIcons color={Colors.WHITE} name="close" size={18} />
        </TouchableOpacity>}
      </View>
    );
  }

  render() {
    let { userPositions, allowEditing } = this.props;
    const { showAll } = this.state;

    const addPosition = () => this.props.navigation.navigate('AddPosition');
    const emptyText = allowEditing ? 'You don\'t have any positions, press the + to add one!' : 'They don\'t have any positions';

    let bottomAction: ReactNode = null;
    if (userPositions.isEmpty()) {
      bottomAction = [
        <Text key={'text'} style={styles.noTraitText}>{ emptyText }</Text>,
        // !!allowEditing && <Button
        //   key={'button'}
        //   buttonStyle={styles.noTraitButton}
        //   title="Add position"
        //   onPress={addPosition}
        //   color={Colors.HIVE_PRIMARY}
        // />,
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
        <View style={styles.pillsOuterContainer}>
          { positionItems }
        </View>
        <View style={styles.traitBottomActionContainer}>
          { bottomAction }
        </View>
        {!!allowEditing && <TouchableOpacity onPress={addPosition} style={styles.addTraitButton}>
          <MaterialIcons name="add-circle" size={32} color={Colors.HIVE_PRIMARY} />
        </TouchableOpacity>}
      </View>
    );
  }
}

interface UserSimpleTraitsProps {
  userSimpleTraits: Immutable.List<UserSimpleTrait>;
  navigation: NavigationScreenProp<void, NavigationStackAction>;
  allowEditing?: boolean; // default false
  // Below are required when allowEditing is true
  removeSimpleTrait?: ActionCreator<ThunkAction<Promise<ProfileActionTypes>, ProfileState, void>>;
  errorToast?(message: string): (dispatch: Dispatch<RootState>) => Promise<void>;
}

export class UserSimpleTraits extends Component<UserSimpleTraitsProps, UserTraitsState> {
  constructor(props: UserSimpleTraitsProps) {
    super(props);

    this.state = { showAll: false }

    this.renderSimpleTrait = this.renderSimpleTrait.bind(this);
  }

  private renderSimpleTrait(trait: UserSimpleTrait) {
    const { allowEditing } = this.props;
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

    const padRight = allowEditing ? { paddingRight: 25 } : null;
    return (
      <View key={ trait.id } style={[styles.pillContainer, padRight]}>
        <Text style={[styles.pillText, styles.boldText]}>{ trait.simpleTraitName }</Text>
        {!!allowEditing && <TouchableOpacity style={styles.traitDelete} onPress={onRemovePress}>
          <MaterialIcons color={Colors.WHITE} name="close" size={18} />
        </TouchableOpacity>}
      </View>
    );
  }

  render() {
    let { userSimpleTraits, allowEditing } = this.props;
    const { showAll } = this.state;

    const addSimpleTrait = () => this.props.navigation.navigate('AddSimpleTrait');
    const emptyText = allowEditing ? 'You don\'t have any traits, press the + to add one!' : 'They don\'t have any traits';

    let bottomAction: ReactNode = null;
    if (userSimpleTraits.isEmpty()) {
      bottomAction = [
        <Text key={'text'} style={styles.noTraitText}>{ emptyText }</Text>,
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
        <View style={styles.pillsOuterContainer}>
          { traitItems }
        </View>
        <View style={styles.traitBottomActionContainer}>
          { bottomAction }
        </View>
        {!!allowEditing && <TouchableOpacity onPress={addSimpleTrait} style={styles.addTraitButton}>
          <MaterialIcons name="add-circle" size={32} color={Colors.HIVE_PRIMARY} />
        </TouchableOpacity>}
      </View>
    );
  }
};

const BUTTON_WIDTH = Window.WIDTH - 80;

export const styles = StyleSheet.create({
  addTraitButton: {
    position: 'absolute',
    top: 0,
    right: 0
  },
  badgeContainer: {
    padding: 5,
    paddingHorizontal: 10,
    borderRadius: 20,
    backgroundColor: Colors.GREEN,
    marginBottom: 5
  },
  badgeText: {
    color: Colors.WHITE,
    fontWeight: '900'
  },
  boldText: {
    fontWeight: '700'
  },
  buttonText: {
    fontSize: 16
  },
  changePassButton: {
    borderColor: Colors.HIVE_SUBDUED,
    borderWidth: .7,
    height: 44,
    margin: 0
  },
  changePassButtonText: {
    color: Colors.HIVE_SUBDUED
  },
  cohortText: {
    fontSize: 16,
    marginTop: 5
  },
  container: {
    paddingTop: 10,
    paddingBottom: 10,
    backgroundColor: 'white',
    minHeight: '100%'
  },
  contentContainer: {
    alignItems: 'center',
    flex: 1,
    flexDirection: 'column',
    padding: 20
  },
  description: {
    fontSize: 18,
    color: Colors.HIVE_SUBDUED
  },
  feedbackButton: {
    backgroundColor: Colors.WHITE,
    alignSelf: 'center',
    width: 200,
    marginVertical: 15,
  },
  feedbackButtonText: {
    color: Colors.HIVE_PRIMARY,
    fontSize: 16
  },
  headerText: {
    textAlign: 'center',
  },
  label: {
    fontSize: 16
  },
  listItem: {
    flex: 1,
    flexDirection: 'row',
    marginTop: 5
  },
  logoutButton: {
    backgroundColor: Colors.HIVE_SUBDUED,
    borderWidth: 0,
    height: 44,
    margin: 0
  },
  logoutButtonText: {
    color: 'white'
  },
  noTraitButton: {
    width: 200,
    marginTop: 10,
    padding: 10
  },
  noTraitText: {
    marginTop: 10
  },
  personalInfoContainer: {
    flex: 1,
    alignItems: 'center',
    width: "100%"
  },
  pickerContainer: {
    marginHorizontal: 15
  },
  pillContainer: {
    borderColor: Colors.HIVE_ACCENT,
    backgroundColor: Colors.HIVE_ACCENT,
    borderWidth: .7,
    marginVertical: 6,
    marginHorizontal: 2.5,
    padding: 8,
    paddingLeft: 10,
    paddingRight: 10,
    borderRadius: 5
  },
  pillsOuterContainer: {
    width: "100%",
    marginVertical: 6
  },
  pillText: {
    color: Colors.WHITE,
    fontWeight: '400',
    fontSize: 14
  },
  profileActionButton: {
    alignSelf: 'center',
    width: BUTTON_WIDTH,
    marginTop: 10,
    padding: 8
  },
  profilePicture: {
    margin: 20
  },
  profileTitle: {
    fontSize: 18,
    marginTop: 10,
    alignSelf: 'flex-end'
  },
  scrollViewContainer: {
    paddingTop: 10
  },
  sectionContainer: {
    width: "100%",
    flex: 1,
    flexDirection: 'column',
    marginTop: 20
  },
  sectionHeader: {
    fontWeight: 'bold',
    fontSize: 24,
    alignSelf: 'flex-start',
    marginBottom: 8
  },
  subHeaderText: {
    fontSize: 18
  },
  traitBottomActionContainer: {
    flexDirection: 'column',
    alignItems: 'center'
  },
  traitDelete: {
    position: 'absolute',
    padding: 5,
    top: 0,
    right: 0
  },
  editSurvey: {
    position: 'absolute',
    padding: 5,
    top: 0,
    right: 25
  },
  value: {
    fontSize: 16,
    color: Colors.HIVE_PRIMARY
  },
  moreAboutYou: {
    marginTop: 20
  },
});
