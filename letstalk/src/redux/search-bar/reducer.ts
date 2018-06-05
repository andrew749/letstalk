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
  ActionTypes,
  TypeKeys,
} from './actions';
import requestToMatchService from '../../services/request-to-match-service';

export interface State {
  readonly value: string;
}

const initialState: State = {
  value: '',
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.UPDATE_VALUE:
      console.log(action.value);
      return {
        ...state,
        value: action.value,
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action.type;
      return state;
  }
};

const updateValue: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (value: string) => {
  return async (dispatch: Dispatch<State>) => {
    return dispatch(updateSearchValue(value));
  };
}

export { updateValue };
