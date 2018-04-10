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
  CredentialElement,
  CredentialOptions,
  CredentialOrganization,
  CredentialOrganizationType,
  CredentialPosition,
  CredentialPositionType,
} from '../../models/credential';
import {
  fetch,
  ActionTypes,
  TypeKeys,
} from './actions';
import requestToMatchService from '../../services/request-to-match-service';

export interface State {
  readonly credentialOptions?: CredentialOptions;
  readonly credentialElements?: Immutable.List<CredentialElement>;
  readonly fetchState: FetchState;
}

function buildCredentialElements(options: CredentialOptions): Immutable.List<CredentialElement> {
  function reduceOrganizations(
    map: Immutable.Map<CredentialOrganizationType, Immutable.List<CredentialOrganization>>,
    org: CredentialOrganization,
  ) {
    return map.update(org.type, Immutable.List(), list => list.push(org))
  }
  function reducePositions(
    map: Immutable.Map<CredentialPositionType, Immutable.List<CredentialPosition>>,
    pos: CredentialPosition,
  ) {
    return map.update(pos.type, Immutable.List(), list => list.push(pos))
  }
  const organizationsByType = options.organizations.reduce(reduceOrganizations, Immutable.Map());
  const positionsByType = options.positions.reduce(reducePositions, Immutable.Map());
  const elements = options.validPairs.flatMap(pair => {
    return organizationsByType.get(pair.organizationType).flatMap(org => {
      return positionsByType.get(pair.positionType).map(pos => {
        return {
          id: [pos.id, org.id] as [number, number],
          value: `${pos.name} at ${org.name}`,
        }
      });
    });
  });
  return elements.toList();
}

const initialState: State = {
  fetchState: initialFetchState,
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.FETCH:
      const credentialOptions = getDataOrCur(action, state.credentialOptions);
      const credentialElements = credentialOptions === state.credentialOptions ?
        state.credentialElements : buildCredentialElements(credentialOptions);

      return {
        ...state,
        fetchState: fetchStateReducer(action),
        credentialOptions,
        credentialElements,
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action.type;
      return state;
  }
};

const fetchCredentialOptions: ActionCreator<
  ThunkAction<Promise<ActionTypes>, State, void>> = () => {
  return async (dispatch: Dispatch<State>) => {
    dispatch(fetch.start());
    try {
      const data = await requestToMatchService.getCredentialOptions();
      return dispatch(fetch.receive(data));
    } catch(e) {
      return dispatch(fetch.error(e.message));
    }
  };
}

export { fetchCredentialOptions };
