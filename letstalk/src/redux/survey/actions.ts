import { Action } from 'redux'
import {Survey, SurveyResponses} from "../../models/survey";
import {
  FetchReceiveAction,
  FetchErrorAction,
  FetchStartAction,
  FetchTypeKeys,
  FetchActionCreators,
} from '../actions';
import { APIError } from '../../services/requests';

export enum TypeKeys {
  FETCH = 'SURVEY/FETCH',
  SET_SURVEY = 'SURVEY/SET_SURVEY',
  SET_SURVEY_RESPONSES = 'SURVEY/SET_SURVEY_RESPONSES',
}

export interface SetSurveyAction extends Action {
  readonly type: TypeKeys.SET_SURVEY;
  readonly survey: Survey;
}

export interface SetSurveyResponsesAction extends Action {
  readonly type: TypeKeys.SET_SURVEY_RESPONSES;
  readonly responses: SurveyResponses;
}

type FetchSurveyReceiveAction = FetchReceiveAction<TypeKeys.FETCH, Survey>;
type FetchSurveyErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type FetchSurveyStartAction = FetchStartAction<TypeKeys.FETCH>;

export function setSurveyResponses(
  responses: SurveyResponses
): SetSurveyResponsesAction {
  return {
    type: TypeKeys.SET_SURVEY_RESPONSES,
    responses,
  };
}

function receive(data: Survey): FetchSurveyReceiveAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.RECEIVE,
    data,
  };
}

function error(error: APIError): FetchSurveyErrorAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.ERROR,
    error,
  };
}

function start(): FetchSurveyStartAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.START,
  };
}

const fetch: FetchActionCreators<TypeKeys.FETCH, Survey> = {
  receive,
  error,
  start,
};

export { fetch }

export type ActionTypes =
  | SetSurveyAction
  | SetSurveyResponsesAction
  | FetchSurveyReceiveAction
  | FetchSurveyErrorAction
  | FetchSurveyStartAction
