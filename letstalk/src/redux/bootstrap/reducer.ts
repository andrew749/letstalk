import Immutable from 'immutable';
import { ActionCreator, Dispatch } from 'redux'
import { ThunkAction } from 'redux-thunk';

import {
  FetchTypeKeys,
  FetchState,
  initialFetchState,
} from '../actions';
import {
  BootstrapData
} from '../../models/bootstrap';
import {
  fetch,
  ActionTypes,
  TypeKeys,
} from './actions';
import profileService from '../../services/profile-service';

export interface State {
  readonly bootstrap?: BootstrapData;
  readonly fetchState: FetchState;
}

const initialState: State = {
  fetchState: initialFetchState,
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.FETCH:
      switch (action.fetchType) {
        case FetchTypeKeys.RECEIVE:
          return {
            ...state,
            fetchState: { state: 'success' },
            bootstrap: action.data,
          };
        case FetchTypeKeys.ERROR:
          return {
            ...state,
            fetchState: { state: 'error', errorMsg: action.errorMsg },
          };
        case FetchTypeKeys.START:
          return {
            ...state,
            fetchState: { state: 'fetching' },
          };
        default:
          // Ensure exhaustiveness of select
          const _: never = action;
          return state;
        }
    default:
      // Ensure exhaustiveness of select
      const _: never = action.type;
      return state;
  }
};

const fetchBootstrap: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = () => {
  return async (dispatch: Dispatch<State>) => {
    dispatch(fetch.start());
    try {
      const data = await profileService.bootstrap();
      return fetch.receive(data);
    } catch(e) {
      return fetch.error(e.message);
    }
  };
}
export { fetchBootstrap };

