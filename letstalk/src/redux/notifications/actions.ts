import Immutable from 'immutable';
import { Action } from 'redux'

import { Notification, NotifState } from '../../models/notification';
import {
  FetchReceiveAction,
  FetchErrorAction,
  FetchStartAction,
  FetchTypeKeys,
  FetchActionCreators,
} from '../actions';
import { APIError } from '../../services/requests';

export enum TypeKeys {
  FETCH        = 'NOTIFICATIONS/FETCH',
  UPDATE_STATE = 'NOTIFICATIONS/UPDATE_STATE',
}

export type Notifications = Immutable.List<Notification>;

type NotificationEditReceiveAction = FetchReceiveAction<TypeKeys.FETCH, Notifications>;
type NotificationEditErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type NotificationEditStartAction = FetchStartAction<TypeKeys.FETCH>;

interface NotificationUpdateStateAction extends Action {
  readonly type: TypeKeys.UPDATE_STATE;
  readonly notificationId: number;
  readonly state: NotifState;
}

function receive(data: Notifications): NotificationEditReceiveAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.RECEIVE,
    data,
  };
}

function error(error: APIError): NotificationEditErrorAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.ERROR,
    error,
  };
}

function start(): NotificationEditStartAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.START,
  };
}

export function notificationUpdateState(
  notificationId: number,
  state: NotifState
): NotificationUpdateStateAction {
  return {
    type: TypeKeys.UPDATE_STATE,
    notificationId,
    state,
  };
}

const fetch: FetchActionCreators<TypeKeys.FETCH, Notifications> = {
  receive,
  error,
  start,
};

export { fetch }

export type ActionTypes =
  | NotificationEditReceiveAction
  | NotificationEditErrorAction
  | NotificationEditStartAction
  | NotificationUpdateStateAction
