import Immutable from 'immutable';
import { Action } from 'redux'

import { CredentialWithId } from '../../models/credential';
import {
  FetchReceiveAction,
  FetchErrorAction,
  FetchStartAction,
  FetchTypeKeys,
  FetchActionCreators,
} from '../actions';

export enum TypeKeys {
  FETCH            = 'CREDENTIALS/FETCH',
  ADD_CREDENTIAL   = 'CREDENTIALS/ADD_CREDENTIAL',
}

export type CredentialStates = 'normal' | 'deleting';

export interface CredentialWithState extends CredentialWithId {
  readonly state: CredentialStates;
}

export type CredentialsWithState = Immutable.List<CredentialWithState>;

type CredentialEditReceiveAction = FetchReceiveAction<TypeKeys.FETCH, CredentialsWithState>;
type CredentialEditErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type CredentialEditStartAction = FetchStartAction<TypeKeys.FETCH>;

export interface CredentialAddAction extends Action {
  readonly type: TypeKeys.ADD_CREDENTIAL;
  readonly credentialWithId: CredentialWithId;
}

function receive(data: CredentialsWithState): CredentialEditReceiveAction {
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

export function credentialAdd(credentialWithId: CredentialWithId): CredentialAddAction {
  return {
    type: TypeKeys.ADD_CREDENTIAL,
    credentialWithId,
  };
}

const fetch: FetchActionCreators<TypeKeys.FETCH, CredentialsWithState> = {
  receive,
  error,
  start,
};

export { fetch }

export type ActionTypes =
  | CredentialEditReceiveAction
  | CredentialEditErrorAction
  | CredentialEditStartAction
  | CredentialAddAction
