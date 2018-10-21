import {
  ActionTypes, setSurvey,
  TypeKeys,
} from './actions';
import {Survey} from "../../models/survey";
import {ActionCreator, Dispatch} from "react-redux";
import {ThunkAction} from "redux-thunk";
import surveyService from "../../services/survey";

export interface State {
  readonly survey?: Survey
}

const initialState: State = { };

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.SET_SURVEY:
      return {
        ...state,
        survey: action.survey,
      };
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
  ThunkAction<Promise<ActionTypes>, State, void>> = () => {
  return async (dispatch: Dispatch<State>) => {
    await dispatch(fetch.start());
    try {
      const survey = await surveyService.getSurvey();
      await dispatch(setSurvey(survey));
      return dispatch(fetch.receive(data));
    } catch(e) {
      return dispatch(fetch.error(e));
    }
  };
};

export { fetchSurvey };
