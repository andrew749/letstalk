import { Action } from 'redux'

import { OnboardingStatus } from '../../models';
import {Survey, SurveyQuestions, SurveyResponses} from "../../models/survey";

export enum TypeKeys {
  SET_STATE = 'SURVEY/SET_STATE',
  SET_QUESTIONS = 'SURVEY/SET_QUESTIONS',
  SET_RESPONSES = 'SURVEY/SET_RESPONSES',
}

export interface SetSurveyStateAction extends Action {
  readonly type: TypeKeys.SET_STATE;
  readonly survey: Survey;
}

export interface SetSurveyQuestionsAction extends Action {
  readonly type: TypeKeys.SET_QUESTIONS;
  readonly questions: SurveyQuestions;
}

export interface SetSurveyResponsesAction extends Action {
  readonly type: TypeKeys.SET_RESPONSES;
  readonly responses: SurveyResponses;
}

export function setSurveyStateAction(
  survey: Survey
): SetSurveyStateAction {
  return {
    type: TypeKeys.SET_STATE,
    survey,
  };
}

export function setSurveyQuestionsAction(
  questions: SurveyQuestions
): SetSurveyQuestionsAction {
  return {
    type: TypeKeys.SET_QUESTIONS,
    questions,
  };
}

export function setSurveyResponsesAction(
  responses: SurveyResponses
): SetSurveyResponsesAction {
  return {
    type: TypeKeys.SET_RESPONSES,
    responses,
  };
}

export type ActionTypes =
  | SetSurveyStateAction
  | SetSurveyQuestionsAction
  | SetSurveyResponsesAction
