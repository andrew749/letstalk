export interface Action<T> {
  readonly type: string;
  readonly payload: T;
};

interface ActionCreator<T> {
  readonly type: string;
  (payload: T): Action<T>;
};

export const actionCreator = <T>(type: string): ActionCreator<T> =>
  Object.assign((payload: T): Action<T> => ({type, payload}), {type})

/**
 * Type guard for actions.
 * https://github.com/Microsoft/TypeScript/wiki/What's-new-in-TypeScript#user-defined-type-guard-functions
 */
export const isType = <T>(action: Action<any>, actionCreator: ActionCreator<T>):
  action is Action<T> => action.type === actionCreator.type

