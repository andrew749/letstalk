import Immutable from 'immutable';

import { receiveMessages } from './actions';
import MessageData from '../../models/message-data';
import { Action, isType } from '../actions';

export type State = {
  messages: Immutable.List<MessageData>,
};

const initialState: State = {
  messages: Immutable.List<MessageData>(),
};

export function reducer(state: State = initialState, action: Action<any>): State {
  if (isType(action, receiveMessages)) {
    return {
      ...state,
      messages: action.payload.messages,
    };
  }

  return state;
};
