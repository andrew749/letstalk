import Immutable from 'immutable';
import { Action } from 'redux';

import MessageData from '../../models/message-data';
import { ActionTypes, TypeKeys } from './actions';

export type State = {
  threads: Immutable.List<MessageData>,
};

const initialState: State = {
  threads: Immutable.List<MessageData>([
    new MessageData(1, "Andrew"),
    new MessageData(2, "Adam"),
    new MessageData(3, "Wojtek"),
  ]),
};

// TODO: Change Action to a union of actual possible action types, when they exist
export function reducer(state: State = initialState, action: ActionTypes) {
  switch (action.type) {
    case TypeKeys.RECEIVE_CONVERSATIONS:
      return {
        ...state,
        threads: action.threads,
      }
    default:
      // Ensure exhaustiveness of select
      const _: never = action.type;
      return state;
  }
}
