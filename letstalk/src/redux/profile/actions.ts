import { Action } from 'redux';

import { ProfileData } from '../../models/profile';
import {
  FetchReceiveAction,
  FetchErrorAction,
  FetchStartAction,
  FetchTypeKeys,
  FetchActionCreators,
} from '../actions';
import { APIError } from '../../services/requests';
import { Survey } from '../../models/survey';
import { UserPosition } from '../../models/position';
import { UserSimpleTrait } from '../../models/simple-trait';

export enum TypeKeys {
  FETCH = 'PROFILE/FETCH',
  POSITION_ADD = 'PROFILE/POSITION_ADD',
  POSITION_REMOVE = 'PROFILE/POSITIVE_REMOVE',
  SIMPLE_TRAIT_ADD = 'PROFILE/SIMPLE_TRAIT_ADD',
  SIMPLE_TRAIT_REMOVE = 'PROFILE/SIMPLE_TRAIT_REMOVE',
  SURVEY_SET = 'PROFILE/SURVEY_SET',
}

type ProfileReceiveAction = FetchReceiveAction<TypeKeys.FETCH, ProfileData>;
type ProfileErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type ProfileStartAction = FetchStartAction<TypeKeys.FETCH>;

export interface PositionAddAction extends Action {
  readonly type: TypeKeys.POSITION_ADD;
  readonly position: UserPosition;
}

export interface PositionRemoveAction extends Action {
  readonly type: TypeKeys.POSITION_REMOVE;
  readonly id: number;
}

export interface SimpleTraitAddAction extends Action {
  readonly type: TypeKeys.SIMPLE_TRAIT_ADD;
  readonly simpleTrait: UserSimpleTrait;
}

export interface SimpleTraitRemoveAction extends Action {
  readonly type: TypeKeys.SIMPLE_TRAIT_REMOVE;
  readonly id: number;
}

export interface SurveySetAction extends Action {
  readonly type: TypeKeys.SURVEY_SET;
  readonly survey: Survey;
}

function receive(data: ProfileData): ProfileReceiveAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.RECEIVE,
    data,
  };
}

function error(error: APIError): ProfileErrorAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.ERROR,
    error,
  };
}

function start(): ProfileStartAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.START,
  };
}

export function positionAdd(position: UserPosition): PositionAddAction {
  return { type: TypeKeys.POSITION_ADD, position };
}

export function positionRemove(id: number): PositionRemoveAction {
  return { type: TypeKeys.POSITION_REMOVE, id };
}

export function simpleTraitAdd(simpleTrait: UserSimpleTrait): SimpleTraitAddAction {
  return { type: TypeKeys.SIMPLE_TRAIT_ADD, simpleTrait };
}

export function simpleTraitRemove(id: number): SimpleTraitRemoveAction {
  return { type: TypeKeys.SIMPLE_TRAIT_REMOVE, id };
}

export function surveySet(survey: Survey): SurveySetAction {
  return { type: TypeKeys.SURVEY_SET, survey };
}

const fetch: FetchActionCreators<TypeKeys.FETCH, ProfileData> = {
  receive,
  error,
  start,
};

export { fetch }

export type ActionTypes =
  | ProfileReceiveAction
  | ProfileErrorAction
  | ProfileStartAction
  | PositionAddAction
  | PositionRemoveAction
  | SimpleTraitAddAction
  | SimpleTraitRemoveAction
  | SurveySetAction
