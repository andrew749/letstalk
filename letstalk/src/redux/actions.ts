import { Action, ActionCreator } from 'redux';

export interface SetterAction<T> extends Action {
  readonly value: T,
}

export interface ConstAction extends Action {}

export const setterActionCreator = <T>(type: string): ActionCreator<SetterAction<T>> => {
  return (value: T) => ({
    type,
    value,
  });
};

export const constActionCreator = (type: string): ActionCreator<ConstAction> => {
  return () => ({ type });
}
