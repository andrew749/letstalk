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
import { ProfileData } from '../../models/profile';
import {
  fetch,
  ActionTypes,
  TypeKeys,
} from './actions';
import profileService from '../../services/profile-service';

export interface State {
  readonly profile?: ProfileData;
  readonly fetchState: FetchState;
}

const initialState: State = {
  fetchState: initialFetchState,
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.FETCH:
      return {
        ...state,
        fetchState: fetchStateReducer(action),
        profile: getDataOrCur(action, state.profile),
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action.type;
      return state;
  }
};

const fetchMatchProfile: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (userId: number) => {
  return async (dispatch: Dispatch<State>) => {
    await dispatch(fetch.start());
    try {
      const data = await profileService.matchProfile(userId);
      return dispatch(fetch.receive(data));
    } catch(e) {
      return dispatch(fetch.error(e));
    }
  };
}

export { fetchMatchProfile };
