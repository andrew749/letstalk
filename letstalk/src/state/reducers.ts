import { RECEIVE_MESSAGES } from './types';
import { BASE_URL } from '../services/constants';
import Requestor from '../services/requests';

const initialState = {
  messages: [],
};

export function threadReducer(state = initialState, action) {
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

function receiveMessages(data) {
  return {
    type: RECEIVE_MESSAGES,
    messages: data
  };
};

export function fetchMessages(userId) {
  return dispatch => (new Requestor(BASE_URL)).get('/messages/get')
    .then(data => dispatch(receiveMessages(data)));
};