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
import { Credentials } from '../../models/credential';
import {
  fetch,
  ActionTypes,
  TypeKeys,
} from './actions';
import requestToMatchService from '../../services/request-to-match-service';

export interface State {
  readonly credentials: Credentials;
  readonly fetchState: FetchState;
}

const initialState: State = {
  fetchState: initialFetchState,
  credentials: Immutable.List(),
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.FETCH:
      return {
        ...state,
        fetchState: fetchStateReducer(action),
        credentials: getDataOrCur(action, state.credentials),
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action.type;
      return state;
  }
};

const fetchCredentialOptions: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = () => {
  return async (dispatch: Dispatch<State>) => {
    await dispatch(fetch.start());
    try {
      const data = await requestToMatchService.getAllCredentials();
      return dispatch(fetch.receive(data));
    } catch(e) {
      return dispatch(fetch.error(e.message));
    }
  };
}

export { fetchCredentialOptions };
