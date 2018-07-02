import Immutable from 'immutable';

import { Cohort } from '../../models/cohort';
import {
  FetchReceiveAction,
  FetchErrorAction,
  FetchStartAction,
  FetchTypeKeys,
  FetchActionCreators,
} from '../actions';
import { APIError } from '../../services/requests';

export enum TypeKeys {
  FETCH = 'COHORTS/FETCH',
}

type Cohorts = Immutable.List<Cohort>;

type CredentialOptionsReceiveAction = FetchReceiveAction<TypeKeys.FETCH, Cohorts>;
type CredentialOptionsErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type CredentialOptionsStartAction = FetchStartAction<TypeKeys.FETCH>;

function receive(data: Cohorts): CredentialOptionsReceiveAction {
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

const fetch: FetchActionCreators<TypeKeys.FETCH, Cohorts> = {
  receive,
  error,
  start,
};

export { fetch }

export type ActionTypes =
  | CredentialOptionsReceiveAction
  | CredentialOptionsErrorAction
  | CredentialOptionsStartAction
