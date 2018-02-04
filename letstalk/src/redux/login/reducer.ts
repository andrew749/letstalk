import Immutable from 'immutable';
import { Action, ActionCreator, Dispatch } from 'redux';
import { ThunkAction } from 'redux-thunk';

import { ActionTypes, TypeKeys, fetchingLoginAction, receivedErrorMsgAction } from './actions';
import { InvalidCredentialsError } from '../../services/sessionService';
import auth from '../../services/auth';

export type State = {
  username: string;
  password: string;
  errorMsg: string | null;
  isFetching: boolean;
};

const initialState: State = {
  username: '',
  password: '',
  errorMsg: null,
  isFetching: false,
};

export const login: ActionCreator<ThunkAction<Promise<void>, State, void>> =
  (username: string, password: string) => {
  return async (dispatch: Dispatch<State>): Promise<void> => {
    dispatch(fetchingLoginAction());
    try {
      await auth.login(username, password);
    } catch (e) {
      if (e.message === InvalidCredentialsError.tpe) {
        // TODO: take string from message itself
        dispatch(receivedErrorMsgAction('Invalid username or password'));
      }
      throw e;
    }
  }
}

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.FETCHING_LOGIN:
      return {
        ...state,
        isFetching: true,
      };
    case TypeKeys.RECEIVED_ERROR_MSG:
      return {
        ...state,
        errorMsg: action.errorMsg,
        isFetching: false,
      };
    case TypeKeys.SET_USERNAME:
      return {
        ...state,
        username: action.username,
      };
    case TypeKeys.SET_PASSWORD:
      return {
        ...state,
        password: action.password,
      };
    case TypeKeys.RESET:
      return initialState;
    default:
      // Ensure exhaustiveness of select
      const _: never = action;
      return state;
  }
};
