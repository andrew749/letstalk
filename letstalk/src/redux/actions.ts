import { Action } from 'redux';

import { APIError, ErrorTypes } from '../services/requests';

export const FETCH_STATE_PREFETCH = 'prefetch';
export const FETCH_STATE_FETCHING = 'fetching';
export const FETCH_STATE_ERROR = 'error';
export const FETCH_STATE_SUCCESS = 'success';

export type States = 'prefetch' | 'fetching' | 'error' | 'success';

export type FetchState = {
  readonly state: States;
  readonly errorMsg?: string;
  readonly errorType?: ErrorTypes;
}

export enum FetchTypeKeys {
  RECEIVE = 'RECEIVE',
  ERROR = 'ERROR',
  START = 'START',
}

interface FetchAction<P> extends Action {
  readonly type: P;
  readonly fetchType: string;
}

export interface FetchReceiveAction<P, D> extends FetchAction<P> {
  readonly fetchType: FetchTypeKeys.RECEIVE;
  readonly data: D;
}

export interface FetchErrorAction<P> extends FetchAction<P> {
  readonly type: P;
  readonly fetchType: FetchTypeKeys.ERROR;
  readonly error: APIError,
}

export interface FetchStartAction<P> extends FetchAction<P> {
  readonly type: P;
  readonly fetchType: FetchTypeKeys.START;
}

type FetchActions<P, D> =
  | FetchReceiveAction<P, D>
  | FetchErrorAction<P>
  | FetchStartAction<P>;

export type FetchActionCreators<P, D> = {
  receive(data: D): FetchReceiveAction<P, D>;
  error(errorMsg: any): FetchErrorAction<P>;
  start(): FetchStartAction<P>;
}

export function getDataOrCur<P, D>(action: FetchActions<P, D>, cur: D): D | null {
  switch (action.fetchType) {
    case FetchTypeKeys.RECEIVE:
      return action.data;
    case FetchTypeKeys.ERROR:
    case FetchTypeKeys.START:
      return cur;
    default:
      // Ensure exhaustiveness of select
      const _: never = action;
  }
}

export function fetchStateReducer<P, D>(action: FetchActions<P, D>): FetchState {
  switch (action.fetchType) {
    case FetchTypeKeys.RECEIVE:
      return { state: FETCH_STATE_SUCCESS };
    case FetchTypeKeys.ERROR:
      return {
        state: FETCH_STATE_ERROR,
        errorMsg: action.error.errorMsg,
        errorType: action.error.errorType,
      };
    case FetchTypeKeys.START:
      return { state: FETCH_STATE_FETCHING };
    default:
      // Ensure exhaustiveness of select
      const _: never = action;
  }
}

function combineErrorTypes(fst: ErrorTypes, snd: ErrorTypes): ErrorTypes {
  if (fst === 'UNAUTHORIZED' || snd === 'UNAUTHORIZED') {
    return 'UNAUTHORIZED';
  } else if (fst === 'INVALID_REQUEST' || snd === 'INVALID_REQUEST') {
    return 'INVALID_REQUEST';
  } else {
    return null;
  }
}

export function combineFetchStates(fst: FetchState, snd: FetchState): FetchState {
  const errorMsg = snd.errorMsg ?
    (fst.errorMsg ? fst.errorMsg + ', ' + snd.errorMsg : snd.errorMsg) :
    fst.errorMsg;
  const errorType = combineErrorTypes(fst.errorType, snd.errorType);
  const possibleStates = [
    FETCH_STATE_PREFETCH,
    FETCH_STATE_FETCHING,
    FETCH_STATE_ERROR,
    FETCH_STATE_SUCCESS,
  ] as Array<States>;
  const fstIdx = possibleStates.indexOf(fst.state);
  const sndIdx = possibleStates.indexOf(snd.state);
  const idx = fstIdx < sndIdx ? fstIdx : sndIdx;
  const state = possibleStates[idx];
  return {
    state,
    errorMsg,
    errorType,
  };
}

const initialFetchState: FetchState = { state: FETCH_STATE_PREFETCH };
export { initialFetchState };
