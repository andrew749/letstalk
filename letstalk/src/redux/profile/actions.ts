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

export enum TypeKeys {
  FETCH = 'PROFILE/FETCH',
  POSITION_REMOVE = 'PROFILE/POSITIVE_REMOVE',
}

type ProfileReceiveAction = FetchReceiveAction<TypeKeys.FETCH, ProfileData>;
type ProfileErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type ProfileStartAction = FetchStartAction<TypeKeys.FETCH>;

export interface PositionRemoveAction extends Action {
  readonly type: TypeKeys.POSITION_REMOVE;
  readonly id: number;
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

export function positionRemove(id: number): PositionRemoveAction {
  return { type: TypeKeys.POSITION_REMOVE, id };
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
  | PositionRemoveAction
