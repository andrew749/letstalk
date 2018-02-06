import { combineReducers, Reducer } from 'redux';

import { reducer as conversation, State as ConversationState } from './conversation/reducer';
import { reducer as login, State as LoginState } from './login/reducer';
import { reducer as thread, State as ThreadState } from './thread/reducer';
import { receiveMessages } from './thread/actions';

export interface RootState {
  conversation: ConversationState;
  login: LoginState;
  thread: ThreadState;
};

const appReducer: Reducer<RootState> = combineReducers<RootState>({
  conversation,
  login,
  thread,
});

export default appReducer;
