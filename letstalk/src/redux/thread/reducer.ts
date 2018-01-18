import Immutable from 'immutable';
import { Action } from 'redux';

import MessageData from '../../models/message-data';
import { ActionTypes, TypeKeys } from './actions';

export type State = {
  messages: Immutable.List<MessageData>,
};

const initialState: State = {
  messages: Immutable.List<MessageData>(),
};

export function reducer(state: State = initialState, action: ActionTypes): State {
  switch (action.type) {
    case TypeKeys.RECEIVE_MESSAGES:
      return {
        ...state,
        messages: action.messages,
      };
    default:
      // Ensure exhaustiveness of select
      const _: never = action.type;
      return state;
  }
};
