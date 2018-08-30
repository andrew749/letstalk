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
  positionRemove,
  ActionTypes,
  TypeKeys,
} from './actions';
import profileService from '../../services/profile-service';
import requestToMatchService from '../../services/request-to-match-service';

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
    case TypeKeys.POSITION_REMOVE:
      const profile = state.profile === null ? null : {
        ...state.profile,
        userPositions: state.profile.userPositions.filter(pos => pos.id !== action.id).toList(),
      }
      return {
        ...state,
        profile,
      }
    default:
      // Ensure exhaustiveness of select
      const _: never = action;
      return state;
  }
};

const fetchProfile: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = () => {
  return async (dispatch: Dispatch<State>) => {
    await dispatch(fetch.start());
    try {
      const data = await profileService.me();
      return dispatch(fetch.receive(data));
    } catch(e) {
      return dispatch(fetch.error(e));
    }
  };
}

const removePosition: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (id: number) => {
  return async (dispatch: Dispatch<State>) => {
    await requestToMatchService.removeUserPosition(id);
    return dispatch(positionRemove(id));
  };
}

export {
  fetchProfile,
  removePosition,
};
