import { RECEIVE_MESSAGES } from './types';

interface State {
  messages: any[]
}

const initialState: State = {
  messages: []
}

export default function threadReducer(state = initialState, action: any) {
  switch (action.type) {
    case RECEIVE_MESSAGES: 
      return {
        ...state,
        messages: action.messages,
      };
    default:
      return state;
  };
};