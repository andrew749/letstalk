import Immutable from 'immutable';

import { Credential } from '../../models/credential';
import {
  FetchReceiveAction,
  FetchErrorAction,
  FetchStartAction,
  FetchTypeKeys,
  FetchActionCreators,
} from '../actions';
import { APIError } from '../../services/requests';

export enum TypeKeys {
  FETCH = 'CREDENTIAL_OPTIONS/FETCH',
}

type Credentials = Immutable.List<Credential>;

type CredentialOptionsReceiveAction = FetchReceiveAction<TypeKeys.FETCH, Credentials>;
type CredentialOptionsErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type CredentialOptionsStartAction = FetchStartAction<TypeKeys.FETCH>;

function receive(data: Credentials): CredentialOptionsReceiveAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.RECEIVE,
    data,
  };
}

function error(error: APIError): CredentialOptionsErrorAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.ERROR,
    error,
  };
}

function start(): CredentialOptionsStartAction {
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
  | CredentialOptionsReceiveAction
  | CredentialOptionsErrorAction
  | CredentialOptionsStartAction
