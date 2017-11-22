import MessageData from '../models/message-data';

// TODO: Add a type
const initialState = {
  threads: [
    new MessageData(1, "Andrew"),
    new MessageData(2, "Adam"),
  ],
};

function conversationsReducer(state = initialState, action) {
  return state;
}

export {
  conversationsReducer,
};
