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
} from 'react-native';

import { Card } from '../components';
import { headerStyle } from './TopHeader';
import Colors from '../services/colors';
import { AnalyticsHelper } from '../services/analytics';
import { ProfileData } from '../models/profile';

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction & { profile: ProfileData }>;
}

export default class EditTraitsSelectorView extends Component<Props> {
  ADD_TRAIT_SELECTOR_VIEW_IDENTIFIER = 'EditTraitsSelectorView';

  static navigationOptions = {
    headerTitle: 'Edit Traits',
    headerStyle,
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.ADD_TRAIT_SELECTOR_VIEW_IDENTIFIER);
    });
  }

  private renderTraitCard(name: string, editView: string, description: string) {
    const onPress = async () => {
      const { navigation } = this.props;
      const profile = navigation.getParam('profile', null);
      await navigation.navigate(editView, profile);
    };
    return (
      <TouchableOpacity onPress={onPress}>
        <Card key={name} style={styles.traitCard}>
          <Text style={styles.traitTypeTitle}>{name}</Text>
          <Text style={styles.traitTypeDescription}>{description}</Text>
        </Card>
      </TouchableOpacity>
    );
  }

  render() {
    // TODO: Copies
    const traitTypes = [
      {
        name: 'Change Cohort',
        editView: 'ChangeCohort',
        description: 'Your program and anticipated graduating year. For students ' +
          'in co-op, this also includes your sequence.',
      },
    ];
    const traitCards = traitTypes.map(({ name, editView, description }) => {
      return this.renderTraitCard(name, editView, description)
    });

    const description =
      'Traits describe who you are. The more we know about you, the better we can help you find ' +
      'others that share common interests and aspirations, leading to awesome new friendships ' +
      'and valuable mentorships.';

    return (
      <ScrollView style={styles.container}>
        <Text>{ description }</Text>
        { traitCards }
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
    padding: 10,
    justifyContent: 'center',
    alignItems: 'center',
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
});
