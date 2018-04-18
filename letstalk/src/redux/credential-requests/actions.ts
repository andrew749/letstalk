import Immutable from 'immutable';
import { Action } from 'redux'

import { CredentialRequest } from '../../models/credential';
import {
  FetchReceiveAction,
  FetchErrorAction,
  FetchStartAction,
  FetchTypeKeys,
  FetchActionCreators,
} from '../actions';

export enum TypeKeys {
  FETCH                = 'CREDENTIAL_REQUESTS/FETCH',
  ADD_CREDENTIAL       = 'CREDENTIAL_REQUESTS/ADD_CREDENTIAL',
  SET_STATE_CREDENTIAL = 'CREDENTIAL_REQUESTS/SET_STATE_CREDENTIAL',
  REMOVE_CREDENTIAL    = 'CREDENTIAL_REQUESTS/REMOVE_CREDENTIAL',
}

export type CredentialRequestStates = 'normal' | 'deleting';

export interface CredentialRequestWithState extends CredentialRequest {
  readonly state: CredentialRequestStates;
}

export type CredentialsWithState = Immutable.List<CredentialRequestWithState>;

type CredentialRequestEditReceiveAction = FetchReceiveAction<TypeKeys.FETCH, CredentialsWithState>;
type CredentialRequestEditErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type CredentialRequestEditStartAction = FetchStartAction<TypeKeys.FETCH>;

export interface CredentialRequestAddAction extends Action {
  readonly type: TypeKeys.ADD_CREDENTIAL;
  readonly credentialRequest: CredentialRequest;
}

export interface CredentialRequestSetStateAction extends Action {
  readonly type: TypeKeys.SET_STATE_CREDENTIAL;
  readonly credentialRequestId: number;
  readonly state: CredentialRequestStates;
}

export interface CredentialRequestRemoveAction extends Action {
  readonly type: TypeKeys.REMOVE_CREDENTIAL;
  readonly credentialRequestId: number;
}

function receive(data: CredentialsWithState): CredentialRequestEditReceiveAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.RECEIVE,
    data,
  };
}

function error(errorMsg: string): CredentialRequestEditErrorAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.ERROR,
    errorMsg,
  };
}

function start(): CredentialRequestEditStartAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.START,
  };
}

export function credentialRequestAdd(
  credentialRequest: CredentialRequest,
): CredentialRequestAddAction {
  return {
    type: TypeKeys.ADD_CREDENTIAL,
    credentialRequest,
  };
}

export function credentialRequestSetState(
  credentialRequestId: number,
  state: CredentialRequestStates,
): CredentialRequestSetStateAction {
  return {
    type: TypeKeys.SET_STATE_CREDENTIAL,
    credentialRequestId,
    state,
  };
}

export function credentialRequestRemove(
  credentialRequestId: number,
): CredentialRequestRemoveAction {
  return {
    type: TypeKeys.REMOVE_CREDENTIAL,
    credentialRequestId,
  };
}

const fetch: FetchActionCreators<TypeKeys.FETCH, CredentialsWithState> = {
  receive,
  error,
  start,
};

export { fetch }

export type ActionTypes =
  | CredentialRequestEditReceiveAction
  | CredentialRequestEditErrorAction
  | CredentialRequestEditStartAction
  | CredentialRequestAddAction
  | CredentialRequestRemoveAction
  | CredentialRequestSetStateAction
