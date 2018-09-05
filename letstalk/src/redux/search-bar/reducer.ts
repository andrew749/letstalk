import Immutable from 'immutable';
import { ActionCreator, Dispatch } from 'redux'
import { ThunkAction } from 'redux-thunk';

import {
  fetchStateReducer,
  FetchTypeKeys,
  FetchState,
  getDataOrCur,
  initialFetchState,
} from '../actions';
import { MultiTrait } from '../../models/multi-trait';
import {
  updateSearchValue,
  updateSearchFocus,
  updateSearchSuggestions,
  ActionTypes,
  SearchListType,
  TypeKeys,
} from './actions';
import requestToMatchService from '../../services/request-to-match-service';

export interface State {
  readonly value: string;
  readonly hasFocus: boolean;
  readonly suggestions: Immutable.List<MultiTrait>;
}

const initialState: State = {
  value: '',
  hasFocus: false,
  suggestions: Immutable.List<MultiTrait>(),
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.UPDATE_VALUE:
      return {
        ...state,
        value: action.value,
      };
    case TypeKeys.UPDATE_FOCUS:
      return {
        ...state,
        hasFocus: action.hasFocus,
      };
    case TypeKeys.UPDATE_SUGGESTIONS:
      return {
        ...state,
        suggestions: action.suggestions,
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action;
      return state;
  }
};

const updateValue: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (value: string) => {
  return async (dispatch: Dispatch<State>) => {
    return dispatch(updateSearchValue(value));
  };
}

const updateFocus: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (hasFocus: boolean) => {
  return async (dispatch: Dispatch<State>) => {
    return dispatch(updateSearchFocus(hasFocus));
  };
}

const updateSuggestions: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (suggestions: Immutable.List<MultiTrait>) => {
  return async (dispatch: Dispatch<State>) => {
    return dispatch(updateSearchSuggestions(suggestions));
  };
}

export { updateFocus, updateValue, updateSuggestions };
