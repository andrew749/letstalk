import Immutable from 'immutable';

import {
  ActionTypes,
  TypeKeys,
} from './actions';
import {Survey, SurveyQuestions} from "../../models/survey";

export interface State {
  readonly survey?: Survey
}

const initialState: State = { };

// TODO need to test this
function updateQuestionsRemovingStaleResponses(oldState: State, questions: SurveyQuestions) {
  const { survey } = oldState;
  const { responses:currentResponses } = survey;
  if (!currentResponses) {
      return { ...oldState, survey: { ...survey, questions } };
  } else {
    const newQuestionKeys = questions.map(({key}) => key).toSet();
    const responses : Immutable.Map<string, string> =
      Immutable.Map(currentResponses.filterNot(({key}) => newQuestionKeys.has(key)));
    return { ...oldState, survey: { ...survey, questions, responses } };
  }
}

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.SET_STATE:
      return {
        ...state,
        survey: action.survey,
      };
    case TypeKeys.SET_QUESTIONS:
      return updateQuestionsRemovingStaleResponses(state, action.questions);
    case TypeKeys.SET_RESPONSES:
      return {
        ...state,
        survey: {
          ...state.survey,
          responses:action.responses,
        }
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action;
      return state;
  }
};
