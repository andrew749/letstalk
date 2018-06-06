// TODO: If this pattern shows up again, maybe generalize the item editing code.
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
  credentialRequestAdd,
  credentialRequestSetState,
  credentialRequestRemove,
  CredentialRequestStates,
  CredentialWithState,
  CredentialsWithState,
  fetch,
  ActionTypes,
  TypeKeys,
} from './actions';
import requestToMatchService from '../../services/request-to-match-service';

export interface State {
  readonly credentialRequestsWithState: CredentialsWithState;
  readonly fetchState: FetchState;
}

const initialState: State = {
  fetchState: initialFetchState,
  credentialRequestsWithState: Immutable.List(),
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  let newCredentialRequestsWithState: CredentialsWithState = null;
  switch (action.type) {
    case TypeKeys.FETCH:
      return {
        ...state,
        fetchState: fetchStateReducer(action),
        credentialRequestsWithState: getDataOrCur(action, state.credentialRequestsWithState),
      };
    case TypeKeys.ADD_CREDENTIAL:
      const newCredentialWithState: CredentialWithState = {
        ...action.credential,
        state: 'normal',
      };
      newCredentialRequestsWithState =
        state.credentialRequestsWithState.push(newCredentialWithState);

      return {
        ...state,
        credentialRequestsWithState: newCredentialRequestsWithState,
      };
    case TypeKeys.SET_STATE_CREDENTIAL:
      newCredentialRequestsWithState = state.credentialRequestsWithState.map(
        credentialRequestWithState => {
          return credentialRequestWithState.id === action.credentialId ?
          { ...credentialRequestWithState, state: action.state } : credentialRequestWithState;
      }).toList();

      return {
        ...state,
        credentialRequestsWithState: newCredentialRequestsWithState,
      };
    case TypeKeys.REMOVE_CREDENTIAL:
      newCredentialRequestsWithState = state.credentialRequestsWithState.filter(
        credentialRequestWithState => {
        return credentialRequestWithState.id !== action.credentialId;
      }).toList();

      return {
        ...state,
        credentialRequestsWithState: newCredentialRequestsWithState,
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action;
      return state;
  }
};

const fetchCredentialRequests: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = () => {
  return async (dispatch: Dispatch<State>) => {
    dispatch(fetch.start());
    try {
      const data = await requestToMatchService.getCredentialRequests();
      const credentialsWithState = data.map(credentialRequest => {
        return {
          ...credentialRequest,
          state: 'normal' as CredentialRequestStates,
        };
      }).toList();
      return dispatch(fetch.receive(credentialsWithState));
    } catch(e) {
      return dispatch(fetch.error(e));
    }
  };
}

const addCredentialRequest: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (credential: Credential) => {
  return async (dispatch: Dispatch<State>) => {
    await requestToMatchService.addCredentialRequest(credential.id);
    return dispatch(credentialRequestAdd(credential));
  };
}

const removeCredentialRequest: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (credentialRequestId: number) => {
  return async (dispatch: Dispatch<State>) => {
    await dispatch(credentialRequestSetState(credentialRequestId, 'deleting'));
    await requestToMatchService.removeCredentialRequest(credentialRequestId);
    return dispatch(credentialRequestRemove(credentialRequestId));
  };
}

export { addCredentialRequest, fetchCredentialRequests, removeCredentialRequest };
