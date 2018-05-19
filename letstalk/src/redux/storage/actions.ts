import Immutable from 'immutable';

import { Action } from 'redux';

export enum TypeKeys {
  STORE = 'STORAGE_SERVICE/STORE',
  RETRIEVE = 'STORAGE_SERVICE/RETRIEVE'
}

export interface StoreAction extends Action {
  readonly type: TypeKeys.STORE;
  readonly key: string;
  readonly data: any;
}

export interface RetrieveAction extends Action {
  readonly type: TypeKeys.RETRIEVE;
  readonly key: string;
}

export function store(key: string, data: any): StoreAction {
  return {
    type: TypeKeys.STORE,
    key: key,
    data: data,
  };
}

export function retrieve(key: string): RetrieveAction {
  return {
    type: TypeKeys.RETRIEVE,
    key: key,
  };
}
