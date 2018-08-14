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
  appendNotifications,
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
  readonly hasAll: boolean;
}

const initialState: State = {
  fetchState: initialFetchState,
  notifications: Immutable.List(),
  hasAll: false,
};

function mergeNotifications(fst: Notifications, snd: Notifications): Notifications {
  return fst.concat(snd).groupBy(n => n.notificationId)
    .map(notifs => notifs.first())
    .sortBy(n => -n.notificationId).toList();
}

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.FETCH:
      const newState: State = {
        ...state,
        fetchState: fetchStateReducer(action),
        notifications: getDataOrCur(action, state.notifications),
      };
      if (action.fetchType === FetchTypeKeys.RECEIVE) {
        return {
          ...newState,
          hasAll: action.data.size < NOTIF_BATCH_SIZE,
        };
      }
      return newState;
    case TypeKeys.APPEND:
      return {
        ...state,
        notifications: mergeNotifications(state.notifications, action.data),
        hasAll: action.data.size < NOTIF_BATCH_SIZE,
      };
    case TypeKeys.UPDATE_STATE:
      const notifIdSet = Immutable.Set(action.notificationIds);
      return {
        ...state,
        notifications: state.notifications.map(notification => {
          if (notifIdSet.contains(notification.notificationId)) {
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

const fetchAdditionalNotifications: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (past: number) => {
  return async (dispatch: Dispatch<State>) => {
    const notifications = await profileService.getNotifications(NOTIF_BATCH_SIZE, past);
    return dispatch(appendNotifications(notifications));
  };
}

const updateNotificationState: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> =
  (notificationIds: Immutable.List<number>, state: NotifState) => {
  return async (dispatch: Dispatch<State>) => {
    await profileService.updateNotificationState(notificationIds, state);
    return dispatch(notificationUpdateState(notificationIds, state));
  };
}

export {
  fetchAdditionalNotifications,
  fetchNewestNotifications,
  updateNotificationState,
};
