import Immutable from 'immutable';

import { Credential } from '../../models/credential';
import {
  FetchReceiveAction,
  FetchErrorAction,
  FetchStartAction,
  FetchTypeKeys,
  FetchActionCreators,
} from '../actions';

export enum TypeKeys {
  FETCH = 'CREDENTIALS/FETCH',
}

type Credentials = Immutable.List<Credential>;

type CredentialEditReceiveAction = FetchReceiveAction<TypeKeys.FETCH, Credentials>;
type CredentialEditErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type CredentialEditStartAction = FetchStartAction<TypeKeys.FETCH>;

function receive(data: Credentials): CredentialEditReceiveAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.RECEIVE,
    data,
  };
}

function error(errorMsg: string): CredentialEditErrorAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.ERROR,
    errorMsg,
  };
}

function start(): CredentialEditStartAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.START,
  };
}

const fetch: FetchActionCreators<TypeKeys.FETCH, Credentials> = {
  receive,
  error,
  start,
};

export { fetch }

export type ActionTypes =
  | CredentialEditReceiveAction
  | CredentialEditErrorAction
  | CredentialEditStartAction
