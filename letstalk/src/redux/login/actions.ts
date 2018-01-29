import Immutable from 'immutable';
import { Action, ActionCreator } from 'redux';

import { SessionToken } from '../../services/sessionService';

export enum TypeKeys {
  RECEIVE_LOGIN = 'RECEIVE_LOGIN',
  RECEIVE_LOGOUT = 'RECEIVE_LOGOUT',
}

export interface ReceiveLoginAction extends Action {
  readonly type: TypeKeys.RECEIVE_LOGIN;
};

export interface ReceiveLogoutAction extends Action {
  readonly type: TypeKeys.RECEIVE_LOGOUT;
};

export const receiveLoginAction: ActionCreator<ReceiveLoginAction> =
  (sessionToken: SessionToken) => {
  return {
    type: TypeKeys.RECEIVE_LOGIN,
    sessionToken,
  };
};

export const receiveLogoutAction: ActionCreator<ReceiveLogoutAction> = () => {
  return {
    type: TypeKeys.RECEIVE_LOGOUT,
  };
};

export type ActionTypes =
  | ReceiveLoginAction
  | ReceiveLogoutAction
