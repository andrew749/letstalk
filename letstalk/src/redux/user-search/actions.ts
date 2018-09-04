import { UserSearchResponse } from '../../models/user-search';
import {
  FetchReceiveAction,
  FetchErrorAction,
  FetchStartAction,
  FetchTypeKeys,
  FetchActionCreators,
} from '../actions';
import { APIError } from '../../services/requests';

export enum TypeKeys {
  FETCH = 'CREDENTIALS/FETCH',
}

type CredentialEditReceiveAction = FetchReceiveAction<TypeKeys.FETCH, UserSearchResponse>;
type CredentialEditErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type CredentialEditStartAction = FetchStartAction<TypeKeys.FETCH>;

function receive(data: UserSearchResponse): CredentialEditReceiveAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.RECEIVE,
    data,
  };
}

function error(error: APIError): CredentialEditErrorAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.ERROR,
    error,
  };
}

function start(): CredentialEditStartAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.START,
  };
}

const fetch: FetchActionCreators<TypeKeys.FETCH, UserSearchResponse> = {
  receive,
  error,
  start,
};

export { fetch }

export type ActionTypes =
  | CredentialEditReceiveAction
  | CredentialEditErrorAction
  | CredentialEditStartAction
