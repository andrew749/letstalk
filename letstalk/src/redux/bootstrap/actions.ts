import { BootstrapData } from '../../models/bootstrap';
import {
  FetchReceiveAction,
  FetchErrorAction,
  FetchStartAction,
  FetchTypeKeys,
  FetchActionCreators,
} from '../actions';
import { APIError } from '../../services/requests';

export enum TypeKeys {
  FETCH              = 'BOOTSTRAP/FETCH',
}

type BootstrapReceiveAction = FetchReceiveAction<TypeKeys.FETCH, BootstrapData>;
type BootstrapErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type BootstrapStartAction = FetchStartAction<TypeKeys.FETCH>;

function receive(data: BootstrapData): BootstrapReceiveAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.RECEIVE,
    data,
  };
}

function error(error: APIError): BootstrapErrorAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.ERROR,
    error,
  };
}

function start(): BootstrapStartAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.START,
  };
}

const fetch: FetchActionCreators<TypeKeys.FETCH, BootstrapData> = {
  receive,
  error,
  start,
};

export { fetch }

export type ActionTypes =
  | BootstrapReceiveAction
  | BootstrapErrorAction
  | BootstrapStartAction
