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
  APPEND       = 'NOTIFICATIONS/APPEND',
}

export type Notifications = Immutable.List<Notification>;

type NotificationEditReceiveAction = FetchReceiveAction<TypeKeys.FETCH, Notifications>;
type NotificationEditErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type NotificationEditStartAction = FetchStartAction<TypeKeys.FETCH>;

interface NotificationUpdateStateAction extends Action {
  readonly type: TypeKeys.UPDATE_STATE;
  readonly notificationIds: Immutable.List<number>;
  readonly state: NotifState;
}

interface NotificationAppendAction extends Action {
  readonly type: TypeKeys.APPEND;
  readonly data: Notifications;
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
  notificationIds: Immutable.List<number>,
  state: NotifState
): NotificationUpdateStateAction {
  return {
    type: TypeKeys.UPDATE_STATE,
    notificationIds,
    state,
  };
}

export function appendNotifications(notifications: Notifications): NotificationAppendAction {
  return {
    type: TypeKeys.APPEND,
    data: notifications,
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
  | NotificationAppendAction
