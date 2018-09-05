import Immutable from 'immutable';
import { Action } from 'redux'

import { MultiTrait } from '../../models/multi-trait';

export const SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS = 'SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS';
export const SEARCH_LIST_TYPE_CREDENTIALS = 'SEARCH_LIST_TYPE_CREDENTIALS';

export type SearchListType =
  | 'SEARCH_LIST_TYPE_CREDENTIAL_REQUESTS'
  | 'SEARCH_LIST_TYPE_CREDENTIALS'

export enum TypeKeys {
  UPDATE_VALUE = 'SEARCH_BAR/UPDATE_VALUE',
  UPDATE_FOCUS = 'SEARCH_BAR/UPDATE_FOCUS',
  UPDATE_SUGGESTIONS = 'SEARCH_BAR/UPDATE_SUGGESTIONS',
}

export interface SearchBarUpdateValueAction extends Action {
  readonly type: TypeKeys.UPDATE_VALUE;
  readonly value: string;
}

export interface SearchBarUpdateFocusAction extends Action {
  readonly type: TypeKeys.UPDATE_FOCUS;
  readonly hasFocus: boolean;
}

export interface SearchBarUpdateSuggestionsAction extends Action {
  readonly type: TypeKeys.UPDATE_SUGGESTIONS;
  readonly suggestions: Immutable.List<MultiTrait>;
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

export function updateSearchSuggestions(
  suggestions: Immutable.List<MultiTrait>,
): SearchBarUpdateSuggestionsAction {
  return {
    type: TypeKeys.UPDATE_SUGGESTIONS,
    suggestions,
  };
}

export type ActionTypes =
  | SearchBarUpdateValueAction
  | SearchBarUpdateFocusAction
  | SearchBarUpdateSuggestionsAction
