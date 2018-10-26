import { Action } from 'redux'

import { UserSearchResponse } from '../../models/user-search';
import {
  FetchReceiveAction,
  FetchErrorAction,
  FetchStartAction,
  FetchTypeKeys,
  FetchActionCreators,
} from '../actions';
import { APIError } from '../../services/requests';
import { CohortV2 } from '../../models/cohort';
import { Position } from '../../models/position';
import { Group } from '../../models/group';
import { SimpleTraitDenormalized } from '../../models/simple-trait';

export enum TypeKeys {
  FETCH = 'CREDENTIALS/FETCH',
  SET_CURRENT_QUERY = 'CREDENTIALS/SET_CURRENT_QUERY',
}

export enum QueryTypes {
  SEARCH_COHORT = "SEARCH_COHORT",
  SEARCH_POSITION = "SEARCH_POSITION",
  SEARCH_SIMPLE_TRAIT = "SEARCH_SIMPLE_TRAIT",
  SEARCH_GROUP = "SEARCH_GROUP",
  YOUR_COHORT = "YOUR_COHORT",
}

interface SearchCohortQuery extends CohortV2 {
  type: QueryTypes.SEARCH_COHORT;
};

interface SearchPositionQuery extends Position {
  type: QueryTypes.SEARCH_POSITION;
};

interface SearchSimpleTraitQuery extends SimpleTraitDenormalized {
  type: QueryTypes.SEARCH_SIMPLE_TRAIT;
};

interface SearchGroupQuery extends Group {
  type: QueryTypes.SEARCH_GROUP;
};

interface YourCohortRecommendation {
  type: QueryTypes.YOUR_COHORT;
};

export type Query =
  | SearchCohortQuery
  | SearchPositionQuery
  | SearchSimpleTraitQuery
  | SearchGroupQuery
  | YourCohortRecommendation

type UserSearchReceiveAction = FetchReceiveAction<TypeKeys.FETCH, UserSearchResponse>;
type UserSearchErrorAction = FetchErrorAction<TypeKeys.FETCH>;
type UserSearchStartAction = FetchStartAction<TypeKeys.FETCH>;

export interface SetCurrentQueryAction extends Action {
  readonly type: TypeKeys.SET_CURRENT_QUERY;
  readonly query: Query;
}

function receive(data: UserSearchResponse): UserSearchReceiveAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.RECEIVE,
    data,
  };
}

function error(error: APIError): UserSearchErrorAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.ERROR,
    error,
  };
}

function start(): UserSearchStartAction {
  return {
    type: TypeKeys.FETCH,
    fetchType: FetchTypeKeys.START,
  };
}

export function setCurrentQuery(query: Query): SetCurrentQueryAction {
  return {
    type: TypeKeys.SET_CURRENT_QUERY,
    query,
  };
}

const fetch: FetchActionCreators<TypeKeys.FETCH, UserSearchResponse> = {
  receive,
  error,
  start,
};

export { fetch }

export type ActionTypes =
  | UserSearchReceiveAction
  | UserSearchErrorAction
  | UserSearchStartAction
  | SetCurrentQueryAction
