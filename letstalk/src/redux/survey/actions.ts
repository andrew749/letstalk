import { Action } from 'redux'

import { OnboardingStatus } from '../../models';
import {Survey, SurveyQuestions, SurveyResponses} from "../../models/survey";

export enum TypeKeys {
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

export function setSurvey(
  survey: Survey
): SetSurveyAction {
  return {
    type: TypeKeys.SET_SURVEY,
    survey,
  };
}

export function setSurveyResponses(
  responses: SurveyResponses
): SetSurveyResponsesAction {
  return {
    type: TypeKeys.SET_SURVEY_RESPONSES,
    responses,
  };
}

export type ActionTypes =
  | SetSurveyAction
  | SetSurveyResponsesAction
