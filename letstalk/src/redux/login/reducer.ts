import Immutable from 'immutable';
import { Action } from 'redux';

import { ActionTypes, TypeKeys } from './actions';

export type State = {
  username: string;
  password: string;
  error: string | null;
  isFetching: boolean;
};

const initialState: State = {
  username: '',
  password: '',
  error: null,
  isFetching: false,
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.RECEIVE_LOGIN:
    case TypeKeys.RECEIVE_LOGOUT:
      return {
        ...state,
        error: null,
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action;
      return state;
  }
};
