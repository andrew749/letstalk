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
import { SubmissionError } from 'redux-form'
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
import ActionButton from "../components/ActionButton";
import surveyService from "../services/survey";

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

  private async onSkipRemaining () {
    // TODO show confirmation modal that says how to complete the survey
    this.onSubmit();
  };

  private async onSubmit () {
    const { survey } = this.props;
    try {
      await surveyService.postSurveyResponses(survey);
      await this.props.fetchSurvey();
      await this.props.navigation.goBack();
    } catch(e) {
      console.error("error submitting responses", e);
      throw new SubmissionError({_error: e.errorMsg});
    }
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
      <View key={question.key} style={styles.questionCard}>
        <Text style={styles.surveyQuestionText}>{question.prompt}</Text>
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
    const all_answered = responses && questions.every((question : SurveyQuestion) => responses.has(question.key));
    return (
      <View style={styles.container}>
        <ScrollView style={styles.scrollView} contentContainerStyle={{alignItems:'center'}}>
          {questions.map(this.renderQuestion)}
        </ScrollView>
        <View style={[styles.footer]}>
          { all_answered
            ? (<ActionButton
                backgroundColor={Colors.HIVE_PRIMARY}
                buttonStyle={[styles.actionButton, styles.submitButton]}
                textStyle={[styles.buttonText, styles.submitButtonText]}
                loading={false}
                title={'Submit Responses'}
                onPress={() => this.onSubmit()}
              />)
            : (<ActionButton
                backgroundColor={Colors.WHITE}
                buttonStyle={[styles.actionButton, styles.skipButton]}
                textStyle={[styles.buttonText, styles.skipButtonText]}
                loading={false}
                title={'Skip Remaining'}
                onPress={() => this.onSkipRemaining()}
              />)
          }
        </View>
      </View>
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
    width: SCREEN_WIDTH - 120,
    marginTop: 10,
    padding: 10
  },
  surveyOptionText: {
    fontSize: 18,
  },
  surveyQuestionText: {
    fontSize: 24,
    color: Colors.HIVE_LIGHT_FONT,
  },
  questionCard: {
    alignItems: 'center',
    width: '90%',
    paddingBottom:20,
  },
  actionButton: {
    width: SCREEN_WIDTH - 80,
    marginTop: 10,
    padding: 8,
    minHeight: 44,
  },
  buttonText: {
    fontSize: 16
  },
  submitButton: {
    borderColor: Colors.HIVE_PRIMARY,
    borderWidth: .7,
    height: 44,
    margin: 0
  },
  submitButtonText: {
    color: Colors.HIVE_MAIN_FONT,
  },
  skipButton: {
    borderColor: Colors.HIVE_SUBDUED,
    borderWidth: .7,
    height: 44,
    margin: 0
  },
  skipButtonText: {
    color: Colors.HIVE_SUBDUED,
  },
  container: {
    display:'flex',
    paddingTop: 10,
    paddingBottom: 10,
    backgroundColor: Colors.WHITE,
    minHeight: '100%',
  },
  scrollView: {
    flex: 1,
    flexDirection: 'column',
    paddingTop: 20,
    paddingBottom: 20,
    alignSelf: 'center',
  },
  footer: {
    alignSelf: 'center',
    alignItems: 'center',
    flex: 1,
    flexDirection: 'column',
  },
});

export default connect(({ survey } : RootState) => survey,
  { fetchSurvey, setSurveyResponsesAction,  })(SurveyView);
