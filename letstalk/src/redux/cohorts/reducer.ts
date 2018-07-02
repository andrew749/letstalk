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
import { Cohort } from '../../models/cohort';
import {
  fetch,
  ActionTypes,
  TypeKeys,
} from './actions';
import profileService from '../../services/profile-service';

export interface State {
  readonly cohorts: Immutable.List<Cohort>;
  readonly fetchState: FetchState;
}

const initialState: State = {
  fetchState: initialFetchState,
  cohorts: Immutable.List(),
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.FETCH:
      return {
        ...state,
        fetchState: fetchStateReducer(action),
        cohorts: getDataOrCur(action, state.cohorts),
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action.type;
      return state;
  }
};

const fetchCohorts: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = () => {
  return async (dispatch: Dispatch<State>) => {
    await dispatch(fetch.start());
    try {
      const data = await profileService.getAllCohorts();
      return dispatch(fetch.receive(data));
    } catch(e) {
      return dispatch(fetch.error(e));
    }
  };
}

export { fetchCohorts };
