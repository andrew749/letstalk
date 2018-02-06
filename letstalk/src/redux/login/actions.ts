import Immutable from 'immutable';
import { Action, ActionCreator } from 'redux';

export enum TypeKeys {
  FETCHING_LOGIN      = 'LOGIN/FETCHING_LOGIN',
  RECEIVED_ERROR_MSG  = 'LOGIN/RECEIVED_ERROR_MSG',
  SET_USERNAME        = 'LOGIN/SET_USERNAME',
  SET_PASSWORD        = 'LOGIN/SET_PASSWORD',
  RESET               = 'LOGIN/RESET',
}

export interface FetchingLoginAction extends Action {
  readonly type: TypeKeys.FETCHING_LOGIN;
};

export interface ReceivedErrorMsgAction extends Action {
  readonly type: TypeKeys.RECEIVED_ERROR_MSG,
  readonly errorMsg: string,
}

export interface SetUsernameAction extends Action {
  readonly type: TypeKeys.SET_USERNAME,
  readonly username: string,
}

export interface SetPasswordAction extends Action {
  readonly type: TypeKeys.SET_PASSWORD,
  readonly password: string,
}

export interface ResetAction extends Action {
  readonly type: TypeKeys.RESET,
}

export const fetchingLoginAction: ActionCreator<FetchingLoginAction> =
  () => ({ type: TypeKeys.FETCHING_LOGIN });

export const receivedErrorMsgAction: ActionCreator<ReceivedErrorMsgAction> =
  (errorMsg: string) => ({
  type: TypeKeys.RECEIVED_ERROR_MSG,
  errorMsg,
});

export const setUsernameAction: ActionCreator<SetUsernameAction> =
  (username: string) => ({
  type: TypeKeys.SET_USERNAME,
  username,
});

export const setPasswordAction: ActionCreator<SetPasswordAction> =
  (password: string) => ({
  type: TypeKeys.SET_PASSWORD,
  password,
});

export const resetAction: ActionCreator<ResetAction> =
  () => ({ type: TypeKeys.RESET });

export type ActionTypes =
  | FetchingLoginAction
  | ReceivedErrorMsgAction
  | SetUsernameAction
  | SetPasswordAction
  | ResetAction
