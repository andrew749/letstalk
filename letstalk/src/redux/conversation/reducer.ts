import Immutable from 'immutable';

import MessageData from '../../models/message-data';
import { Action } from '../actions';

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

export function reducer(state: State = initialState, action: Action<any>) {
  return state;
}
