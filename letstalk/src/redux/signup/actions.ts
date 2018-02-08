import { Action, ActionCreator } from 'redux';

import {
  constActionCreator,
  setterActionCreator,
  SetterAction,
} from '../actions';

export enum TypeKeys {
  RECEIVED_ERROR_MSG  = 'SIGNUP/RECEIVED_ERROR_MSG',
  SET_FIRST_NAME      = 'SIGNUP/SET_FIRST_NAME',
  SET_LAST_NAME       = 'SIGNUP/SET_LAST_NAME',
  SET_EMAIL           = 'SIGNUP/SET_EMAIL',
  SET_PASSWORD        = 'SIGNUP/SET_PASSWORD',
  SET_PHONE_NUMBER    = 'SIGNUP/SET_PHONE_NUMBER',
  RESET               = 'SIGNUP/RESET',
}

export interface ReceivedErrorMsgAction extends SetterAction<string> {
  readonly type: TypeKeys.RECEIVED_ERROR_MSG,
}

export interface SetFirstNameAction extends SetterAction<string> {
  readonly type: TypeKeys.SET_FIRST_NAME,
}

export interface SetLastNameAction extends SetterAction<string> {
  readonly type: TypeKeys.SET_LAST_NAME,
}

export interface SetEmailAction extends SetterAction<string> {
  readonly type: TypeKeys.SET_EMAIL,
}

export interface SetPasswordAction extends SetterAction<string> {
  readonly type: TypeKeys.SET_PASSWORD,
}

export interface SetPhoneNumberAction extends SetterAction<string> {
  readonly type: TypeKeys.SET_PHONE_NUMBER,
}

export interface ResetAction extends Action {
  readonly type: TypeKeys.RESET,
}

export const receivedErrorMsgAction = setterActionCreator<string>(TypeKeys.RECEIVED_ERROR_MSG);
export const setEmailAction         = setterActionCreator<string>(TypeKeys.SET_EMAIL);
export const setPasswordAction      = setterActionCreator<string>(TypeKeys.SET_PASSWORD);
export const setFirstNameAction     = setterActionCreator<string>(TypeKeys.SET_FIRST_NAME);
export const setLastNameAction      = setterActionCreator<string>(TypeKeys.SET_LAST_NAME);
export const setPhoneNumberAction   = setterActionCreator<string>(TypeKeys.SET_PHONE_NUMBER);

export const resetAction = constActionCreator(TypeKeys.RESET);

export type ActionTypes =
  | ReceivedErrorMsgAction
  | SetFirstNameAction
  | SetLastNameAction
  | SetEmailAction
  | SetPasswordAction
  | SetPhoneNumberAction
  | ResetAction
