import React, { Component } from 'react';
import { Text } from 'react-native';
import {
  NavigationScreenProp,
  NavigationStackAction,
} from "react-navigation";

import { headerStyle } from '../TopHeader';
import { AnalyticsHelper } from '../../services/analytics';

interface Props {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

export default class ChangeCohortView extends Component<Props> {
  CHANGE_COHORT_VIEW_IDENTIFIER = 'ChangeCohortView';

  static navigationOptions = {
    headerTitle: 'Change Cohort',
    headerStyle,
  }

  async componentDidMount() {
    this.props.navigation.addListener('willFocus', (route) => {
      AnalyticsHelper.getInstance().recordPage(this.CHANGE_COHORT_VIEW_IDENTIFIER);
    });
  }

  render() {
    return (
      <Text>
        Change da cohort
      </Text>
    );
  }
}
