import { Action } from 'redux';

export type FetchState = {
  state: 'prefetch' | 'fetching' | 'error' | 'success';
  errorMsg?: string;
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
  readonly errorMsg: string;
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
      return { state: 'success' };
    case FetchTypeKeys.ERROR:
      return { state: 'error', errorMsg: action.errorMsg };
    case FetchTypeKeys.START:
      return { state: 'fetching' };
    default:
      // Ensure exhaustiveness of select
      const _: never = action;
  }
}

const initialFetchState: FetchState = { state: 'prefetch' };
export { initialFetchState };
