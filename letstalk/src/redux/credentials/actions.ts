import Immutable from 'immutable';
import { Action } from 'redux'

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
  FETCH                = 'CREDENTIALS/FETCH',
  ADD_CREDENTIAL       = 'CREDENTIALS/ADD_CREDENTIAL',
  SET_STATE_CREDENTIAL = 'CREDENTIALS/SET_STATE_CREDENTIAL',
  REMOVE_CREDENTIAL    = 'CREDENTIALS/REMOVE_CREDENTIAL',
}

export type CredentialStates = 'normal' | 'deleting';

export interface CredentialWithState extends Credential {
  readonly state: CredentialStates;
}

export type CredentialsWithState = Immutable.List<CredentialWithState>;

type CredentialEditReceiveAction = FetchReceiveAction<TypeKeys.FETCH, CredentialsWithState>;
type CredentialEditErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type CredentialEditStartAction = FetchStartAction<TypeKeys.FETCH>;

export interface CredentialAddAction extends Action {
  readonly type: TypeKeys.ADD_CREDENTIAL;
  readonly credential: Credential;
}

export interface CredentialSetStateAction extends Action {
  readonly type: TypeKeys.SET_STATE_CREDENTIAL;
  readonly credentialId: number;
  readonly state: CredentialStates;
}

export interface CredentialRemoveAction extends Action {
  readonly type: TypeKeys.REMOVE_CREDENTIAL;
  readonly credentialId: number;
}

function receive(data: CredentialsWithState): CredentialEditReceiveAction {
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

export function credentialAdd(credential: Credential): CredentialAddAction {
  return {
    type: TypeKeys.ADD_CREDENTIAL,
    credential,
  };
}

export function credentialSetState(
  credentialId: number,
  state: CredentialStates,
): CredentialSetStateAction {
  return {
    type: TypeKeys.SET_STATE_CREDENTIAL,
    credentialId,
    state,
  };
}

export function credentialRemove(credentialId: number): CredentialRemoveAction {
  return {
    type: TypeKeys.REMOVE_CREDENTIAL,
    credentialId,
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
  | CredentialRemoveAction
  | CredentialSetStateAction
