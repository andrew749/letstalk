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
import { Credential } from '../../models/credential';
import {
  updateSearchValue,
  updateSearchFocus,
  ActionTypes,
  TypeKeys,
} from './actions';
import requestToMatchService from '../../services/request-to-match-service';

export interface State {
  readonly value: string;
  readonly hasFocus: boolean;
}

const initialState: State = {
  value: '',
  hasFocus: false,
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  console.log(action);
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

export { updateFocus, updateValue };
