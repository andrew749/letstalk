import { ProfileData } from '../../models/profile';
import {
  FetchReceiveAction,
  FetchErrorAction,
  FetchStartAction,
  FetchTypeKeys,
  FetchActionCreators,
} from '../actions';

export enum TypeKeys {
  FETCH = 'PROFILE/FETCH',
}

type ProfileReceiveAction = FetchReceiveAction<TypeKeys.FETCH, ProfileData>;
type ProfileErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type ProfileStartAction = FetchStartAction<TypeKeys.FETCH>;

function receive(data: ProfileData): ProfileReceiveAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.RECEIVE,
    data,
  };
}

function error(errorMsg: string): ProfileErrorAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.ERROR,
    errorMsg,
  };
}

function start(): ProfileStartAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.START,
  };
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
