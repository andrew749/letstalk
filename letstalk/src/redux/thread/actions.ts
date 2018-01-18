import Immutable from 'immutable';
import { Action, ActionCreator } from 'redux';

import MessageData from '../../models/message-data';

export enum TypeKeys {
  RECEIVE_MESSAGES = 'RECEIVE_MESSAGES',
}

export interface ReceiveMessagesAction extends Action {
  readonly type: TypeKeys.RECEIVE_MESSAGES;
  readonly messages: Immutable.List<MessageData>;
};

export const receiveMessages: ActionCreator<ReceiveMessagesAction> =
  (messages: Immutable.List<MessageData>) => {
  return {
    type: TypeKeys.RECEIVE_MESSAGES,
    messages,
  };
};

export type ActionTypes =
  | ReceiveMessagesAction
