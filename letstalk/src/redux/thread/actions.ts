import Immutable from 'immutable';

import MessageData from '../../models/message-data';
import { actionCreator } from '../actions';

// TODO: Figure out a nice way to make sure these are unique.
// TODO: Use enum???
const RECEIVE_MESSAGES = 'THREAD_RECEIVE_MESSAGES';

interface ReceiveMessagesPayload {
  readonly messages: Immutable.List<MessageData>;
};

export const receiveMessages =
  actionCreator<ReceiveMessagesPayload>(RECEIVE_MESSAGES);
