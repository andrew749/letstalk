import { RECEIVE_MESSAGES } from './types';
import Requestor from '../../services/requests';
import { BASE_URL } from '../../services/constants'

function receiveMessages(data: any[]) {
    return {
      type: RECEIVE_MESSAGES,
      messages: data
    };
  };
  
export function fetchMessages(userId: string) {
    return (dispatch: any) => (new Requestor(BASE_URL)).get('/messages/get')
        .then((data: any[]) => dispatch(receiveMessages(data)));
};