import Immutable from 'immutable';
import { Action } from 'redux'

export enum TypeKeys {
  UPDATE_VALUE = 'SEARCH_BAR/UPDATE_VALUE',
  UPDATE_FOCUS = 'SEARCH_BAR/UPDATE_FOCUS',
}

export interface SearchBarUpdateValueAction extends Action {
  readonly type: TypeKeys.UPDATE_VALUE;
  readonly value: string;
}

export interface SearchBarUpdateFocusAction extends Action {
  readonly type: TypeKeys.UPDATE_FOCUS;
  readonly hasFocus: boolean;
}

export function updateSearchValue(value: string): SearchBarUpdateValueAction {
  return {
    type: TypeKeys.UPDATE_VALUE,
    value,
  };
}

export function updateSearchFocus(hasFocus: boolean): SearchBarUpdateFocusAction {
  return {
    type: TypeKeys.UPDATE_FOCUS,
    hasFocus,
  };
}

export type ActionTypes =
  | SearchBarUpdateValueAction
  | SearchBarUpdateFocusAction
