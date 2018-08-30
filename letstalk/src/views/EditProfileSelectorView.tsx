import React, { Component } from 'react';
import {
  NavigationScreenProp,
  NavigationStackAction,
} from "react-navigation";
import {
  Text,
  TouchableOpacity,
  ScrollView,
  StyleSheet,
  View,
} from 'react-native';

import { Card, Header } from '../components';
import { headerStyle } from './TopHeader';
import Colors from '../services/colors';
import { AnalyticsHelper } from '../services/analytics';
import { ProfileData } from '../models/profile';

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction & { profile: ProfileData }>;
}

export default class EditProfileSelectorView extends Component<Props> {
  EDIT_PROFILE_SELECTOR_VIEW_IDENTIFIER = 'EditProfileSelectorView';

  static navigationOptions = {
    headerTitle: 'Edit Profile',
    headerStyle,
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.EDIT_PROFILE_SELECTOR_VIEW_IDENTIFIER);
    });
  }

  private renderTraitCard(name: string, editView: string, description: string) {
    const onPress = async () => {
      const { navigation } = this.props;
      const profile = navigation.getParam('profile', null);
      await navigation.navigate(editView, profile);
    };
    return (
      <TouchableOpacity key={name} onPress={onPress}>
        <Card style={styles.traitCard}>
          <View style={styles.nameContainer}>
            <Text style={styles.traitTypeTitle}>{name}</Text>
          </View>
          <Text style={styles.traitTypeDescription}>{description}</Text>
        </Card>
      </TouchableOpacity>
    );
  }

  render() {
    const editPersonalInfoCard = this.renderTraitCard(
      'Edit Personal Info',
      'ProfileEdit',
      'Edit info such as your name, phone number, mentorship preference, hometown and bio. Also, ' +
      'change your profile picture.',
    );

    // TODO: Copies
    const traitTypes = [
      {
        name: 'Change Cohort',
        editView: 'ChangeCohort',
        description: 'Your program and anticipated graduating year. For students ' +
          'in co-op, this also includes your sequence.',
      },
      {
        name: 'Add Position',
        editView: 'AddPosition',
        description: 'Add positions at companies, clubs or sports teams that you hold or have ' +
          'held in the past.',
      },
      {
        name: 'Add Trait',
        editView: 'AddSimpleTrait',
        description: 'Add any other traits such as your interests, hobbies, experiences, etc.',
      },
    ];
    const traitCards = traitTypes.map(({ name, editView, description }) => {
      return this.renderTraitCard(name, editView, description)
    });

    const traitDescription =
      'Traits describe who you are. The more we know about you, the better we can help you find ' +
      'others that share common interests and aspirations, leading to awesome new friendships ' +
      'and valuable mentorships.';

    return (
      <ScrollView style={styles.container}>
        <Header>Personal Info</Header>
        <View style={styles.buttonContainer}>
          { editPersonalInfoCard }
        </View>
        <Header>Traits</Header>
        <Text style={styles.descriptionText}>{ traitDescription }</Text>
        <View style={[styles.buttonContainer, { paddingBottom: 20 }]}>
          { traitCards }
        </View>
      </ScrollView>
    );
  }
}

const styles = StyleSheet.create({
  container: {
    padding: 10,
  },
  traitCard: {
    marginHorizontal: 0,
    marginVertical: 5,
    padding: 10,
  },
  buttonContainer: {
    paddingVertical: 5,
  },
  traitTypeTitle: {
    color: Colors.HIVE_MAIN_FONT,
    fontSize: 16,
    fontWeight: '900',
  },
  traitTypeDescription: {
    marginTop: 5,
    color: Colors.HIVE_MAIN_FONT,
    fontSize: 14,
  },
  nameContainer: {
    alignItems: 'center',
  },
  descriptionText: {
    paddingHorizontal: 10,
    paddingBottom: 10,
  },
});
