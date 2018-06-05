import Immutable from 'immutable';
import { Action } from 'redux'

export enum TypeKeys {
  UPDATE_VALUE = 'SEARCH_BAR/UPDATE_VALUE',
}

export interface SearchBarUpdateValueAction extends Action {
  readonly type: TypeKeys.UPDATE_VALUE;
  readonly value: string;
}

export function updateSearchValue(value: string): SearchBarUpdateValueAction {
  return {
    type: TypeKeys.UPDATE_VALUE,
    value,
  };
}

export type ActionTypes =
  | SearchBarUpdateValueAction
