import { combineReducers, createStore, applyMiddleware } from 'redux';
import thunk from 'redux-thunk';

import { conversationsReducer } from './state/conversations';

const store = createStore(combineReducers({
  conversationsReducer,
}), applyMiddleware(thunk));
export default store;
