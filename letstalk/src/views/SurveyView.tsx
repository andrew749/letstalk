import React, { Component, ReactNode } from 'react';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import {Picker, ScrollView, ScrollViewComponent, StyleSheet, Text, View} from 'react-native';
import {
  NavigationScreenProp,
  NavigationStackAction,
  NavigationActions
} from 'react-navigation';
import { FormValidationMessage } from 'react-native-elements';
import Immutable from 'immutable';

import { RootState } from '../redux';
import { State } from '../redux/onboarding/reducer';
import {
  Survey,
  SurveyResponses,
} from '../models';
import {
  setOnboardingStatusAction,
  SetOnboardingStatusAction,
} from '../redux/onboarding/actions';
import Loading from './Loading';
import {
  ActionButton,
  Emoji,
  FormP,
  FormProps,
  Header,
  InfoText,
  LabeledFormInput,
  ModalPicker,
  Rating,
} from '../components';
import profileService, {
  PersonalityVector,
  UserVectorPreferenceType,
} from '../services/profile-service';
import { State as BootstrapState, fetchBootstrap } from '../redux/bootstrap/reducer';
import { ActionTypes as BootstrapActionTypes} from '../redux/bootstrap/actions';
import { State as SurveyState } from '../redux/survey/reducer';
import {
  ActionTypes as SurveyActionTypes,
  setSurveyQuestionsAction,
  setSurveyResponsesAction,
  setSurveyStateAction
} from '../redux/survey/actions';
import {
  getCohortId,
  programOptions,
  sequenceOptions,
  gradYearOptions,
  ValueLabel,
} from '../models/cohort';
import { AnalyticsHelper } from '../services';
import Colors from '../services/colors';
import {headerStyle, headerTintColor, headerTitleStyle} from "./TopHeader";
import {SurveyQuestion} from "../models/survey";

interface DispatchActions {
  fetchBootstrap: ActionCreator<ThunkAction<Promise<BootstrapActionTypes>, BootstrapState, void>>;
  setSurveyStateAction: ActionCreator<ThunkAction<Promise<SurveyActionTypes>, SurveyState, void>>;
  setSurveyQuestionsAction: ActionCreator<ThunkAction<Promise<SurveyActionTypes>, SurveyState, void>>;
  setSurveyResponsesAction: ActionCreator<ThunkAction<Promise<SurveyActionTypes>, SurveyState, void>>;
}

interface Props extends SurveyState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class SurveyView extends Component<Props> {
  SURVEY_VIEW_IDENTIFIER = "SurveyView";

  static navigationOptions = {
    headerTitle: 'SurveyView',
    headerStyle, 
    headerTitleStyle, 
    headerTintColor 
  }

  constructor(props: Props) {
    super(props);

    this.onSubmit = this.onSubmit.bind(this);
    this.load = this.load.bind(this);
    this.renderBody = this.renderBody.bind(this);
  }

  async componentDidMount() {
    AnalyticsHelper.getInstance().recordPage(this.ONBOARDING_VIEW_IDENTIFIER);
    await this.load();
  }

  private async load() {
    await this.props.fetchBootstrap();
  }

  async componentWillReceiveProps(nextProps: Props) {
    if (nextProps.bootstrap && nextProps.bootstrap.survey) {
      if (!this.props.survey) {
        await this.props.setSurveyStateAction(nextProps.bootstrap.survey)
      } else {
        await this.props.setSurveyQuestionsAction(nextProps.bootstrap.survey.questions)
      }
    }
  }

  async onSkipRemaining() {
    // TODO show confirmation modal that says how to complete the survey
  }

  async onSubmit() {
    try {
      return
      // TODO
    } catch(e) {
      throw new SubmissionError({_error: e.errorMsg});
    }
  }

  renderQuestion(question: SurveyQuestion) {
    const response = !responses ? null : responses.get(question.key);
    return (
      <View key={question.key}>
        <Header>{question.prompt}</Header>
      </View>
    );
  };

  renderBody() {
    const survey = this.props.survey;
    if (!survey) {
      return <Header>TODO hol up fam</Header>
    }
    const { questions, responses } = survey;
    return (
      <ViewPager
      >
        {questions.map(renderQuestion)}
      </ViewPager>
    );
  }

  render() {
    const {
      state,
      errorMsg,
      errorType,
    } = this.props.cohorts.fetchState;
    return (
      <Loading
        state={state}
        errorMsg={errorMsg}
        errorType={errorType}
        load={this.load}
        renderBody={this.renderBody}
        navigation={this.props.navigation}
      />
    );
  }
}

const styles = StyleSheet.create({
  actionButton: {
    marginBottom: 10,
  },
  cohortForm: {
    paddingBottom: 100,
  },
  hint: {
    color: 'gray',
    fontSize: 14,
    marginTop: -10,
    marginLeft: 10,
  },
  longForm: {
    height: 40
  }
});

export default connect(({ bootstrap }: RootState) => bootstrap,
  { fetchBootstrap, setSurveyStateAction, setSurveyQuestionsAction, setSurveyResponsesAction,  })(SurveyView);
