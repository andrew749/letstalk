import Immutable from 'immutable';
import React, { Component } from 'react';
import { connect, ActionCreator } from 'react-redux';
import { ThunkAction } from 'redux-thunk';
import { Dimensions, Platform, StyleSheet, TouchableOpacity, View, ScrollView } from 'react-native';
import {
  NavigationScreenProp,
  NavigationStackAction,
  StackNavigator,
  NavigationActions
} from 'react-navigation';
import { RootState } from '../redux';
import { SubmissionError } from 'redux-form'
import Loading from './Loading';
import {fetchSurvey, State as SurveyState} from '../redux/survey/reducer';
import {
  ActionTypes as SurveyActionTypes,
  setSurveyResponsesAction,
  SetSurveyResponsesAction
} from '../redux/survey/actions';
import { AnalyticsHelper } from '../services';
import {headerStyle, headerTintColor, headerTitleStyle} from "./TopHeader";
import {SurveyOption, SurveyQuestion, SurveyResponses} from "../models/survey";
import {Text} from "react-native-elements";
import Colors from "../services/colors";
import ActionButton from "../components/ActionButton";
import surveyService, { GROUP_GENERIC } from "../services/survey";
import {
  State as ProfileState,
  fetchProfile,
  setSurvey
} from '../redux/profile/reducer';
import { ActionTypes as ProfileActionTypes } from '../redux/profile/actions';

interface DispatchActions {
  fetchSurvey: ActionCreator<ThunkAction<Promise<SurveyActionTypes>, SurveyState, void>>;
  setSurvey: ActionCreator<ThunkAction<Promise<ProfileActionTypes>, ProfileState, void>>;
  setSurveyResponsesAction(surveyResponses: SurveyResponses): SetSurveyResponsesAction;
}

interface NavigationParams {
  currentQuestion: number;
}

interface NestedProps extends SurveyState, DispatchActions {
  navigation: NavigationScreenProp<void, NavigationStackAction & NavigationParams>;
}

class NestedSurveyViewComponent extends Component<NestedProps> {
  SURVEY_VIEW_IDENTIFIER = "SurveyView";

  static navigationOptions = {
    headerTitle: 'Questionnaire',
    headerStyle,
    headerTitleStyle,
    headerTintColor
  };

  constructor(props: NestedProps) {
    super(props);

    this.onSubmit = this.onSubmit.bind(this);
    this.load = this.load.bind(this);
    this.renderBody = this.renderBody.bind(this);
  }

  async componentDidMount() {
    AnalyticsHelper.getInstance().recordPage(this.SURVEY_VIEW_IDENTIFIER);
  }

  async load () {
  }

  private currentQuestion () {
    return this.props.navigation.getParam('currentQuestion', 0);
  }

  private navigateNextQuestion () {
    const { questions } = this.props.survey;
    const nextQuestion = this.currentQuestion() + 1;
    if (nextQuestion >= questions.size) {
      return;
    }
    this.props.navigation.push("SurveyView", {currentQuestion: nextQuestion});
  }

  private async onSubmit () {
    const { survey } = this.props;
    if (!survey.responses) {
      await this.props.navigation.pop(survey.questions.size);
    } else {
      try {
        await surveyService.postSurveyResponses(survey);
        await this.props.fetchSurvey(survey.group);
        if (survey.group !== GROUP_GENERIC) {
          await this.props.setSurvey(survey);
        }
        await this.props.navigation.pop(survey.questions.size);
      } catch(e) {
        console.error("error submitting responses", e);
        throw new SubmissionError({_error: e.errorMsg});
      }
    }
  }

  // Submit having not answered all questions.
  private async onSkipRemaining () {
    await this.onSubmit();
  }

  private async updateResponse (questionKey : string, optionKey : string) {
    let { responses } = this.props.survey;
    responses = !responses
      ? Immutable.Map({[questionKey]: optionKey})
      : responses.set(questionKey, optionKey);
    this.props.setSurveyResponsesAction(responses);
    setTimeout(() => this.navigateNextQuestion(), 200);
  }

  private renderSurveyOption (question: SurveyQuestion, option: SurveyOption, isSelected: boolean) {
    return (
      <TouchableOpacity
        key={option.key}
        style={[styles.surveyOption, isSelected && styles.surveyOptionSelected ]}
        onPress={() => this.updateResponse(question.key, option.key)}>
        <Text style={[styles.surveyOptionText, isSelected && styles.surveyOptionSelectedText]}>
          {option.text}
        </Text>
      </TouchableOpacity>)
  }

  private renderQuestion (question: SurveyQuestion) {
    const { responses } = this.props.survey;
    const response = !responses ? null : responses.get(question.key);
    const { options } = question;
    return (
      <View key={question.key} style={styles.questionCard}>
        <View style={styles.questionPromptContainer}>
          <Text style={styles.surveyQuestionText}>{question.prompt}</Text>
        </View>
        { options.map(option => this.renderSurveyOption(question, option, response === option.key)) }
      </View>
    );
  }

  private renderBody () {
    const { survey } = this.props;
    if (!survey) {
      return <View/>;
    }
    const currentQuestion = this.currentQuestion();
    const { questions, responses } = survey;
    const all_answered = responses && questions.every((question : SurveyQuestion) => responses.has(question.key));
    const quesNum = questions.size;
    const ansNum = responses ? responses.size : 0;
    return (
      <ScrollView style={styles.container}>
        <View style={styles.contentContainer}>
          {this.renderQuestion(questions.get(currentQuestion))}
        </View>
        <View style={[styles.footer]}>
          { currentQuestion < questions.size - 1
          ? (<ActionButton
              backgroundColor={Colors.WHITE}
              buttonStyle={[styles.actionButton, styles.skipButton]}
              textStyle={[styles.buttonText, styles.skipButtonText]}
              loading={false}
              title={'Skip Question'}
              onPress={() => this.navigateNextQuestion()}
            />)
          : all_answered
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
                title={'Submit ' + ansNum + '/' + quesNum + ' Responses'}
                onPress={() => this.onSkipRemaining()}
              />)
        }
        </View>
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
    width: '90%',
    marginTop: 10,
    padding: 10,
    backgroundColor: 'white',
    borderColor: Colors.HIVE_PRIMARY,
    borderWidth: .7,
    borderRadius: 4
  },
  surveyOptionText: {
    color: Colors.HIVE_SUBDUED,
    fontSize: 15,
  },
  surveyOptionSelected: {
    backgroundColor: Colors.HIVE_PRIMARY,
  },
  surveyOptionSelectedText: {
    color: 'white'
  },
  surveyQuestionText: {
    fontSize: 18,
    color: 'white',
  },
  questionCard: {
    width: SCREEN_WIDTH - 80,
    alignItems: 'center',
  },
  questionPromptContainer: {
    backgroundColor: Colors.HIVE_PRIMARY,
    borderRadius: 4,
    padding: 18,
    marginBottom: 20
  },
  actionButton: {
    width: SCREEN_WIDTH - 80,
    marginTop: 10,
    marginBottom: 30,
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
    flex: 1,
    flexDirection: 'column',
    paddingTop: 10,
    paddingBottom: 10,
    backgroundColor: Colors.WHITE,
    minHeight: '100%',
  },
  contentContainer: {
    paddingTop: 20,
    paddingBottom: 20,
    alignSelf: 'center',
  },
  footer: {
    alignSelf: 'center',
    alignItems: 'center',
  },
});

const SurveyView = connect(({ survey } : RootState) => survey,
  { fetchSurvey, setSurveyResponsesAction, setSurvey })(NestedSurveyViewComponent);

export default SurveyView;
