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

export type FetchActionCreators<P, D> = {
  receive(data: D): FetchReceiveAction<P, D>;
  error(errorMsg: any): FetchErrorAction<P>;
  start(): FetchStartAction<P>;
}

const initialFetchState: FetchState = { state: 'prefetch' };
export { initialFetchState };
