import Immutable from 'immutable';
import { Action } from 'redux'

export const SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS = 'SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS';
export const SEARCH_LIST_TYPE_CREDENTIALS = 'SEARCH_LIST_TYPE_CREDENTIALS';

export type SearchListType =
  | 'SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS'
  | 'SEARCH_LIST_TYPE_CREDENTIALS'

export enum TypeKeys {
  UPDATE_VALUE = 'SEARCH_BAR/UPDATE_VALUE',
  UPDATE_FOCUS = 'SEARCH_BAR/UPDATE_FOCUS',
  UPDATE_SEARCH_LIST_TYPE = 'SEARCH_BAR/UPDATE_SEARCH_LIST_TYPE',
}

export interface SearchBarUpdateValueAction extends Action {
  readonly type: TypeKeys.UPDATE_VALUE;
  readonly value: string;
}

export interface SearchBarUpdateFocusAction extends Action {
  readonly type: TypeKeys.UPDATE_FOCUS;
  readonly hasFocus: boolean;
}

export interface SearchBarUpdateSearchListTypeAction extends Action {
  readonly type: TypeKeys.UPDATE_SEARCH_LIST_TYPE;
  readonly listType: SearchListType;
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

export function updateSearchListType(listType: SearchListType): SearchBarUpdateSearchListTypeAction {
  return {
    type: TypeKeys.UPDATE_SEARCH_LIST_TYPE,
    listType,
  };
}

export type ActionTypes =
  | SearchBarUpdateValueAction
  | SearchBarUpdateFocusAction
  | SearchBarUpdateSearchListTypeAction
