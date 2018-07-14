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
  notificationUpdateState,
  Notifications,
  fetch,
  ActionTypes,
  TypeKeys,
} from './actions';
import profileService from '../../services/profile-service';
import { NotifState } from '../../models/notification';

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
    case TypeKeys.UPDATE_STATE:
      return {
        ...state,
        notifications: state.notifications.map(notification => {
          if (notification.notificationId === action.notificationId) {
            return { ...notification, state: action.state };
          } else {
            return notification;
          }
        }).toList(),
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action;
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

const updateNotificationState: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (notificationId: number, state: NotifState) => {
  return async (dispatch: Dispatch<State>) => {
    await profileService.updateNotificationState(notificationId, state);
    return dispatch(notificationUpdateState(notificationId, state));
  };
}

export { fetchNewestNotifications, updateNotificationState };
