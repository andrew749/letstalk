import Immutable from 'immutable'
import {
  ActionTypes,
  TypeKeys,
  fetch
} from './actions';
import {Survey} from "../../models/survey";
import {ActionCreator, Dispatch} from "react-redux";
import {ThunkAction} from "redux-thunk";
import surveyService from "../../services/survey";
import {FetchState, fetchStateReducer, getDataOrCur, initialFetchState} from "../actions";

export interface State {
  readonly survey?: Survey
  readonly fetchState: FetchState;
}

const initialState: State = {
  fetchState: initialFetchState,
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.FETCH:
      const survey = getDataOrCur(action, state.survey);
      if (!survey) {
        return {
          ...state,
          fetchState: fetchStateReducer(action),
          survey,
        };
      } else {
        const responses = survey.responses ? Immutable.Map(survey.responses) : null;
        return {
          ...state,
          fetchState: fetchStateReducer(action),
          survey: {
            ...survey,
            responses,
            questions: Immutable.List(survey.questions),
          }
        };
      }
    case TypeKeys.SET_SURVEY_RESPONSES:
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
}

const fetchSurvey: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (group: string) => {
  return async (dispatch: Dispatch<State>) => {
    await dispatch(fetch.start());
    try {
      const data = await surveyService.getSurvey(group);
      return dispatch(fetch.receive(data));
    } catch(e) {
      return dispatch(fetch.error(e));
    }
  };
};

export { fetchSurvey };
