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
  positionAdd,
  positionRemove,
  simpleTraitAdd,
  simpleTraitRemove,
  ActionTypes,
  TypeKeys,
} from './actions';
import profileService from '../../services/profile-service';
import requestToMatchService from '../../services/request-to-match-service';
import { AddUserPositionRequest } from '../../services/request-to-match-service';
import { UserPosition } from '../../models/position';
import { UserSimpleTrait } from '../../models/simple-trait';

export interface State {
  readonly profile?: ProfileData;
  readonly fetchState: FetchState;
}

const initialState: State = {
  fetchState: initialFetchState,
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  let profile: ProfileData = null;
  switch (action.type) {
    case TypeKeys.FETCH:
      return {
        ...state,
        fetchState: fetchStateReducer(action),
        profile: getDataOrCur(action, state.profile),
      };
    case TypeKeys.POSITION_ADD:
      profile = state.profile === null ? null : {
        ...state.profile,
        userPositions: state.profile.userPositions.push(action.position),
      }
      return {
        ...state,
        profile,
      }
    case TypeKeys.POSITION_REMOVE:
      profile = state.profile === null ? null : {
        ...state.profile,
        userPositions: state.profile.userPositions.filter(pos => pos.id !== action.id).toList(),
      }
      return {
        ...state,
        profile,
      }
    case TypeKeys.SIMPLE_TRAIT_ADD:
      console.log(action.simpleTrait);
      profile = state.profile === null ? null : {
        ...state.profile,
        userSimpleTraits: state.profile.userSimpleTraits.push(action.simpleTrait),
      }
      return {
        ...state,
        profile,
      }
    case TypeKeys.SIMPLE_TRAIT_REMOVE:
      profile = state.profile === null ? null : {
        ...state.profile,
        userSimpleTraits: state.profile.userSimpleTraits.filter(trait => {
          return trait.id !== action.id;
        }).toList(),
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

const addPosition: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (req: AddUserPositionRequest) => {
  return async (dispatch: Dispatch<State>) => {
    const position = await requestToMatchService.addUserPosition(req);
    return dispatch(positionAdd(position));
  };
}

const removePosition: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (id: number) => {
  return async (dispatch: Dispatch<State>) => {
    await requestToMatchService.removeUserPosition(id);
    return dispatch(positionRemove(id));
  };
}

const addSimpleTraitById: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (id: number) => {
  return async (dispatch: Dispatch<State>) => {
    const simpleTrait = await requestToMatchService.addUserSimpleTraitById(id);
    return dispatch(simpleTraitAdd(simpleTrait));
  };
}

const addSimpleTraitByName: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (name: string) => {
  return async (dispatch: Dispatch<State>) => {
    const simpleTrait = await requestToMatchService.addUserSimpleTraitByName(name);
    return dispatch(simpleTraitAdd(simpleTrait));
  };
}

const removeSimpleTrait: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = (id: number) => {
  return async (dispatch: Dispatch<State>) => {
    await requestToMatchService.removeUserSimpleTrait(id);
    return dispatch(simpleTraitRemove(id));
  };
}

export {
  fetchProfile,
  addPosition,
  removePosition,
  addSimpleTraitById,
  addSimpleTraitByName,
  removeSimpleTrait,
};
