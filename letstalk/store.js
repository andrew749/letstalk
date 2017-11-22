import { combineReducers, createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';

import { conversationsReducer } from './state/conversations';
import { threadReducer } from './state/thread';

const store = createStore(combineReducers({
  conversationsReducer,
  threadReducer,
}), applyMiddleware(thunk));
export default store;
