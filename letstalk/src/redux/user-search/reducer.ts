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
  fetch,
  ActionTypes,
  TypeKeys,
  Query,
  QueryTypes,
  setCurrentQuery,
} from './actions';
import { UserSearchResponse } from '../../models/user-search';
import userSearchService from '../../services/user-search-service';

export interface State {
  readonly currentQuery: Query
  readonly response: UserSearchResponse | null;
  readonly fetchState: FetchState;
}

const initialState: State = {
  currentQuery: { type: QueryTypes.YOUR_COHORT },
  fetchState: initialFetchState,
  response: null,
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.FETCH:
      return {
        ...state,
        fetchState: fetchStateReducer(action),
        response: getDataOrCur(action, state.response),
      };
    case TypeKeys.SET_CURRENT_QUERY:
      return {
        ...state,
        currentQuery: action.query,
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action;
      return state;
  }
};

const searchByCohort: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (cohortId: number, size: number) => {
  return async (dispatch: Dispatch<State>) => {
    dispatch(fetch.start());
    try {
      const data = await userSearchService.searchByCohort({ cohortId, size });
      return dispatch(fetch.receive(data));
    } catch(e) {
      return dispatch(fetch.error(e));
    }
  };
}

const searchByMyCohort: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (size: number) => {
  return async (dispatch: Dispatch<State>) => {
    dispatch(fetch.start());
    try {
      const data = await userSearchService.searchByMyCohort({ size });
      return dispatch(fetch.receive(data));
    } catch(e) {
      return dispatch(fetch.error(e));
    }
  };
}

const searchByPosition: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> =
  (roleId: number, organizationId: number, size: number) => {
    return async (dispatch: Dispatch<State>) => {
      dispatch(fetch.start());
      try {
        const data = await userSearchService.searchByPosition({ roleId, organizationId, size });
        return dispatch(fetch.receive(data));
      } catch(e) {
        return dispatch(fetch.error(e));
      }
    };
  }

const searchBySimpleTrait: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (simpleTraitId: number, size: number) => {
  return async (dispatch: Dispatch<State>) => {
    dispatch(fetch.start());
    try {
      const data = await userSearchService.searchBySimpleTrait({ simpleTraitId, size });
      return dispatch(fetch.receive(data));
    } catch(e) {
      return dispatch(fetch.error(e));
    }
  };
}

const setQuery: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (query: Query) => {
  return async (dispatch: Dispatch<State>) => {
    return dispatch(setCurrentQuery(query));
  };
}

export {
  searchByCohort,
  searchByMyCohort,
  searchByPosition,
  searchBySimpleTrait,
  setQuery,
};
