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
  Notifications,
  fetch,
  ActionTypes,
  TypeKeys,
} from './actions';
import profileService from '../../services/profile-service';

const NOTIF_BATCH_SIZE = 10;

export interface State {
  readonly notifications: Notifications;
  readonly fetchState: FetchState;
}

const initialState: State = {
  fetchState: initialFetchState,
  notifications: Immutable.List(),
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.FETCH:
      return {
        ...state,
        fetchState: fetchStateReducer(action),
        notifications: getDataOrCur(action, state.notifications),
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action.type;
      return state;
  }
};

const fetchNewestNotifications: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = () => {
  return async (dispatch: Dispatch<State>) => {
    dispatch(fetch.start());
    try {
      const notifications = await profileService.getNotifications(NOTIF_BATCH_SIZE);
      return dispatch(fetch.receive(notifications));
    } catch(e) {
      return dispatch(fetch.error(e));
    }
  };
}

export { fetchNewestNotifications };
