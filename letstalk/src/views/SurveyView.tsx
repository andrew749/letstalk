import Immutable from 'immutable';
import React, { Component } from 'react';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import {Dimensions, ScrollView, StyleSheet, TouchableOpacity, TouchableOpacityComponent, View} from 'react-native';
import {
  NavigationScreenProp,
  NavigationStackAction,
} from 'react-navigation';
import { RootState } from '../redux';
import Loading from './Loading';
import {
  Header,
} from '../components';
import {fetchSurvey, State as SurveyState} from '../redux/survey/reducer';
import {
  ActionTypes as SurveyActionTypes,
  setSurveyResponsesAction,
  SetSurveyResponsesAction
} from '../redux/survey/actions';
import { AnalyticsHelper } from '../services';
import {headerStyle, headerTintColor, headerTitleStyle} from "./TopHeader";
import {SurveyQuestion, SurveyResponses} from "../models/survey";
import {Text} from "react-native-elements";
import Colors from "../services/colors";

interface DispatchActions {
  fetchSurvey: ActionCreator<ThunkAction<Promise<SurveyActionTypes>, SurveyState, void>>;
  setSurveyResponsesAction(surveyResponses: SurveyResponses): SetSurveyResponsesAction;
}

interface Props extends SurveyState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction>;
}

class SurveyView extends Component<Props> {
  SURVEY_VIEW_IDENTIFIER = "SurveyView";

  static navigationOptions = {
    headerTitle: 'Matching Survey',
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

  updateResponse = async (questionKey : string, optionKey : string) => {
    let { responses } = this.props.survey;
    responses = !responses
      ? Immutable.Map({[questionKey]: optionKey})
      : responses.set(questionKey, optionKey);
    await this.props.setSurveyResponsesAction(responses);
  };

  renderQuestion = (question: SurveyQuestion) => {
    const { responses } = this.props.survey;
    const response = !responses ? null : responses.get(question.key);
    const { options } = question;
    return (
      <View key={question.key} >
        <Header>{question.prompt}</Header>
        { options.map(option => {
          const isSelected = response === option.key;
          let { backgroundColor, textColor } = isSelected
            ? { backgroundColor: Colors.HIVE_PRIMARY, textColor: Colors.HIVE_MAIN_FONT }
            : { backgroundColor: Colors.HIVE_BG, textColor: Colors.HIVE_SUBDUED }
          return (
            <TouchableOpacity
              key={option.key}
              style={[styles.surveyOption, {backgroundColor}]}
              onPress={() => this.updateResponse(question.key, option.key)}>
            <Text style={[styles.surveyOptionText, {color: textColor}]}>
              {option.text}
            </Text>
          </TouchableOpacity>)
          })}
      </View>
    );
  };

  renderBody = () => {
    const survey = this.props.survey;
    if (!survey) {
      return <View/>;
    }
    const { questions, responses } = survey;
    return (
      <ScrollView>
        {questions.map(this.renderQuestion)}
      </ScrollView>
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

const SCREEN_WIDTH = Dimensions.get('window').width;
const styles = StyleSheet.create({
  surveyOption: {
    alignSelf: 'center',
    width: SCREEN_WIDTH - 120,
    marginTop: 10,
    padding: 10
  },
  surveyOptionText: {
    fontSize: 18,
  },
  actionButton: {
    alignSelf: 'center',
    width: SCREEN_WIDTH - 80,
    marginTop: 10,
    padding: 8
  },
});

export default connect(({ survey } : RootState) => survey,
  { fetchSurvey, setSurveyResponsesAction,  })(SurveyView);
