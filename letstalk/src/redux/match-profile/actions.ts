import { MatchProfileData } from '../../models/profile';
import {
  FetchReceiveAction,
  FetchErrorAction,
  FetchStartAction,
  FetchTypeKeys,
  FetchActionCreators,
} from '../actions';
import { APIError } from '../../services/requests';

export enum TypeKeys {
  FETCH = 'MATCH_PROFILE/FETCH',
}

type ProfileReceiveAction = FetchReceiveAction<TypeKeys.FETCH, MatchProfileData>;
type ProfileErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type ProfileStartAction = FetchStartAction<TypeKeys.FETCH>;

function receive(data: MatchProfileData): ProfileReceiveAction {
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

const fetch: FetchActionCreators<TypeKeys.FETCH, MatchProfileData> = {
  receive,
  error,
  start,
};

export { fetch }

export type ActionTypes =
  | ProfileReceiveAction
  | ProfileErrorAction
  | ProfileStartAction
