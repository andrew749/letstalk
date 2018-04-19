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
import { Credential, CredentialWithId } from '../../models/credential';
import {
  credentialAdd,
  credentialSetState,
  credentialRemove,
  CredentialStates,
  CredentialWithState,
  CredentialsWithState,
  fetch,
  ActionTypes,
  TypeKeys,
} from './actions';
import requestToMatchService from '../../services/request-to-match-service';

export interface State {
  readonly credentialsWithState: CredentialsWithState;
  readonly fetchState: FetchState;
}

const initialState: State = {
  fetchState: initialFetchState,
  credentialsWithState: Immutable.List(),
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  let newCredentialsWithState: CredentialsWithState = null;
  switch (action.type) {
    case TypeKeys.FETCH:
      return {
        ...state,
        fetchState: fetchStateReducer(action),
        credentialsWithState: getDataOrCur(action, state.credentialsWithState),
      };
    case TypeKeys.ADD_CREDENTIAL:
      const newCredentialWithState: CredentialWithState = {
        ...action.credentialWithId,
        state: 'normal',
      };
      newCredentialsWithState = state.credentialsWithState ?
        state.credentialsWithState.push(newCredentialWithState) :
        Immutable.List([newCredentialWithState]);

      return {
        ...state,
        credentialsWithState: newCredentialsWithState,
      };
    case TypeKeys.SET_STATE_CREDENTIAL:
      newCredentialsWithState = state.credentialsWithState.map(credentialWithState => {
        return credentialWithState.credentialId === action.credentialId ?
          { ...credentialWithState, state: action.state } : credentialWithState;
      }).toList();

      return {
        ...state,
        credentialsWithState: newCredentialsWithState,
      };
    case TypeKeys.REMOVE_CREDENTIAL:
      newCredentialsWithState = state.credentialsWithState.filter(credentialWithState => {
        return credentialWithState.credentialId !== action.credentialId;
      }).toList();

      return {
        ...state,
        credentialsWithState: newCredentialsWithState,
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action;
      return state;
  }
};

const fetchCredentials: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = () => {
  return async (dispatch: Dispatch<State>) => {
    dispatch(fetch.start());
    try {
      const data = await requestToMatchService.getCredentials();
      const credentialsWithState = data.map(credentialWithId => {
        return {
          ...credentialWithId,
          state: 'normal' as CredentialStates,
        };
      }).toList();
      return dispatch(fetch.receive(credentialsWithState));
    } catch(e) {
      return dispatch(fetch.error(e.message));
    }
  };
}

// TODO: Make your own ActionCreator type, since this never checks the types of the param.
const addCredential: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (credential: Credential) => {
  return async (dispatch: Dispatch<State>) => {
    const credentialId = await requestToMatchService.addCredential(credential);
    return dispatch(credentialAdd({ ...credential, credentialId }));
  };
}

const removeCredential: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (credentialId: number) => {
  return async (dispatch: Dispatch<State>) => {
    await dispatch(credentialSetState(credentialId, 'deleting'));
    await requestToMatchService.removeCredential(credentialId);
    return dispatch(credentialRemove(credentialId));
  };
}

export { addCredential, fetchCredentials, removeCredential };
