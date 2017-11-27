import { combineReducers } from 'redux';

import conversationsReducer from './conversation/reducer';
import threadReducer from './thread/reducer';

const appReducer = combineReducers({
    conversationsReducer,
    threadReducer
});

// Setup root reducer
const rootReducer = (state: any, action: any) => {
    const newState = (action.type === 'RESET') ? undefined : state;
    return appReducer(newState, action);
};

export default rootReducer;
