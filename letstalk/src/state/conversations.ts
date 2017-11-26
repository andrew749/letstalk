import MessageData from '../models/message-data';

// TODO: Add a type
const initialState = {
  threads: [
    new MessageData(1, "Andrew"),
    new MessageData(2, "Adam"),
  ],
};

export function conversationsReducer(state = initialState, action) {
  return state;
}