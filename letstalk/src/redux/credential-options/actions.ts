import { CredentialOptions } from '../../models/credential';
import {
  FetchReceiveAction,
  FetchErrorAction,
  FetchStartAction,
  FetchTypeKeys,
  FetchActionCreators,
} from '../actions';

export enum TypeKeys {
  FETCH = 'CREDENTIAL_OPTIONS/FETCH',
}

type CredentialOptionsReceiveAction = FetchReceiveAction<TypeKeys.FETCH, CredentialOptions>;
type CredentialOptionsErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type CredentialOptionsStartAction = FetchStartAction<TypeKeys.FETCH>;

function receive(data: CredentialOptions): CredentialOptionsReceiveAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.RECEIVE,
    data,
  };
}

function error(errorMsg: string): CredentialOptionsErrorAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.ERROR,
    errorMsg,
  };
}

function start(): CredentialOptionsStartAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.START,
  };
}

const fetch: FetchActionCreators<TypeKeys.FETCH, CredentialOptions> = {
  receive,
  error,
  start,
};

export { fetch }

export type ActionTypes =
  | CredentialOptionsReceiveAction
  | CredentialOptionsErrorAction
  | CredentialOptionsStartAction
