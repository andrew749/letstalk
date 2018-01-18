import Immutable from 'immutable';
import { Action, ActionCreator } from 'redux';

import MessageData from '../../models/message-data';

export enum TypeKeys {
  RECEIVE_CONVERSATIONS = 'RECEIVE_CONVERSATIONS',
}

interface ReceiveConversationsAction extends Action {
  readonly type: TypeKeys.RECEIVE_CONVERSATIONS;
  readonly threads: Immutable.List<MessageData>;
};

export const receiveConversations: ActionCreator<ReceiveConversationsAction> =
  (threads: Immutable.List<MessageData>) => {
  return {
    type: TypeKeys.RECEIVE_CONVERSATIONS,
    threads,
  };
};

export type ActionTypes =
  | ReceiveConversationsAction
