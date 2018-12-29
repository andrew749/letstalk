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
import {
  fetch,
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
    default:
      // Ensure exhaustiveness of select
      const _: never = action.type;
      return state;
  }
};

const fetchBootstrap: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = () => {
  return async (dispatch: Dispatch<State>) => {
    await dispatch(fetch.start());
    try {
      const data = await profileService.bootstrap();
      return dispatch(fetch.receive(data));
    } catch(e) {
      return dispatch(fetch.error(e));
    }
  };
}

export { fetchBootstrap };
