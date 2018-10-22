import React, { Component } from 'react';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import { StyleSheet, View} from 'react-native';
import {
  NavigationScreenProp,
  NavigationStackAction,
} from 'react-navigation';

import { RootState } from '../redux';
import {
  Survey,
  SurveyResponses,
} from '../models';
import Loading from './Loading';
import {
  Header,
} from '../components';
import {fetchSurvey, State as SurveyState} from '../redux/survey/reducer';
import {
  ActionTypes as SurveyActionTypes,
  setSurveyQuestionsAction,
  setSurveyResponsesAction,
  setSurveyStateAction
} from '../redux/survey/actions';
import { AnalyticsHelper } from '../services';
import {headerStyle, headerTintColor, headerTitleStyle} from "./TopHeader";
import {SurveyQuestion} from "../models/survey";

interface DispatchActions {
  fetchSurvey: ActionCreator<ThunkAction<Promise<SurveyActionTypes>, SurveyState, void>>;
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
    AnalyticsHelper.getInstance().recordPage(this.SURVEY_VIEW_IDENTIFIER);
    await this.load();
  }

  private async load() {
    await this.props.fetchSurvey();
  }

  onSkipRemaining = async () => {
    // TODO show confirmation modal that says how to complete the survey
  };

  onSubmit = async () => {
    const { responses } = this.props.survey;
    // TODO
  };

  renderQuestion = (question: SurveyQuestion) => {
    const { responses } = this.props.survey;
    const response = !responses ? null : responses.get(question.key);
    return (
      <View key={question.key}>
        <Header>{question.prompt}</Header>
      </View>
    );
  };

  renderBody = () => {
    const survey = this.props.survey;
    if (!survey) {
      return <Header>TODO hol up fam</Header>
    }
    const { questions, responses } = survey;
    return (
      <ViewPager>
        {questions.map(this.renderQuestion)}
      </ViewPager>
    );
  };

  render() {
    const {
      state,
      errorMsg,
      errorType,
    } = this.props.fetchState;
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
  { fetchSurvey, setSurveyStateAction, setSurveyQuestionsAction, setSurveyResponsesAction,  })(SurveyView);
