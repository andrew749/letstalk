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
import {
  BootstrapData
} from '../../models/bootstrap';
import { setOnboardingStatusAction } from '../onboarding/actions';
import {
  fetch,
  removeRtmMatchesAction,
  ActionTypes,
  TypeKeys,
} from './actions';
import {
  USER_TYPE_ASKER,
  USER_TYPE_ANSWERER,
} from '../../models/user';
import profileService from '../../services/profile-service';
import requestToMatchService from '../../services/request-to-match-service';

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
      return {
        ...state,
        fetchState: fetchStateReducer(action),
        bootstrap: getDataOrCur(action, state.bootstrap),
      };
    case TypeKeys.REMOVE_RTM_MATCHES:
      if (state.bootstrap === null) return state;
      const relationships = state.bootstrap.relationships.filterNot(relationship => {
        return (relationship.userType === USER_TYPE_ASKER ||
          relationship.userType === USER_TYPE_ANSWERER) &&
          relationship.userId === action.userId;
      }).toList();
      return {
        ...state,
        bootstrap: {
          ...state.bootstrap,
          relationships
        },
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action;
      return state;
  }
};

const fetchBootstrap: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = () => {
  return async (dispatch: Dispatch<State>) => {
    await dispatch(fetch.start());
    try {
      const data = await profileService.bootstrap();
      await dispatch(setOnboardingStatusAction(data.onboardingStatus));
      return dispatch(fetch.receive(data));
    } catch(e) {
      console.log(e);
      return dispatch(fetch.error(e));
    }
  };
}

const removeRtmMatches: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (userId: number) => {
  return async (dispatch: Dispatch<State>) => {
    await requestToMatchService.removeRtmMatches(userId);
    return dispatch(removeRtmMatchesAction(userId));
  }
}

export { fetchBootstrap, removeRtmMatches };
